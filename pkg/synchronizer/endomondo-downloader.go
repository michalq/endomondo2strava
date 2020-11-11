package synchronizer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// Result represent single resoult of scrapping
type Result struct {
	From     time.Time
	To       time.Time
	Workouts []endomondo.WorkoutsResponseData
}

// EndomondoDownloader finds workouts
type EndomondoDownloader struct {
	endomondoClient *endomondo.Client
}

// NewEndomondoDownloader creates new instance
func NewEndomondoDownloader(endomondoClient *endomondo.Client) *EndomondoDownloader {
	return &EndomondoDownloader{endomondoClient}
}

// FindAllBetween finds all workouts between provided dates
// Finds month by month from start date to end date, cause endomondo has problem with longer periods
func (e *EndomondoDownloader) FindAllBetween(startAt, endAt time.Time) (int, chan Result, chan error) {
	resultsChan := make(chan Result)
	errorsChan := make(chan error)
	startTime := startAt
	iterations := 0
	for startTime.Before(endAt) {
		endTime := startTime.AddDate(0, 1, 0)
		go func(startTime, endTime time.Time) {
			workouts, err := e.endomondoClient.Workouts(endomondo.WorkoutsQueryParams{
				After:  startTime.Format(time.RFC3339),
				Before: endTime.Format(time.RFC3339),
			})
			if err != nil {
				errorsChan <- err
				return
			}
			resultsChan <- Result{From: startTime, To: endTime, Workouts: workouts.Data}
		}(startTime, endTime)
		startTime = startTime.AddDate(0, 1, 0)
		iterations++
	}

	return iterations, resultsChan, errorsChan
}

// DownloadWorkouts downloads all workouts returned by endomondo
func (e *EndomondoDownloader) DownloadWorkouts(filePath string, workouts []endomondo.WorkoutsResponseData, format string) {
	for _, workout := range workouts {
		e.DownloadWorkout(filePath, workout.ID, format)
	}
}

// DownloadWorkout performs downloading single workout
func (e *EndomondoDownloader) DownloadWorkout(filePath string, workoutID int64, format string) (string, error) {
	workoutBuf, err := e.endomondoClient.ExportWorkout(workoutID, format)
	defer workoutBuf.Close()
	if err != nil {
		return "", err
	}
	fullPath := fmt.Sprintf("%s/%d.%s", filePath, workoutID, strings.ToLower(format))
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Err", err)
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, workoutBuf)
	return fullPath, nil
}
