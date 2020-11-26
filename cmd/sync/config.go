package main

import (
	"github.com/michalq/endo2strava/internal/controllers"
)

type configuration struct {
	endomondoEmail        string
	endomondoPass         string
	endomondoExportFormat string
	stravaClientID        string
	stravaClientSecret    string
	action                controllers.SynchronizationActions
}
