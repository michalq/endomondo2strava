package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michalq/endo2strava/pkg/dao"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
	"github.com/michalq/endo2strava/pkg/migration"
	"github.com/michalq/endo2strava/pkg/strava-client"
	"github.com/michalq/endo2strava/pkg/synchronizer"
)

const (
	// WorkoutsPath stores path where workouts will be downloaded
	WorkoutsPath = "./tmp/workouts"
)

var (
	ctx        = context.Background()
	httpClient = &http.Client{}
)

func main() {
	// Loading input
	fmt.Println("Hello world!")
	fmt.Println("---")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var steps synchronizer.SynchronizationSteps
	for _, step := range strings.Split(os.Getenv("STEP"), ",") {
		steps = append(steps, synchronizer.SynchronizationStep(step))
	}
	config := configuration{
		startAt:               os.Getenv("START_AT"),
		endAt:                 os.Getenv("END_AT"),
		endomondoEmail:        os.Getenv("ENDOMONDO_EMAIL"),
		endomondoPass:         os.Getenv("ENDOMONDO_PASS"),
		endomondoExportFormat: os.Getenv("ENDOMONDO_EXPORT_FORMAT"),
		stravaClientID:        os.Getenv("STRAVA_CLIENT_ID"),
		stravaClientSecret:    os.Getenv("STRAVA_CLIENT_SECRET"),
		step:                  steps,
	}

	// Loading deps
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	if _, err := endomondoClient.Authorize(config.endomondoEmail, config.endomondoPass); err != nil {
		log.Fatalf("Endomondo authorization failed (%s).\n", err)
	}
	endomondoDownloader := synchronizer.NewEndomondoDownloader(endomondoClient, WorkoutsPath, config.endomondoExportFormat, func(l string) { fmt.Println(l) })

	stravaClient := strava.NewClient(ctx, httpClient, "https://www.strava.com", config.stravaClientID, config.stravaClientSecret)
	db, err := sql.Open("sqlite3", "file:./tmp/db.sqlite")
	if err != nil {
		log.Fatalf("Database connection failed (%s).\n", err)
	}
	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migrations fail (%s).", err)
	}
	workoutsRepository := dao.NewWorkouts(db)
	// Validate input
	startTime, err := time.Parse(time.RFC3339, config.startAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Err", err)
	}
	endTime, err := time.Parse(time.RFC3339, config.endAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Err", err)
	}
	if config.endomondoExportFormat != string(endomondo.ExportFormatGPX) && config.endomondoExportFormat != string(endomondo.ExportFormatTCX) {
		log.Fatalf("Format not supported, supported format [%s, %s]", endomondo.ExportFormatTCX, endomondo.ExportFormatGPX)
	}

	fmt.Printf("Grant access to strava:\n%s\n\n...and copy code that will be after ?code= in redirected url:\n", stravaClient.GenerateAuthorizationURL())
	stravaCode := bufio.NewScanner(os.Stdin)
	stravaCode.Scan()
	fmt.Println("---")
	stravaClient, err = stravaClient.Authorize(stravaCode.Text())
	if err != nil {
		fmt.Println("Err", err)
	} else {
		fmt.Println("Strava authorized successfully!")
	}
	stravaUploader := synchronizer.NewStravaUploader(stravaClient)

	// Run
	fmt.Println("---")
	if config.step.Has(synchronizer.StepExport) {
		workouts := endomondoDownloader.DownloadAllBetween(startTime, endTime)
		if err := workoutsRepository.SaveAll(workouts); err != nil {
			fmt.Println("Error while saving workouts to db", err)
		}
	} else {
		fmt.Println("Skipping export")
	}

	if config.step.Has(synchronizer.StepImport) {
		workouts, err := workoutsRepository.FindAll()
		if err != nil {
			log.Fatalf("Error fetching workouts from db (%s)", err)
		}
		var toImport []synchronizer.Workout
		for _, workout := range workouts {
			if workout.UploadStarted == 0 {
				toImport = append(toImport, workout)
			}
		}
		uploaded, err := stravaUploader.UploadAll(toImport)
		if err != nil {
			fmt.Println(err)
		}
		for _, workout := range uploaded {
			if err := workoutsRepository.Update(&workout); err != nil {
				fmt.Println("Err", err)
			}
		}
		// TODO verify started import whether ended
		fmt.Printf("\n---\nSynchronized %d/%d workouts\n", len(uploaded)+(len(workouts)-len(toImport)), len(workouts))
	} else {
		fmt.Println("Skipping import")
	}
}
