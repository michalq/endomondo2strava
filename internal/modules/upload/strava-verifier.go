package upload

import (
	"fmt"
	"strconv"

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
	toVerification := make([]workouts.Workout, 0)
	for _, workout := range allWorkouts {
		if workout.StravaID == "" || workout.UploadEnded == 1 {
			continue
		}
		if len(toVerification) >= requestLimit {
			break
		}
		toVerification = append(toVerification, workout)
	}
	verifiedChan := make(chan workouts.Workout)
	errorsChan := make(chan error)
	for _, workout := range toVerification {
		go func(workout workouts.Workout) {
			resp, err := authorizedClient.GetUpload(workout.StravaID)
			if err != nil {
				errorsChan <- err
				return
			}
			workout.StravaActivityID = strconv.Itoa(resp.ActivityID)
			workout.StravaError = resp.Error
			workout.StravaStatus = resp.Status
			workout.UploadEnded = 1
			verifiedChan <- workout
		}(workout)
	}
	verifiedWorkouts := make([]workouts.Workout, 0)
	for range toVerification {
		select {
		case workout := <-verifiedChan:
			verifiedWorkouts = append(verifiedWorkouts, workout)
			s.logger.Info(fmt.Sprintf(
				"Workout %s verified, Status: %s, Error: %s",
				workout.StravaActivityID, workout.StravaStatus, workout.StravaError,
			))
		case err := <-errorsChan:
			s.logger.Warning(fmt.Sprintf("Error while verification (%s).", err.Error()))
		}
	}

	if err := s.workoutsRepository.SaveAll(verifiedWorkouts); err != nil {
		return err
	}

	return nil
}
