package export

import (
	"fmt"
	"math"
	"strconv"

	"github.com/michalq/endo2strava/internal/modules/workouts"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// Exporter does export of workouts
type Exporter struct {
	endomondoDownloader *Downloader
	workoutsRepository  workouts.Workouts
}

// NewExporter creates new instance of Exporter
func NewExporter(endomondoDownloader *Downloader, workoutsRepository workouts.Workouts) *Exporter {
	return &Exporter{endomondoDownloader, workoutsRepository}
}

// RetrieveWorkouts finds and save in database all existing workouts and save workout file in given format
func (e *Exporter) RetrieveWorkouts(authorizedClient *endomondo.Client, format string) error {

	allWorkouts, err := e.fetchAllWorkoutsByPage(authorizedClient, 100)
	if err != nil {
		return err
	}
	if err := e.workoutsRepository.SaveAll(allWorkouts); err != nil {
		return fmt.Errorf("Error while saving workouts to db [%s]", err)
	}
	downloadedWorkouts, err := e.endomondoDownloader.DownloadAll(authorizedClient, format, allWorkouts)
	if err != nil {
		return err
	}
	if err := e.workoutsRepository.SaveAll(downloadedWorkouts); err != nil {
		return fmt.Errorf("Error while updating workouts in db [%s]", err)
	}
	return nil
}

// RetrieveDetails search workout by workout to retrieve rest of necessary data like title, description, pictures etc
func (e *Exporter) RetrieveDetails(authorizedClient *endomondo.Client) error {
	return nil
}

func (e *Exporter) fetchAllWorkoutsByPage(authorizedClient *endomondo.Client, workoutsPageLimit int) ([]workouts.Workout, error) {

	var allWorkouts []workouts.Workout
	workoutsResponse, err := e.fetchWorkouts(authorizedClient, 0, 1)
	if err != nil {
		return nil, err
	}
	pages := int(math.Ceil(float64(workoutsResponse.Paging.Total) / float64(workoutsPageLimit)))
	workoutsResponseChan := make(chan *endomondo.WorkoutsResponse)
	errChan := make(chan error)
	for page := 0; page < pages; page++ {
		go func(workoutsResponseChan chan<- *endomondo.WorkoutsResponse, errChan chan<- error, page int) {
			workoutsResponse, err := e.fetchWorkouts(authorizedClient, page, workoutsPageLimit)
			if err != nil {
				errChan <- err
				return
			}
			workoutsResponseChan <- workoutsResponse
		}(workoutsResponseChan, errChan, page)
	}

	for page := 0; page < pages; page++ {
		select {
		case workoutsResponse := <-workoutsResponseChan:
			for _, workout := range workoutsResponse.Data {
				allWorkouts = append(allWorkouts, workouts.Workout{EndomondoID: strconv.FormatInt(workout.ID, 10)})
			}
		case _ = <-errChan:
			//
		}
	}

	return allWorkouts, nil
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
