package export

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/michalq/endo2strava/internal/modules/common"
	"github.com/michalq/endo2strava/internal/modules/users"

	"github.com/michalq/endo2strava/internal/modules/workouts"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// Status status of export
type Status struct {
	// All how many workouts found
	All int
	// Downloaded how many workouts were downloaded
	Downloaded int
}

// Exporter does export of workouts
type Exporter struct {
	endomondoDownloader *Downloader
	workoutsRepository  workouts.Workouts
	usersRepository     users.Users
	logger              common.Logger
}

// NewExporter creates new instance of Exporter
func NewExporter(endomondoDownloader *Downloader, workoutsRepository workouts.Workouts, usersRepository users.Users, logger common.Logger) *Exporter {
	return &Exporter{endomondoDownloader, workoutsRepository, usersRepository, logger}
}

// FindWorkoutsQuantity simply finds workouts quantity
func (e *Exporter) FindWorkoutsQuantity(authorizedClient *endomondo.Client, user *users.User) (int, error) {
	workoutsResponse, err := e.fetchWorkouts(authorizedClient, 0, 1)
	if err != nil {
		return 0, err
	}
	return workoutsResponse.Paging.Total, nil
}

// FindWorkouts only finds general information about workouts and save them in database
func (e *Exporter) FindWorkouts(authorizedClient *endomondo.Client, workoutsQuantity int) error {

	foundWorkouts, err := e.workoutsRepository.FindAll()
	if err != nil {
		return err
	}
	if len(foundWorkouts) == workoutsQuantity {
		return nil
	}
	allWorkouts, err := e.fetchAllWorkoutsByPage(authorizedClient, 100, workoutsQuantity)
	if err != nil {
		return err
	}
	if err := e.workoutsRepository.SaveAll(allWorkouts); err != nil {
		return fmt.Errorf("Error while saving workouts to db [%s]", err)
	}
	return nil
}

// DownloadWorkouts downloads not downloaded workouts
func (e *Exporter) DownloadWorkouts(authorizedClient *endomondo.Client, format string) error {

	var workoutsToDownload []workouts.Workout
	allWorkouts, err := e.workoutsRepository.FindAll()
	if err != nil {
		return err
	}
	for _, workout := range allWorkouts {
		if workout.Path == "" {
			workoutsToDownload = append(workoutsToDownload, workout)
		}
	}
	downloadedWorkouts, err := e.endomondoDownloader.DownloadAll(authorizedClient, format, workoutsToDownload)
	if err != nil {
		return err
	}
	if err := e.workoutsRepository.SaveAll(downloadedWorkouts); err != nil {
		return fmt.Errorf("Error while updating workouts in db [%s]", err)
	}
	return nil
}

// FindWorkoutsDetails search workout by workout to retrieve rest of necessary data like title, description, pictures etc
func (e *Exporter) FindWorkoutsDetails(authorizedClient *endomondo.Client) error {

	allWorkouts, err := e.workoutsRepository.FindAll()
	if err != nil {
		return err
	}

	workoutsChan := make(chan workouts.Workout)
	workoutErrChan := make(chan error)
	var processed []workouts.Workout
	workoutsUpdated := 0
	for _, workout := range allWorkouts {
		if workout.DetailsExported == 1 {
			continue
		}
		workoutsUpdated++
		go func(workout workouts.Workout, workoutsChan chan<- workouts.Workout, workoutErrChan chan<- error) {
			details, err := authorizedClient.Workout(workout.EndomondoIDAsInt())
			if err != nil {
				workoutErrChan <- err
				return
			}
			var pictures []string
			for _, picture := range details.Pictures {
				pictures = append(pictures, picture.URL)
			}
			workout.Title = details.Title
			workout.Description = details.Message
			workout.Hashtags = strings.Join(details.Hashtags, ",")
			workout.Pictures = strings.Join(pictures, ",")
			workout.DetailsExported = 1
			workoutsChan <- workout
		}(workout, workoutsChan, workoutErrChan)
	}

	for i := 0; i < workoutsUpdated; i++ {
		select {
		case workout := <-workoutsChan:
			processed = append(processed, workout)
			e.logger.Info(fmt.Sprintf("Found details for %s", workout.EndomondoID))
		case err := <-workoutErrChan:
			e.logger.Warning(fmt.Sprintln("Err", err))
		}
	}
	if err := e.workoutsRepository.SaveAll(processed); err != nil {
		return fmt.Errorf("Error while saving workouts to db [%s]", err)
	}
	return nil
}

func (e *Exporter) fetchAllWorkoutsByPage(authorizedClient *endomondo.Client, workoutsPageLimit, workoutsQuantity int) ([]workouts.Workout, error) {

	var allWorkouts []workouts.Workout
	pages := int(math.Ceil(float64(workoutsQuantity) / float64(workoutsPageLimit)))
	workoutsResponseChan := make(chan *endomondo.WorkoutsResponse)
	errChan := make(chan error)
	for page := 0; page < pages; page++ {
		go func(workoutsResponseChan chan<- *endomondo.WorkoutsResponse, errChan chan<- error, page int) {
			workoutsResponse, err := e.fetchWorkouts(authorizedClient, page*workoutsPageLimit, workoutsPageLimit)
			if err != nil {
				errChan <- err
				return
			}
			workoutsResponseChan <- workoutsResponse
		}(workoutsResponseChan, errChan, page)
	}

	errCollection := common.NewErrorCollection()
	for page := 0; page < pages; page++ {
		select {
		case workoutsResponse := <-workoutsResponseChan:
			for _, workout := range workoutsResponse.Data {
				allWorkouts = append(allWorkouts, workouts.Workout{EndomondoID: strconv.FormatInt(workout.ID, 10)})
			}
		case err := <-errChan:
			e.logger.Warning(err.Error())
			errCollection.Append(err)
		}
	}

	return allWorkouts, errCollection
}

func (e *Exporter) fetchWorkouts(authorizedClient *endomondo.Client, offset, limit int) (*endomondo.WorkoutsResponse, error) {

	workoutsResponse, err := authorizedClient.Workouts(endomondo.WorkoutsQueryParams{
		Offset: &offset,
		Limit:  &limit,
	})
	if err != nil {
		return nil, err
	}
	return workoutsResponse, nil
}
