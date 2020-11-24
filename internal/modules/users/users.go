package users

// Users repository for users
type Users interface {
	Save(*User) error
	Update(*User) error
	FindOneByID(string) (*User, error)
}
