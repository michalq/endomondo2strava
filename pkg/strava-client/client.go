package strava

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
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
func (c *Client) ImportWorkout(upload UploadParameters) (*UploadResponse, error) {

	if c.auth == nil {
		return nil, errors.New("not authorized to strava")
	}
	file, err := os.Open(upload.File)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileContents)
	_ = writer.WriteField("description", upload.Description)
	_ = writer.WriteField("trainer", upload.Trainer)
	_ = writer.WriteField("commute", upload.Commute)
	_ = writer.WriteField("data_type", upload.DataType)
	_ = writer.WriteField("external_id", upload.ExternalID)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	apiURL := fmt.Sprintf("%s/api/v3/uploads", c.baseURL)
	req, _ := http.NewRequestWithContext(c.ctx, "POST", apiURL, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.auth.AccessToken))
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusCreated {
		uploadResponse := &UploadResponse{}
		if err := json.Unmarshal(respBody, uploadResponse); err != nil {
			return nil, err
		}
		return uploadResponse, nil
	}

	return nil, fmt.Errorf("unexpected response [%s)", string(respBody))
}

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

// AuthorizeDirectly authorize directly with prepared token data, skipping granting access
func (c *Client) AuthorizeDirectly(tokenData *AuthTokenData) (*Client, error) {

	expiresAt := time.Unix(tokenData.ExpiresAt, 0)
	if expiresAt.After(time.Now()) {
		return &Client{ctx: c.ctx, httpClient: c.httpClient, baseURL: c.baseURL, clientID: c.clientID, clientSecret: c.clientSecret, auth: tokenData}, nil
	}
	query := make(url.Values)
	query.Set("grant_type", "refresh_token")
	query.Set("client_id", c.clientID)
	query.Set("client_secret", c.clientSecret)
	query.Set("refresh_token", tokenData.RefreshToken)

	authTokenData, err := c.makeAuthRequest(query)
	if err != nil {
		return c, err
	}
	return &Client{ctx: c.ctx, httpClient: c.httpClient, baseURL: c.baseURL, clientID: c.clientID, clientSecret: c.clientSecret, auth: authTokenData}, nil
}

// Authorization returns authorization data
func (c *Client) Authorization() *AuthTokenData {
	return c.auth
}

// Authorize authorizes client
func (c *Client) Authorize(code string) (*Client, error) {

	query := make(url.Values)
	query.Set("grant_type", "authorization_code")
	query.Set("client_id", c.clientID)
	query.Set("client_secret", c.clientSecret)
	query.Set("code", code)
	authTokenData, err := c.makeAuthRequest(query)
	if err != nil {
		return c, err
	}
	return &Client{ctx: c.ctx, httpClient: c.httpClient, baseURL: c.baseURL, clientID: c.clientID, clientSecret: c.clientSecret, auth: authTokenData}, nil
}

func (c *Client) makeAuthRequest(query url.Values) (*AuthTokenData, error) {
	apiURL := fmt.Sprintf("%s/api/v3/oauth/token", c.baseURL)
	req, _ := http.NewRequestWithContext(c.ctx, "POST", apiURL, strings.NewReader(query.Encode()))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("api returned unexpected status code %d (%s)", resp.StatusCode, string(body))
	}
	authTokenData := &AuthTokenData{}
	if err := json.Unmarshal(body, authTokenData); err != nil {
		return nil, err
	}
	return authTokenData, nil
}
