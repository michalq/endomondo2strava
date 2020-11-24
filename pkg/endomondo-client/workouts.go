package endomondo

// ExportFormat represents format of exporting data
type ExportFormat string

const (
	// ExportFormatTCX Garmin format with more data
	ExportFormatTCX ExportFormat = "TCX"
	// ExportFormatGPX Basic format
	ExportFormatGPX ExportFormat = "GPX"
)

// WorkoutType type of workout
type WorkoutType int

// WorkoutsQueryParams query params passed to workouts endpoint
type WorkoutsQueryParams struct {
	Before string
	After  string
	Offset *int
	Limit  *int
}

// WorkoutsResponse represents response payload of all workouts
type WorkoutsResponse struct {
	Data   []WorkoutsResponseData `json:"data"`
	Paging struct {
		Next     string `json:"next"`
		Total    int    `json:"total"`
		Previous string `json:"previous"`
	} `json:"paging"`
}

// WorkoutsResponseData single workout data
type WorkoutsResponseData struct {
	ID                   int64       `json:"id"`
	Sport                WorkoutType `json:"sport"`
	Expand               string      `json:"expand"`
	StartTime            string      `json:"start_time"`
	LocalStartTime       string      `json:"local_start_time"`
	Distance             float64     `json:"distance"`
	Duration             float64     `json:"duration"`
	SpeedAvg             float64     `json:"speed_avg"`
	SpeedMax             float64     `json:"speed_max"`
	AltitudeMin          float64     `json:"altitude_min"`
	AltitudeMax          float64     `json:"altitude_max"`
	Ascent               float64     `json:"ascent"`
	Descent              float64     `json:"descent"`
	PbCount              int         `json:"pb_count"`
	Calories             float64     `json:"calories"`
	IsLive               bool        `json:"is_live"`
	IncludeInStats       bool        `json:"include_in_stats"`
	CanFBShareViaBackend bool        `json:"can_fb_share_via_backend"`
}
