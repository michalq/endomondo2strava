package strava

// GetUploadResponse response of get upload request
type GetUploadResponse struct {
	ID         string `json:"id_str"`
	ActivityID int    `json:"activity_id"`
	ExternalID string `json:"external_id"`
	Error      string `json:"error"`
	Status     string `json:"status"`
}
