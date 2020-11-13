package synchronizer

import "github.com/michalq/endo2strava/pkg/strava-client"

// StravaUploader uploads workouts into strava
type StravaUploader struct {
	stravaClient *strava.Client
}

// NewStravaUploader creates instance of StravaUploader
func NewStravaUploader(stravaClient *strava.Client) *StravaUploader {
	return &StravaUploader{stravaClient}
}

// UploadAll uploads all provided workouts to strava
func (s *StravaUploader) UploadAll(workouts []Workout) int {
	return 0
}
