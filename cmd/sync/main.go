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
	"github.com/michalq/endo2strava/internal/dal"
	"github.com/michalq/endo2strava/internal/migration"
	"github.com/michalq/endo2strava/internal/modules/export"
	"github.com/michalq/endo2strava/internal/modules/report"
	"github.com/michalq/endo2strava/internal/modules/upload"
	"github.com/michalq/endo2strava/internal/modules/users"

	"github.com/michalq/endo2strava/pkg/endomondo-client"
	"github.com/michalq/endo2strava/pkg/strava-client"
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

// StdOutLogger logger that prints simple in std out
type StdOutLogger struct{}

// Info log information
func (StdOutLogger) Info(l string) {
	fmt.Println(l)
}

func main() {
	// Loading input
	fmt.Println("Endomondo 2 Strava synchronizer.\n")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var actions controllers.SynchronizationActions
	for _, action := range strings.Split(os.Getenv("STEP"), ",") {
		actions = append(actions, controllers.SynchronizationAction(action))
	}
	config := configuration{
		endomondoEmail:        os.Getenv("ENDOMONDO_EMAIL"),
		endomondoPass:         os.Getenv("ENDOMONDO_PASS"),
		endomondoExportFormat: os.Getenv("ENDOMONDO_EXPORT_FORMAT"),
		stravaClientID:        os.Getenv("STRAVA_CLIENT_ID"),
		stravaClientSecret:    os.Getenv("STRAVA_CLIENT_SECRET"),
		action:                actions,
	}

	// Loading deps
	// logger := &StdOutLogger{}
	simpleLogger := func(l string) { fmt.Println(l) }
	db, err := sql.Open("sqlite3", "file:./tmp/db.sqlite")
	if err != nil {
		log.Fatalf("Database connection failed (%s).\n", err)
	}
	workoutsRepository := dal.NewWorkouts(db)
	usersRepository := dal.NewUsers(db)
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	stravaClient := strava.NewClient(ctx, httpClient, "https://www.strava.com", config.stravaClientID, config.stravaClientSecret)
	endomondoExporter := export.NewExporter(export.NewDownloader(filesPath, simpleLogger), workoutsRepository, usersRepository, simpleLogger)
	stravaImporter := upload.NewStravaUploader(workoutsRepository, simpleLogger)
	usersManager := users.NewManager(usersRepository)
	exportOrchestrator := export.NewOrchestrator(endomondoExporter, usersRepository)
	reportGenerator := report.NewGenerator(workoutsRepository)
	exportController := controllers.NewExportController(endomondoClient, usersManager, reportGenerator, exportOrchestrator)
	importController := controllers.NewImportController(stravaImporter, stravaClient, usersManager, usersRepository)

	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migrations fail (%s).", err)
	}

	// Validate input
	if config.endomondoExportFormat != string(endomondo.ExportFormatGPX) && config.endomondoExportFormat != string(endomondo.ExportFormatTCX) {
		log.Fatalf("Format not supported, supported format [%s, %s]", endomondo.ExportFormatTCX, endomondo.ExportFormatGPX)
	}

	if config.action.Has(controllers.ActionExport) {
		exportController.ExportAction(controllers.ExportInput{
			Email: config.endomondoEmail, Pass: config.endomondoPass, Format: config.endomondoExportFormat,
		})
	}

	if config.action.Has(controllers.ActionImport) {
		importController.ImportAction(controllers.ImportInput{
			ClientID: config.stravaClientID, ClientSecret: config.stravaClientSecret,
		})
	}
}
