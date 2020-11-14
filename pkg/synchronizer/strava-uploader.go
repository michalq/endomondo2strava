package synchronizer

import (
	"fmt"

	"github.com/michalq/endo2strava/pkg/strava-client"
)

// StravaUploader uploads workouts into strava
type StravaUploader struct {
	stravaClient *strava.Client
}

// NewStravaUploader creates instance of StravaUploader
func NewStravaUploader(stravaClient *strava.Client) *StravaUploader {
	return &StravaUploader{stravaClient}
}

// UploadAll uploads all provided workouts to strava
func (s *StravaUploader) UploadAll(workouts []Workout) ([]Workout, error) {
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
		return uploaded, nil
	}
	return uploaded, nil
}
