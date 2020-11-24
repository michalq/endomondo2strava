package workouts

import (
	"strconv"
)

// Workout represents single workout
type Workout struct {
	EndomondoID string
	StravaID    string
	Path        string
	Ext         string
	Title       string
	Description string
	Hashtags    string
	Pictures    string
	// DetailsExported flag if workout details were exported
	DetailsExported int
	UploadStarted   int
	UploadEnded     int
}

// EndomondoIDAsInt id as int
func (w *Workout) EndomondoIDAsInt() int {
	id, _ := strconv.ParseInt(w.EndomondoID, 10, 32)
	return int(id)
}
