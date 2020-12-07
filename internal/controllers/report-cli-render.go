package controllers

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/michalq/endo2strava/internal/modules/report"
)

func renderCliReport(summary *report.Report) {
	var foundDetails, foundPhotos string
	if summary.FoundWorkouts > summary.FoundDetails {
		foundDetails = color.RedString(strconv.Itoa(summary.FoundDetails))
		foundPhotos = color.RedString(strconv.Itoa(summary.FoundPhotos))
	} else {
		foundDetails = color.GreenString(strconv.Itoa(summary.FoundDetails))
		foundPhotos = color.GreenString(strconv.Itoa(summary.FoundPhotos))
	}
	var downloadedQuantity string
	if summary.FoundWorkouts > summary.Downloaded {
		downloadedQuantity = color.RedString(strconv.Itoa(summary.Downloaded))
	} else {
		downloadedQuantity = color.GreenString(strconv.Itoa(summary.Downloaded))
	}
	var importStartedQuantity string
	if summary.FoundWorkouts > summary.ImportStarted {
		importStartedQuantity = color.RedString(strconv.Itoa(summary.ImportStarted))
	} else {
		importStartedQuantity = color.GreenString(strconv.Itoa(summary.ImportStarted))
	}
	var importEndedQuantity string
	if summary.FoundWorkouts > summary.Imported {
		importEndedQuantity = color.RedString(strconv.Itoa(summary.Imported))
	} else {
		importEndedQuantity = color.GreenString(strconv.Itoa(summary.Imported))
	}
	var importVerifiedQuantity string
	if summary.FoundWorkouts > summary.Verified {
		importVerifiedQuantity = color.RedString(strconv.Itoa(summary.Verified))
	} else {
		importVerifiedQuantity = color.GreenString(strconv.Itoa(summary.Verified))
	}
	var importErrors string
	if summary.ImportErrors == 0 {
		importErrors = color.BlueString(strconv.Itoa(summary.ImportErrors))
	} else {
		importErrors = color.RedString(strconv.Itoa(summary.ImportErrors))
	}
	fmt.Printf(
		"----------------------------\n"+
			"Found workouts\t\t%s\n"+
			"----------------------------\n"+
			"Downloaded details\t%s\n"+
			"----------------------------\n"+
			"Found photos\t\t%s\n"+
			"----------------------------\n"+
			"Downloaded workouts\t%s\n"+
			"----------------------------\n"+
			"Sent to import\t\t%s\n"+
			"----------------------------\n"+
			"Imported\t\t%s\n"+
			"----------------------------\n"+
			"Verified\t\t%s\n"+
			"----------------------------\n"+
			"Import errors\t\t%s\n"+
			"----------------------------\n",
		color.GreenString(strconv.Itoa(summary.FoundWorkouts)),
		foundDetails,
		foundPhotos,
		downloadedQuantity,
		importStartedQuantity,
		importEndedQuantity,
		importVerifiedQuantity,
		importErrors,
	)
}
