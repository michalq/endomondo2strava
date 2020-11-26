package users

import "fmt"

// Manager manages users
type Manager struct {
	usersRepository Users
}

// NewManager creates new instance of manager
func NewManager(usersRepository Users) *Manager {
	return &Manager{usersRepository}
}

// FindOrCreate finds user and if not exists creates one
func (m *Manager) FindOrCreate(id string) (*User, error) {
	user, err := m.usersRepository.FindOneByID(id)
	if err != nil {
		return nil, fmt.Errorf("could not find user (%s)", err)
	}
	if user == nil {
		user = &User{ID: id, StravaAccessExpiresAt: 0, StravaAccessToken: "", StravaRefreshToken: "", WorkoutsQuantity: 0}
		if err := m.usersRepository.Save(user); err != nil {
			return nil, err
		}
	}
	return user, nil
}
