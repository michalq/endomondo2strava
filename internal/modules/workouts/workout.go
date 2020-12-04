package workouts

import (
	"strconv"
	"strings"
)

// Workout represents single workout
type Workout struct {
	EndomondoID string
	// StravaID upload ID, it is not workout ID
	StravaID string
	// Path where is saved workout file
	Path string
	// Ext extension of workout
	Ext         string
	Title       string
	Description string
	// Hashtags comma separated hashtags
	Hashtags string
	// Pictures comma separated pictures
	Pictures string
	// DetailsExported flag if workout details were exported
	DetailsExported int
	// UploadStarted flag if upload was started
	UploadStarted int
	// UploadEnded flag if upload was ended
	UploadEnded int

	// StravaActivityID real workout id that is filled after verification
	StravaActivityID string
	// StravaStatus text message from verification if success
	StravaStatus string
	// StravaError text message from verification if error
	StravaError string
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
