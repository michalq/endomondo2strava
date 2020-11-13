package strava

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// ScopeActivityWrite Access to create manual activities and uploads, and access to edit any activities that are visible to the app, based on activity read access level
const ScopeActivityWrite = "activity:write"

// Client simple strava client
type Client struct {
	ctx          context.Context
	httpClient   *http.Client
	baseURL      string
	clientID     string
	clientSecret string
	auth         *AuthTokenData
}

// NewClient creates instance of strava client
func NewClient(ctx context.Context, httpClient *http.Client, baseURL, clientID, clientSecret string) *Client {
	return &Client{ctx: ctx, httpClient: httpClient, baseURL: baseURL, clientID: clientID, clientSecret: clientSecret}
}

// ImportWorkout send workout to strava
//
// Upload Activity (createUpload)
// Uploads a new data file to create an activity from. Requires activity:write scope.
// POST /uploads
//
// More: https://developers.strava.com/docs/reference/#api-models-Upload
func (c *Client) ImportWorkout(importID string) {}

// GenerateAuthorizationURL generates url that user must accept access
func (c *Client) GenerateAuthorizationURL() string {

	query := make(url.Values)
	query.Set("client_id", c.clientID)
	query.Set("scope", ScopeActivityWrite)
	query.Set("response_type", "code")
	query.Set("approval_prompt ", "auto")
	query.Set("redirect_uri", "http://127.0.0.1:5000/authorization")
	return fmt.Sprintf("%s/oauth/authorize?%s", c.baseURL, query.Encode())
}

// Authorize authorizes client
func (c *Client) Authorize(code string) error {

	query := make(url.Values)
	query.Set("grant_type", "authorization_code")
	query.Set("client_id", c.clientID)
	query.Set("client_secret", c.clientSecret)
	query.Set("code", code)
	url := fmt.Sprintf("%s/api/v3/oauth/token", c.baseURL)
	req, _ := http.NewRequestWithContext(c.ctx, "POST", url, strings.NewReader(query.Encode()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("api returned unexpected status code %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	authTokenData := &AuthTokenData{}
	if err := json.Unmarshal(body, authTokenData); err != nil {
		return err
	}
	c.auth = authTokenData
	return err
}
