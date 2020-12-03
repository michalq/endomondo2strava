package upload

import (
	"fmt"

	"github.com/michalq/endo2strava/internal/modules/common"
	"github.com/michalq/endo2strava/internal/modules/workouts"
	"github.com/michalq/endo2strava/pkg/strava-client"
)

// StravaVerifier is responsible for verifying uploaded workouts
type StravaVerifier struct {
	workoutsRepository workouts.Workouts
	logger             common.Logger
}

// NewStravaVerifier creates new instance of strava verifier
func NewStravaVerifier(workoutsRepository workouts.Workouts, logger common.Logger) *StravaVerifier {
	return &StravaVerifier{workoutsRepository, logger}
}

// Verify checks all workouts that are started uploading if uploads ends
func (s *StravaVerifier) Verify(authorizedClient *strava.Client, requestLimit int) error {

	allWorkouts, err := s.workoutsRepository.FindAll()
	if err != nil {
		return err
	}
	for _, workout := range allWorkouts {
		if workout.StravaID == "" {
			continue
		}
		resp, err := authorizedClient.GetUpload(workout.StravaID)
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", resp)
		return nil
	}

	return nil
}
