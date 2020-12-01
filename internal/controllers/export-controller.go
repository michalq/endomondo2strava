package controllers

import (
	"fmt"
	"log"

	"github.com/michalq/endo2strava/internal/modules/report"

	"github.com/michalq/endo2strava/internal/modules/users"

	"github.com/michalq/endo2strava/internal/modules/export"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// ExportController controller for export
type ExportController struct {
	endomondoClient *endomondo.Client
	userManager     *users.Manager
	reportGenerator *report.Generator
	orchestrator    *export.Orchestrator
}

// NewExportController creates new export controller instance
func NewExportController(endomondoClient *endomondo.Client, userManager *users.Manager, reportGenerator *report.Generator, orchestrator *export.Orchestrator) *ExportController {
	return &ExportController{endomondoClient, userManager, reportGenerator, orchestrator}
}

// ExportInput input for controller
type ExportInput struct {
	UserID string
	Email  string
	Pass   string
	Format string
}

// ExportAction handles exporting workouts from endomondo
func (e *ExportController) ExportAction(input ExportInput) {

	fmt.Println("Running export 🚀")
	user, err := e.userManager.FindOrCreate(input.UserID)
	if err != nil {
		log.Fatalln(err)
	}
	authorizedClient, err := e.endomondoClient.Authorize(input.Email, input.Pass)
	if err != nil {
		log.Fatal(err)
	}
	if err := e.orchestrator.Run(authorizedClient, user, input.Format); err != nil {
		log.Fatalln(err)
	}

	report, err := e.reportGenerator.Generate()
	renderCliReport(report)
	if err != nil {
		log.Fatalln(err)
	}
}
