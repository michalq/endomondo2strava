package dao

import (
	"database/sql"

	"github.com/michalq/endo2strava/pkg/synchronizer"
)

// Users repository for users
type Users struct {
	db *sql.DB
}

// NewUsers creates new instance of Users
func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

// Save insert or update user
func (u *Users) Save(user *synchronizer.User) error {
	if usr, _ := u.FindOneByID(user.ID); usr != nil {
		return u.Update(user)
	}
	stmt, _ := u.db.Prepare("INSERT INTO users (id, strava_access_token, strava_refresh_token, strava_access_expires_at) VALUES (?, ?, ?, ?)")
	_, err := stmt.Exec(user.ID, user.StravaAccessToken, user.StravaRefreshToken, user.StravaAccessExpiresAt)
	return err
}

// Update updates user
func (u *Users) Update(user *synchronizer.User) error {
	stmt, _ := u.db.Prepare("UPDATE users SET strava_access_token=?, strava_refresh_token=?, strava_access_expires_at=? WHERE id=?")
	_, err := stmt.Exec(user.StravaAccessToken, user.StravaRefreshToken, user.StravaAccessExpiresAt, user.ID)
	return err
}

// FindOneByID finds user by id
func (u *Users) FindOneByID(ID string) (*synchronizer.User, error) {
	user := &synchronizer.User{}
	err := u.db.
		QueryRow("SELECT id, strava_access_token, strava_refresh_token, strava_access_expires_at FROM users WHERE id=?", ID).
		Scan(&user.ID, &user.StravaAccessToken, &user.StravaRefreshToken, &user.StravaAccessExpiresAt)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return user, nil
}
