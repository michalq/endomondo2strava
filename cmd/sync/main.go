package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michalq/endo2strava/internal/controllers"
	"github.com/michalq/endo2strava/internal/dao"
	"github.com/michalq/endo2strava/internal/migration"
	"github.com/michalq/endo2strava/internal/modules/export"
	"github.com/michalq/endo2strava/internal/modules/upload"

	"github.com/michalq/endo2strava/pkg/endomondo-client"
	"github.com/michalq/endo2strava/pkg/strava-client"
	"github.com/michalq/endo2strava/pkg/synchronizer"
)

const (
	filesPath = "./tmp/workouts"
	// UserID in single context runtime this value doesn't matter, it is generated randomly
	UserID = "9b85e5d7-6a4a-4c07-82bd-67c2f7e920d5"
)

var (
	ctx        = context.Background()
	httpClient = &http.Client{}
)

func main() {
	// Loading input
	fmt.Println("Hello world!")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var steps synchronizer.SynchronizationSteps
	for _, step := range strings.Split(os.Getenv("STEP"), ",") {
		steps = append(steps, synchronizer.SynchronizationStep(step))
	}
	config := configuration{
		endomondoEmail:        os.Getenv("ENDOMONDO_EMAIL"),
		endomondoPass:         os.Getenv("ENDOMONDO_PASS"),
		endomondoExportFormat: os.Getenv("ENDOMONDO_EXPORT_FORMAT"),
		stravaClientID:        os.Getenv("STRAVA_CLIENT_ID"),
		stravaClientSecret:    os.Getenv("STRAVA_CLIENT_SECRET"),
		step:                  steps,
	}

	// Loading deps
	simpleLogger := func(l string) { fmt.Println(l) }
	db, err := sql.Open("sqlite3", "file:./tmp/db.sqlite")
	if err != nil {
		log.Fatalf("Database connection failed (%s).\n", err)
	}
	workoutsRepository := dao.NewWorkouts(db)
	usersRepository := dao.NewUsers(db)
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	stravaClient := strava.NewClient(ctx, httpClient, "https://www.strava.com", config.stravaClientID, config.stravaClientSecret)
	endomondoExporter := export.NewExporter(export.NewDownloader(filesPath, simpleLogger), workoutsRepository)
	stravaImporter := upload.NewStravaUploader(workoutsRepository, simpleLogger)

	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migrations fail (%s).", err)
	}

	// Validate input
	if config.endomondoExportFormat != string(endomondo.ExportFormatGPX) && config.endomondoExportFormat != string(endomondo.ExportFormatTCX) {
		log.Fatalf("Format not supported, supported format [%s, %s]", endomondo.ExportFormatTCX, endomondo.ExportFormatGPX)
	}

	if config.step.Has(synchronizer.StepExport) {
		controllers.ExportController(controllers.ExportInput{
			Email: config.endomondoEmail, Pass: config.endomondoPass, Format: config.endomondoExportFormat,
		}, endomondoExporter, endomondoClient)
	} else {
		fmt.Println("Skipping export")
	}

	if config.step.Has(synchronizer.StepImport) {
		controllers.ImportController(controllers.ImportInput{
			ClientID: config.stravaClientID, ClientSecret: config.stravaClientSecret,
		}, stravaImporter, stravaClient, usersRepository)
	} else {
		fmt.Println("Skipping import")
	}
}
