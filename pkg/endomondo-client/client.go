package endomondo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Client simple endomondo client
type Client struct {
	ctx        context.Context
	httpClient *http.Client
	baseURL    string
	userID     int64
	authToken  string
}

// NewClient creates new endomondo client
func NewClient(ctx context.Context, httpClient *http.Client, baseURL string) *Client {
	return &Client{ctx: ctx, httpClient: httpClient, baseURL: baseURL, userID: 0, authToken: ""}
}

// Authorize authorizes user
// It creates new authorized object of client.
func (c *Client) Authorize(email, pass string) (*Client, error) {

	var authToken string
	requestBody, _ := json.Marshal(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
	}{
		Email:    email,
		Password: pass,
		Remember: true,
	})
	resp, err := http.Post(fmt.Sprintf("%s/rest/session", c.baseURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	authResponse := &AuthResponse{}
	if err := json.Unmarshal(body, authResponse); err != nil {
		return nil, err
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "USER_TOKEN" {
			authToken = cookie.Value
			break
		}
	}
	return &Client{ctx: c.ctx, httpClient: c.httpClient, baseURL: c.baseURL, userID: authResponse.ID, authToken: authToken}, nil
}

// Workout find workout details
func (c *Client) Workout(workoutID int) (*WorkoutDetailsResponse, error) {

	response := &WorkoutDetailsResponse{}
	_, err := c.makeGETRequestAndReadBody(response, fmt.Sprintf("%s/rest/v1/users/%d/workouts/%d", c.baseURL, c.userID, workoutID))
	if err != nil {
		return nil, err
	}
	return response, nil
}

// Workouts finds all workouts
func (c *Client) Workouts(queryParams WorkoutsQueryParams) (*WorkoutsResponse, error) {

	query := make(url.Values)
	if queryParams.Offset != nil {
		query.Set("offset", strconv.Itoa(*queryParams.Offset))
	}
	if queryParams.Limit != nil {
		query.Set("limit", strconv.Itoa(*queryParams.Limit))
	}
	if queryParams.Before != "" {
		query.Set("before", queryParams.Before)
	}
	if queryParams.After != "" {
		query.Set("after", queryParams.After)
	}
	response := &WorkoutsResponse{}
	_, err := c.makeGETRequestAndReadBody(response, fmt.Sprintf("%s/rest/v1/users/%d/workouts/history?%s", c.baseURL, c.userID, query.Encode()))
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ExportWorkout export workout
func (c *Client) ExportWorkout(workoutID int64, format string) (io.ReadCloser, error) {

	query := make(url.Values)
	query.Set("format", format)
	resp, err := c.makeGETRequest(fmt.Sprintf("%s/rest/v1/users/%d/workouts/%d/export?%s", c.baseURL, c.userID, workoutID, query.Encode()))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("non 200 status code when downloading workout")
	}

	return resp.Body, err
}

func (c *Client) makeGETRequest(url string) (*http.Response, error) {
	req, _ := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	req.AddCookie(&http.Cookie{Name: "USER_TOKEN", Value: c.authToken})
	return c.httpClient.Do(req)
}

func (c *Client) makeGETRequestAndReadBody(responseBody interface{}, url string) (*http.Response, error) {
	resp, err := c.makeGETRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, responseBody); err != nil {
		return nil, err
	}
	return resp, nil
}
