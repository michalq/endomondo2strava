package synchronizer

import (
	"fmt"
	"io"
	"os"
	"strconv"
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
	workoutsPath    string
	workoutsFormat  string
	logger          func(l string)
}

// NewEndomondoDownloader creates new instance
func NewEndomondoDownloader(
	endomondoClient *endomondo.Client,
	workoutsPath,
	workoutsFormat string,
	logger func(l string),
) *EndomondoDownloader {
	return &EndomondoDownloader{endomondoClient, workoutsPath, workoutsFormat, logger}
}

// DownloadAllBetween downloads all workouts between provided dates
// Finds month by month from start date to end date, cause endomondo has problem with longer periods
func (e *EndomondoDownloader) DownloadAllBetween(startAt, endAt time.Time) []Workout {
	var downloadedWorkouts []Workout
	resultsChan := make(chan Result)
	errorsChan := make(chan error)
	startTime := startAt
	workoutsListRoutines := 0

	for startTime.Before(endAt) {
		workoutsListRoutines++
		endTime := startTime.AddDate(0, 1, 0)
		go e.fetchWorkoutsBetween(startTime, endTime, resultsChan, errorsChan)
		startTime = startTime.AddDate(0, 1, 0)
	}

	workoutsChan := make(chan Workout)
	workoutErrorChan := make(chan error)
	allWorkouts := 0
	for i := 0; i < workoutsListRoutines; i++ {
		select {
		case result := <-resultsChan:
			e.logger(fmt.Sprintf("Between %s and %s found %d workouts", result.From.Format("2006-01-02"), result.To.Format("2006-01-02"), len(result.Workouts)))
			allWorkouts += len(result.Workouts)
			for _, workout := range result.Workouts {
				go e.downloadSingleWorkout(workout.ID, workoutsChan, workoutErrorChan)
			}
		case err := <-errorsChan:
			fmt.Println("Err", err)
		}
	}

	for i := 0; i < allWorkouts; i++ {
		select {
		case workout := <-workoutsChan:
			e.logger(fmt.Sprintf("Downloaded workout %s", workout.ID))
			downloadedWorkouts = append(downloadedWorkouts, workout)
		case err := <-workoutErrorChan:
			e.logger(fmt.Sprintln("Err", err))
		}
	}

	e.logger(fmt.Sprintf("Downloaded %d workouts out of %d", len(downloadedWorkouts), allWorkouts))

	return downloadedWorkouts
}

func (e *EndomondoDownloader) downloadSingleWorkout(
	workoutID int64,
	workoutChan chan<- Workout,
	errorChan chan<- error,
) {
	workout, err := e.downloadWorkout(workoutID)
	if err != nil {
		errorChan <- err
	} else {
		workoutChan <- *workout
	}
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

func (e *EndomondoDownloader) downloadWorkout(workoutID int64) (*Workout, error) {
	workoutBuf, err := e.endomondoClient.ExportWorkout(workoutID, e.workoutsFormat)
	defer workoutBuf.Close()
	if err != nil {
		return nil, err
	}
	fullPath := fmt.Sprintf("%s/%d.%s", e.workoutsPath, workoutID, strings.ToLower(e.workoutsFormat))
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(out, workoutBuf)
	return &Workout{ID: strconv.FormatInt(workoutID, 10), Path: fullPath, Ext: strings.ToLower(e.workoutsFormat)}, nil
}
