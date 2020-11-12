package synchronizer

import "github.com/michalq/endo2strava/pkg/strava-client"

// StravaUploader uploads workouts into strava
type StravaUploader struct {
	stravaClient *strava.Client
}
