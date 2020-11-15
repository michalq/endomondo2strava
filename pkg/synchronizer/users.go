package synchronizer

// User represent object of user
type User struct {
	ID                    string
	StravaAccessToken     string
	StravaRefreshToken    string
	StravaAccessExpiresAt int64
}

// Users repository for users
type Users interface {
	Save(*User) error
	Update(*User) error
	FindOneByID(string) (*User, error)
}
