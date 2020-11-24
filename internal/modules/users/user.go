package users

// User represent object of user
type User struct {
	ID                    string
	StravaAccessToken     string
	StravaRefreshToken    string
	StravaAccessExpiresAt int64
}
