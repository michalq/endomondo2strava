package export

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/michalq/endo2strava/internal/modules/workouts"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// WorkoutsPath where workouts files will be saved
type WorkoutsPath string

// Result represent single resoult of scrapping
type Result struct {
	Workouts []endomondo.WorkoutsResponseData
}

// Downloader finds workouts
type Downloader struct {
	workoutsPath WorkoutsPath
	logger       func(l string)
}

// NewDownloader creates new instance
func NewDownloader(
	workoutsPath WorkoutsPath,
	logger func(l string),
) *Downloader {
	return &Downloader{workoutsPath, logger}
}

// DownloadAll downloads all provided workouts
func (e *Downloader) DownloadAll(
	authorizedClient *endomondo.Client,
	format string,
	workoutsToDownload []workouts.Workout,
) ([]workouts.Workout, error) {
	var downloadedWorkouts []workouts.Workout

	workoutsChan := make(chan workouts.Workout)
	workoutErrorChan := make(chan error)
	for _, workout := range workoutsToDownload {
		go e.downloadSingleWorkout(authorizedClient, format, workout, workoutsChan, workoutErrorChan)
	}

	for range workoutsToDownload {
		select {
		case workout := <-workoutsChan:
			e.logger(fmt.Sprintf("Downloaded workout %s", workout.EndomondoID))
			downloadedWorkouts = append(downloadedWorkouts, workout)
		case err := <-workoutErrorChan:
			e.logger(fmt.Sprintln("Err", err))
		}
	}

	return downloadedWorkouts, nil
}

func (e *Downloader) downloadSingleWorkout(
	authorizedClient *endomondo.Client,
	format string,
	workout workouts.Workout,
	workoutChan chan<- workouts.Workout,
	errorChan chan<- error,
) {
	downloadedWorkout, err := e.downloadWorkout(authorizedClient, format, workout)
	if err != nil {
		errorChan <- err
	} else {
		workoutChan <- *downloadedWorkout
	}
}

func (e *Downloader) downloadWorkout(authorizedClient *endomondo.Client, format string, workout workouts.Workout) (*workouts.Workout, error) {
	workoutID, err := strconv.ParseInt(workout.EndomondoID, 10, 64)
	if err != nil {
		return nil, err
	}
	workoutBuf, err := authorizedClient.ExportWorkout(workoutID, format)
	defer workoutBuf.Close()
	if err != nil {
		return nil, err
	}
	fullPath := fmt.Sprintf("%s/%s.%s", e.workoutsPath, workout.EndomondoID, strings.ToLower(format))
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	_, err = io.Copy(out, workoutBuf)

	workout.Path = fullPath
	workout.Ext = strings.ToLower(format)
	return &workout, nil
}
