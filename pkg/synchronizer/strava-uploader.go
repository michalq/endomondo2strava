package synchronizer

import (
	"fmt"
	"log"
	"strings"
	"time"

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
	logger             func(string)
}

// NewStravaUploader creates instance of StravaUploader
func NewStravaUploader(stravaClient *strava.Client, workoutsRepository Workouts, logger func(string)) *StravaUploader {
	return &StravaUploader{stravaClient, workoutsRepository, logger}
}

// UploadAll uploads all provided workouts to strava
func (s *StravaUploader) UploadAll() error {
	workouts, err := s.workoutsRepository.FindAll()
	if err != nil {
		return err
	}
	workoutCount := len(workouts)
	var toImport []Workout
	for _, workout := range workouts {
		if workout.UploadStarted == 0 {
			toImport = append(toImport, workout)
		}
	}

	totalUploaded := 0
	for len(toImport) > 0 {
		uploaded, err := s.uploadMany(toImport)
		totalUploaded += len(uploaded)
		fmt.Printf("\n---\nUploaded: %d, Skipped: %d (due to pending or ended import), All: %d\n",
			totalUploaded, workoutCount, workoutCount-totalUploaded)

		if err != nil {
			rateLimitExceeded := strings.Contains(err.Error(), "Rate Limit Exceeded")
			if rateLimitExceeded {
				log.Println(err.Error())
			} else {
				return err
			}
		}

		for _, workout := range uploaded {
			if updateErr := s.workoutsRepository.Update(&workout); updateErr != nil {
				fmt.Println("Err", updateErr)
			}
			workouts, err := s.workoutsRepository.FindAll()
			if err != nil {
				return err
			}
			for _, workout := range workouts {
				if workout.UploadStarted == 0 {
					toImport = append(toImport, workout)
				}
			}
		}

		if len(toImport) > 0 {
			if totalUploaded%1000 == 0 {
				log.Println("Exceeded api limit of 1000 calls per day. You can cancel and restart or")
				log.Println("wait for 24 hours and then this will continue automatically")
				time.Sleep(time.Hour * 24)
			} else {
				log.Println("Exceeded api limit of 100 calls per 15 mins. You can cancel and restart or")
				log.Println("wait for 15 minutes and then this will continue automatically")
				time.Sleep(time.Minute * 15)
			}
		}
	}

	return nil
}

func (s *StravaUploader) uploadMany(workouts []Workout) ([]Workout, error) {
	uploadedChan := make(chan Workout)
	errorsChan := make(chan error)

	for i := 0; i < len(workouts) && i < 100; i++ {
		go s.uploadSingleWorkout(workouts[i], uploadedChan, errorsChan)
	}
	var uploaded []Workout
	for range workouts {
		select {
		case workout := <-uploadedChan:
			uploaded = append(uploaded, workout)
			s.logger(fmt.Sprintf("Sent workout to strava, endomondo id %s, strava id %s", workout.EndomondoID, workout.StravaID))
		case err := <-errorsChan:
			return uploaded, err
		}
	}
	return uploaded, nil
}

func (s *StravaUploader) uploadSingleWorkout(
	workout Workout,
	uploadedChan chan<- Workout,
	errorsChan chan<- error,
) {
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
		errorsChan <- err
		return
	}
	workout.StravaID = uploadResponse.ID
	workout.UploadStarted = 1
	uploadedChan <- workout
}
