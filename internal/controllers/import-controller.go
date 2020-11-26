package controllers

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/michalq/endo2strava/internal/modules/upload"
	"github.com/michalq/endo2strava/internal/modules/users"
	"github.com/michalq/endo2strava/pkg/strava-client"
)

// ImportController controller for import
type ImportController struct {
	stravaUploader  *upload.StravaUploader
	stravaClient    *strava.Client
	userManager     *users.Manager
	usersRepository users.Users
}

// NewImportController creates new import controller instance
func NewImportController(
	stravaUploader *upload.StravaUploader,
	stravaClient *strava.Client,
	userManager *users.Manager,
	usersRepository users.Users,
) *ImportController {
	return &ImportController{stravaUploader, stravaClient, userManager, usersRepository}
}

// ImportInput input passed to import controller
type ImportInput struct {
	UserID       string
	ClientID     string
	ClientSecret string
}

// ImportAction handles import data to strava
func (i *ImportController) ImportAction(input ImportInput) {
	user, err := i.userManager.FindOrCreate(input.UserID)
	if err != nil {
		log.Fatalln(err)
	}
	var authorizedClient *strava.Client
	if user.StravaAccessToken != "" {
		authorizedClient, err = i.stravaClient.AuthorizeDirectly(&strava.AuthTokenData{
			AccessToken:  user.StravaAccessToken,
			RefreshToken: user.StravaRefreshToken,
			ExpiresAt:    user.StravaAccessExpiresAt,
		})
		if err != nil {
			log.Fatalf("Strava authorization fail (%s).", err)
		}
	} else {
		fmt.Printf("Grant access to strava:\n%s\n\n...and copy code that will be after ?code= in redirected url:\n", i.stravaClient.GenerateAuthorizationURL())
		stravaCode := bufio.NewScanner(os.Stdin)
		stravaCode.Scan()
		authorizedClient, err = i.stravaClient.Authorize(stravaCode.Text())
		if err != nil {
			log.Fatalf("Strava authorization fail (%s).", err)
		}
	}
	user.StravaRefreshToken = authorizedClient.Authorization().AccessToken
	user.StravaAccessToken = authorizedClient.Authorization().RefreshToken
	user.StravaAccessExpiresAt = authorizedClient.Authorization().ExpiresAt
	if err := i.usersRepository.Update(user); err != nil {
		log.Fatalf("Cannot update user (%s)", err)
	}

	fmt.Println("Starting import")
	status, err := i.stravaUploader.UploadAll(authorizedClient)
	if err != nil {
		fmt.Println(err)
	}
	// TODO verify started import whether ended
	fmt.Printf("\n---\nUploaded: %d, Skipped: %d (due to pending or ended import), All: %d\n", status.Uploaded, status.Skipped, status.All)
}
