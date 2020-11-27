package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/michalq/endo2strava/pkg/dao"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
	"github.com/michalq/endo2strava/pkg/migration"
	"github.com/michalq/endo2strava/pkg/strava-client"
	"github.com/michalq/endo2strava/pkg/synchronizer"
)

type configuration struct {
	startAt               string
	endAt                 string
	endomondoEmail        string
	endomondoPass         string
	endomondoExportFormat string
	stravaClientID        string
	stravaClientSecret    string
	step                  synchronizer.SynchronizationSteps
}

const (
	// WorkoutsPath stores path where workouts will be downloaded
	WorkoutsPath = "./tmp/workouts"
	// UserID in single context runtime this value doesn't matter, it is generated randomly
	UserID = "9b85e5d7-6a4a-4c07-82bd-67c2f7e920d5"
)

var (
	ctx        = context.Background()
	httpClient = &http.Client{}
	stravaCode = ""
)

func main() {
	// Loading input
	fmt.Println("Starting local server!")
	setupLocalServer()

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
	simpleLogger := func(l string) { fmt.Println(l) }
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	if _, err := endomondoClient.Authorize(config.endomondoEmail, config.endomondoPass); err != nil {
		log.Fatalf("Endomondo authorization failed (%s).\n", err)
	}
	endomondoDownloader := synchronizer.NewEndomondoDownloader(endomondoClient, WorkoutsPath, config.endomondoExportFormat, simpleLogger)

	stravaClient := strava.NewClient(ctx, httpClient, "https://www.strava.com", config.stravaClientID, config.stravaClientSecret)
	db, err := sql.Open("sqlite3", "file:./tmp/db.sqlite")
	if err != nil {
		log.Fatalf("Database connection failed (%s).\n", err)
	}
	if err := migration.Migrate(db); err != nil {
		log.Fatalf("Migrations fail (%s).", err)
	}
	workoutsRepository := dao.NewWorkouts(db)
	usersRepository := dao.NewUsers(db)

	// Validate input
	startTime, err := time.Parse(time.RFC3339, config.startAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Input error", err)
	}
	endTime, err := time.Parse(time.RFC3339, config.endAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Input error", err)
	}
	if config.endomondoExportFormat != string(endomondo.ExportFormatGPX) && config.endomondoExportFormat != string(endomondo.ExportFormatTCX) {
		log.Fatalf("Format not supported, supported format [%s, %s]", endomondo.ExportFormatTCX, endomondo.ExportFormatGPX)
	}

	// Find user for synchronization session
	user, err := usersRepository.FindOneByID(UserID)
	if err != nil {
		log.Fatalf("Couldn't find user (%s).", err)
	}
	if user == nil {
		user = &synchronizer.User{ID: UserID, StravaAccessExpiresAt: 0, StravaAccessToken: "", StravaRefreshToken: ""}
		usersRepository.Save(user)
	}

	authURL := stravaClient.GenerateAuthorizationURL()
	openBrowser(authURL)
	for i := 0; i < 12 && stravaCode == ""; i++ {
		log.Println("Waiting for user authorization in browser, sleeping 10 seconds")
		time.Sleep(time.Second * 10)
	}
	if stravaCode == "" {
		log.Fatalln("Timed out after waiting 2 minutes for auth code")
	}

	stravaClient, err = stravaClient.Authorize(stravaCode)
	if err != nil {
		log.Fatalf("Strava authorization fail (%s).", err)
	}

	user.StravaRefreshToken = stravaClient.Authorization().AccessToken
	user.StravaAccessToken = stravaClient.Authorization().RefreshToken
	user.StravaAccessExpiresAt = stravaClient.Authorization().ExpiresAt
	if err := usersRepository.Update(user); err != nil {
		log.Fatalf("Cannot update user (%s)", err)
	}
	stravaUploader := synchronizer.NewStravaUploader(stravaClient, workoutsRepository, simpleLogger)

	// Run
	fmt.Println("---")
	if config.step.Has(synchronizer.StepExport) {
		fmt.Println("Starting export")
		workouts := endomondoDownloader.DownloadAllBetween(startTime, endTime)
		if err := workoutsRepository.SaveAll(workouts); err != nil {
			fmt.Println("Error while saving workouts to db", err)
		}
	} else {
		fmt.Println("Skipping export")
	}

	if config.step.Has(synchronizer.StepImport) {
		fmt.Println("Starting import")
		err := stravaUploader.UploadAll()
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Skipping import")
	}
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		data, _ := ioutil.ReadFile("/proc/version")
		if data != nil && strings.Contains(string(data), "-microsoft-") {
			err = exec.Command("wslview", url).Start()
		} else {
			err = exec.Command("xdg-open", url).Start()
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func setupLocalServer() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/authorization", func(w http.ResponseWriter, r *http.Request) {
		val := r.FormValue("code")
		if val != "" {
			fmt.Fprintln(w, "Return to the console to watch progress of the transfer")
			stravaCode = val
		} else {
			fmt.Fprintf(w, "Missing auth code in response")
		}
	})

	go func() {
		fmt.Println("serving on 5000")
		err := http.ListenAndServe(":5000", r)
		if err != nil {
			panic("ListenAndServe: " + err.Error())
		}
	}()
}
