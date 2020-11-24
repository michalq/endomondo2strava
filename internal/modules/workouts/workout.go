package workouts

// Workout represents single workout
type Workout struct {
	EndomondoID   string
	StravaID      string
	Path          string
	Ext           string
	UploadStarted int
	UploadEnded   int
}
