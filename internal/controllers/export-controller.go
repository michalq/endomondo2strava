package controllers

import (
	"log"

	"github.com/michalq/endo2strava/internal/modules/export"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// ExportInput input for controller
type ExportInput struct {
	Email  string
	Pass   string
	Format string
}

// ExportController handles exporting workouts from endomondo
func ExportController(exportInput ExportInput, endomondoExporter *export.Exporter, endomondoClient *endomondo.Client) {
	endomondoClient, err := endomondoClient.Authorize(exportInput.Email, exportInput.Pass)
	if err != nil {
		log.Fatal(err)
	}
	if err := endomondoExporter.RetrieveWorkouts(endomondoClient, exportInput.Format); err != nil {
		log.Fatal(err)
	}
}
