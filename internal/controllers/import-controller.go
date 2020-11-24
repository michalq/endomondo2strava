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

// ImportInput input passed to import controller
type ImportInput struct {
	UserID       string
	ClientID     string
	ClientSecret string
}

// ImportController handles import data to strava
func ImportController(input ImportInput, stravaUploader *upload.StravaUploader, stravaClient *strava.Client, usersRepository users.Users) {
	user, err := usersRepository.FindOneByID(input.UserID)
	if err != nil {
		log.Fatalf("Couldn't find user (%s).", err)
	}
	if user == nil {
		user = &users.User{ID: input.UserID, StravaAccessExpiresAt: 0, StravaAccessToken: "", StravaRefreshToken: ""}
		usersRepository.Save(user)
	}
	if user.StravaAccessToken != "" {
		stravaClient, err = stravaClient.AuthorizeDirectly(&strava.AuthTokenData{
			AccessToken:  user.StravaAccessToken,
			RefreshToken: user.StravaRefreshToken,
			ExpiresAt:    user.StravaAccessExpiresAt,
		})
		if err != nil {
			log.Fatalf("Strava authorization fail (%s).", err)
		}
	} else {
		fmt.Printf("Grant access to strava:\n%s\n\n...and copy code that will be after ?code= in redirected url:\n", stravaClient.GenerateAuthorizationURL())
		stravaCode := bufio.NewScanner(os.Stdin)
		stravaCode.Scan()
		stravaClient, err = stravaClient.Authorize(stravaCode.Text())
		if err != nil {
			log.Fatalf("Strava authorization fail (%s).", err)
		}
	}
	user.StravaRefreshToken = stravaClient.Authorization().AccessToken
	user.StravaAccessToken = stravaClient.Authorization().RefreshToken
	user.StravaAccessExpiresAt = stravaClient.Authorization().ExpiresAt
	if err := usersRepository.Update(user); err != nil {
		log.Fatalf("Cannot update user (%s)", err)
	}

	fmt.Println("Starting import")
	status, err := stravaUploader.UploadAll(stravaClient)
	if err != nil {
		fmt.Println(err)
	}
	// TODO verify started import whether ended
	fmt.Printf("\n---\nUploaded: %d, Skipped: %d (due to pending or ended import), All: %d\n", status.Uploaded, status.Skipped, status.All)
}
