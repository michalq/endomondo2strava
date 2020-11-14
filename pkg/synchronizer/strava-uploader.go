package synchronizer

import (
	"fmt"

	"github.com/michalq/endo2strava/pkg/strava-client"
)

// UploadStatus status of current import
type UploadStatus struct {
	// Uploaded in current run
	Uploaded int
	// Skipped due to pending import
	Skipped int
	// All workouts imported and not imported
	All int
}

// StravaUploader uploads workouts into strava
type StravaUploader struct {
	stravaClient       *strava.Client
	workoutsRepository Workouts
}

// NewStravaUploader creates instance of StravaUploader
func NewStravaUploader(stravaClient *strava.Client, workoutsRepository Workouts) *StravaUploader {
	return &StravaUploader{stravaClient, workoutsRepository}
}

// UploadAll uploads all provided workouts to strava
func (s *StravaUploader) UploadAll() (*UploadStatus, error) {
	workouts, err := s.workoutsRepository.FindAll()
	if err != nil {
		return nil, err
	}
	var toImport []Workout
	for _, workout := range workouts {
		if workout.UploadStarted == 0 {
			toImport = append(toImport, workout)
		}
	}

	uploaded, err := s.uploadMany(toImport)
	if len(uploaded) == 0 && err != nil {
		return nil, err
	}
	for _, workout := range uploaded {
		if err := s.workoutsRepository.Update(&workout); err != nil {
			fmt.Println("Err", err)
		}
	}
	return &UploadStatus{
		Uploaded: len(uploaded),
		All:      len(workouts),
		Skipped:  len(workouts) - len(toImport),
	}, err
}

func (s *StravaUploader) uploadMany(workouts []Workout) ([]Workout, error) {
	var uploaded []Workout
	for _, workout := range workouts {
		uploadResponse, err := s.stravaClient.ImportWorkout(strava.UploadParameters{
			ExternalID:  workout.EndomondoID,
			Name:        fmt.Sprintf("Endomondo %s", workout.EndomondoID),
			Description: fmt.Sprintf("Workout imported from endomondo"),
			File:        workout.Path,
			Commute:     "0",
			DataType:    workout.Ext,
			Trainer:     "0",
		})
		if err != nil {
			return nil, err
		}
		workout.StravaID = uploadResponse.ID
		workout.UploadStarted = 1
		uploaded = append(uploaded, workout)
	}
	return uploaded, nil
}
