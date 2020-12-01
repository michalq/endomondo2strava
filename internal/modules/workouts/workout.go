package workouts

import (
	"strconv"
	"strings"
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

// HashtagsList returns hashtags as an array of string
func (w *Workout) HashtagsList() []string {
	if len(w.Hashtags) > 0 {
		return strings.Split(",", w.Hashtags)
	}
	return []string{}
}

// EndomondoIDAsInt id as int
func (w *Workout) EndomondoIDAsInt() int {
	id, _ := strconv.ParseInt(w.EndomondoID, 10, 32)
	return int(id)
}
