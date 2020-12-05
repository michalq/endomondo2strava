package controllers

import (
	"log"

	"github.com/michalq/endo2strava/internal/modules/report"
)

// ReportController manage report of export import
type ReportController struct {
	reportGenerator *report.Generator
}

// NewReportController creates instance of ReportController
func NewReportController(reportGenerator *report.Generator) *ReportController {
	return &ReportController{reportGenerator}
}

// ReportAction renders report
func (r *ReportController) ReportAction() {
	report, err := r.reportGenerator.Generate()
	if err != nil {
		log.Fatalln(err)
	}
	renderCliReport(report)
}
