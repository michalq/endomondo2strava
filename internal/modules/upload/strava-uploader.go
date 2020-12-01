package upload

import (
	"fmt"
	"strings"

	"github.com/michalq/endo2strava/internal/modules/workouts"
	"github.com/michalq/endo2strava/pkg/strava-client"
)

// Status status of current import
type Status struct {
	// Uploaded in current run
	Uploaded int
	// Skipped due to pending import
	Skipped int
	// All workouts imported and not imported
	All int
}

// StravaUploader uploads workouts into strava
type StravaUploader struct {
	workoutsRepository workouts.Workouts
	logger             func(string)
}

// NewStravaUploader creates instance of StravaUploader
func NewStravaUploader(workoutsRepository workouts.Workouts, logger func(string)) *StravaUploader {
	return &StravaUploader{workoutsRepository, logger}
}

// UploadAll uploads all provided workouts to strava
func (s *StravaUploader) UploadAll(authorizedClient *strava.Client, requestLimit int) (*Status, error) {
	allWorkouts, err := s.workoutsRepository.FindAll()
	if err != nil {
		return nil, err
	}
	var toImport []workouts.Workout
	for _, workout := range allWorkouts {
		if workout.UploadStarted == 0 {
			toImport = append(toImport, workout)
		}
	}
	if len(toImport) > requestLimit {
		toImport = toImport[0:requestLimit]
	}

	uploaded, err := s.uploadMany(authorizedClient, toImport)
	if len(uploaded) == 0 && err != nil {
		return nil, err
	}
	for _, workout := range uploaded {
		if err := s.workoutsRepository.Update(&workout); err != nil {
			fmt.Println("Err", err)
		}
	}
	return &Status{
		Uploaded: len(uploaded),
		All:      len(allWorkouts),
		Skipped:  len(allWorkouts) - len(toImport),
	}, err
}

func (s *StravaUploader) uploadMany(authorizedClient *strava.Client, workoutsToUpload []workouts.Workout) ([]workouts.Workout, error) {
	uploadedChan := make(chan workouts.Workout)
	errorsChan := make(chan error)
	for _, workout := range workoutsToUpload {
		go s.uploadSingleWorkout(authorizedClient, workout, uploadedChan, errorsChan)
	}
	var uploaded []workouts.Workout
	var errCollection = NewErrorCollection()
	for range workoutsToUpload {
		select {
		case workout := <-uploadedChan:
			uploaded = append(uploaded, workout)
			s.logger(fmt.Sprintf("Send workout to strava, endomondo id %s, strava id %s", workout.EndomondoID, workout.StravaID))
		case err := <-errorsChan:
			// Collect errors, but don't stop rest of routines
			s.logger(fmt.Sprintf("Error occured while uploading %s", err))
			errCollection.Append(err)
		}
	}
	return uploaded, errCollection
}

func (s *StravaUploader) uploadSingleWorkout(
	authorizedClient *strava.Client,
	workout workouts.Workout,
	uploadedChan chan<- workouts.Workout,
	errorsChan chan<- error,
) {
	var hashtags []string
	if len(workout.HashtagsList()) > 0 {
		for _, hashtag := range workout.HashtagsList() {
			hashtags = append(hashtags, "#"+hashtag)
		}
	}
	uploadResponse, err := authorizedClient.ImportWorkout(strava.UploadParameters{
		ExternalID:  workout.EndomondoID,
		Name:        fmt.Sprintf(workout.Title),
		Description: fmt.Sprintf("%s %s (Endomondo id %s)", workout.Description, strings.Join(hashtags, " "), workout.EndomondoID),
		File:        workout.Path,
		Commute:     "0",
		DataType:    workout.Ext,
		Trainer:     "0",
	})
	if err != nil {
		errorsChan <- err
		return
	}
	workout.StravaID = uploadResponse.ID
	workout.UploadStarted = 1
	uploadedChan <- workout
}
