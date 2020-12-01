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
	var importedQuantity string
	if summary.FoundWorkouts > summary.Imported {
		importedQuantity = color.RedString(strconv.Itoa(summary.Imported))
	} else {
		importedQuantity = color.GreenString(strconv.Itoa(summary.Imported))
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
			"Imported workouts\t%s\n"+
			"----------------------------\n",
		color.GreenString(strconv.Itoa(summary.FoundWorkouts)),
		foundDetails,
		foundPhotos,
		downloadedQuantity,
		importedQuantity,
	)
}
