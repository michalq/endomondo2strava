package endomondo

// WorkoutDetailsResponse represents single workouts details api response
// It doesn't represent full response, just few important fields to export like title, name, picture or hashtags.
type WorkoutDetailsResponse struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	// Message is simple description of workout
	Message  string   `json:"message"`
	Hashtags []string `json:"hashtags"`
	Pictures []string `json:"picture"`
}
