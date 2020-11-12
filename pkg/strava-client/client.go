package strava

// Client simple strava client
type Client struct {
	accessToken string
}

// NewClient creates instance of strava client
func NewClient(accessToken string) *Client {
	return &Client{accessToken}
}

// ImportWorkout send workout to strava
//
// Upload Activity (createUpload)
// Uploads a new data file to create an activity from. Requires activity:write scope.
// POST /uploads
//
// More: https://developers.strava.com/docs/reference/#api-models-Upload
func (c *Client) ImportWorkout(importID string) {}

// Authorize authorizes client
func (c *Client) Authorize() {}
