package synchronizer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
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
	workoutsPath    string
	workoutsFormat  string
}

// NewEndomondoDownloader creates new instance
func NewEndomondoDownloader(endomondoClient *endomondo.Client, workoutsPath, workoutsFormat string) *EndomondoDownloader {
	return &EndomondoDownloader{endomondoClient, workoutsPath, workoutsFormat}
}

// DownloadAllBetween downloads all workouts between provided dates
// Finds month by month from start date to end date, cause endomondo has problem with longer periods
func (e *EndomondoDownloader) DownloadAllBetween(startAt, endAt time.Time) {
	resultsChan := make(chan Result)
	errorsChan := make(chan error)
	startTime := startAt
	iterations := 0

	for startTime.Before(endAt) {
		iterations++
		endTime := startTime.AddDate(0, 1, 0)
		go e.fetchWorkoutsBetween(startTime, endTime, resultsChan, errorsChan)
		startTime = startTime.AddDate(0, 1, 0)
	}

	var wg sync.WaitGroup
	workoutsDownloaded, allWorkouts := 0, 0
	for i := 0; i < iterations; i++ {
		select {
		case result := <-resultsChan:
			fmt.Printf("Between %s and %s found %d workouts\n", result.From.Format("2006-01-02"), result.To.Format("2006-01-02"), len(result.Workouts))
			allWorkouts += len(result.Workouts)
			for _, workout := range result.Workouts {
				wg.Add(1)
				go func(workoutID int64) {
					if err := e.DownloadWorkout(workoutID); err != nil {
						fmt.Println("Err", err)
					} else {
						workoutsDownloaded++
						fmt.Printf("Downloaded workout ID %d\t[%d / %d]\n", workoutID, workoutsDownloaded, allWorkouts)
					}
				}(workout.ID)
			}
		case err := <-errorsChan:
			fmt.Println("Err", err)
		}
	}
	wg.Wait()
	fmt.Printf("\n---\nAll done.\nDownloaded %d workouts out of %d", workoutsDownloaded, allWorkouts)
}

func (e *EndomondoDownloader) fetchWorkoutsBetween(
	startTime,
	endTime time.Time,
	resultsChan chan<- Result,
	errorsChan chan<- error,
) {

	workouts, err := e.endomondoClient.Workouts(endomondo.WorkoutsQueryParams{
		After:  startTime.Format(time.RFC3339),
		Before: endTime.Format(time.RFC3339),
	})
	if err != nil {
		errorsChan <- err
		return
	}
	resultsChan <- Result{From: startTime, To: endTime, Workouts: workouts.Data}
}

// DownloadWorkout performs downloading single workout
func (e *EndomondoDownloader) DownloadWorkout(workoutID int64) error {
	workoutBuf, err := e.endomondoClient.ExportWorkout(workoutID, e.workoutsFormat)
	defer workoutBuf.Close()
	if err != nil {
		return err
	}
	fullPath := fmt.Sprintf("%s/%d.%s", e.workoutsPath, workoutID, strings.ToLower(e.workoutsFormat))
	out, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Err", err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, workoutBuf)
	return nil
}
