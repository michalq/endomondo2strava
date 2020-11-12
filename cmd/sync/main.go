package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
}

func main() {
	// Loading input
	fmt.Println("Hello world!")
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
	}

	// Loading deps
	ctx := context.Background()
	httpClient := &http.Client{}
	endomondoClient := endomondo.NewClient(ctx, httpClient, "https://www.endomondo.com")
	if _, err := endomondoClient.Authorize(in.endomondoEmail, in.endomondoPass); err != nil {
		log.Fatalln("Err", err)
	}
	endomondoDownloader := synchronizer.NewEndomondoDownloader(endomondoClient, WorkoutsPath, in.endomondoExportFormat)

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

	// Run
	endomondoDownloader.DownloadAllBetween(startTime, endTime)
}
