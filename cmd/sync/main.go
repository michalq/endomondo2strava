package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/michalq/endo2strava/pkg/strava-client"

	"github.com/michalq/endo2strava/pkg/synchronizer"

	"github.com/joho/godotenv"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

const (
	// WorkoutsPath stores path where workouts will be downloaded
	WorkoutsPath = "./tmp/workouts"
)

type input struct {
	startAt               string
	endAt                 string
	endomondoEmail        string
	endomondoPass         string
	endomondoExportFormat string
	stravaClientID        string
	stravaClientSecret    string
}

func main() {
	// Loading input
	fmt.Println("Hello world!")
	fmt.Println("---")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	in := input{
		startAt:               os.Getenv("START_AT"),
		endAt:                 os.Getenv("END_AT"),
		endomondoEmail:        os.Getenv("ENDOMONDO_EMAIL"),
		endomondoPass:         os.Getenv("ENDOMONDO_PASS"),
		endomondoExportFormat: os.Getenv("ENDOMONDO_EXPORT_FORMAT"),
		stravaClientID:        os.Getenv("STRAVA_CLIENT_ID"),
		stravaClientSecret:    os.Getenv("STRAVA_CLIENT_SECRET"),
	}

	// Loading deps
	ctx := context.Background()
	httpClient := &http.Client{}
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	if _, err := endomondoClient.Authorize(in.endomondoEmail, in.endomondoPass); err != nil {
		log.Fatalln("Err", err)
	}
	stravaClient := strava.NewClient(ctx, httpClient, "https://www.strava.com", in.stravaClientID, in.stravaClientSecret)
	endomondoDownloader := synchronizer.NewEndomondoDownloader(endomondoClient, WorkoutsPath, in.endomondoExportFormat, func(l string) { fmt.Println(l) })
	stravaUploader := synchronizer.NewStravaUploader(stravaClient)

	// Validate input
	startTime, err := time.Parse(time.RFC3339, in.startAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Err", err)
	}
	endTime, err := time.Parse(time.RFC3339, in.endAt+"T00:00:00.000Z")
	if err != nil {
		log.Fatalln("Err", err)
	}
	if in.endomondoExportFormat != string(endomondo.ExportFormatGPX) && in.endomondoExportFormat != string(endomondo.ExportFormatTCX) {
		log.Fatalf("Format not supported, supported format [%s, %s]", endomondo.ExportFormatTCX, endomondo.ExportFormatGPX)
	}

	fmt.Printf("Grant access to strava:\n%s\n\n...and copy code that will be after ?code= in redirected url:\n", stravaClient.GenerateAuthorizationURL())
	stravaCode := bufio.NewScanner(os.Stdin)
	stravaCode.Scan()
	fmt.Println("---")
	if err := stravaClient.Authorize(stravaCode.Text()); err != nil {
		log.Fatalln("Err", err)
	} else {
		fmt.Println("Strava authorized successfully!")
	}

	// Run
	fmt.Println("---")
	workouts := endomondoDownloader.DownloadAllBetween(startTime, endTime)
	uploaded := stravaUploader.UploadAll(workouts)
	fmt.Printf("\n---\nSynchronized %d/%d workouts\n", uploaded, len(workouts))
}
