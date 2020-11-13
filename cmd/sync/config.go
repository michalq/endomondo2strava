package main

import "github.com/michalq/endo2strava/pkg/synchronizer"

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
