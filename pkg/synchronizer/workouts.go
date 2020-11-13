package synchronizer

// Workouts is an repository for workouts export/import data
type Workouts interface {
	// SaveAll save all workouts
	SaveAll(workouts []Workout) (duplicated int, added int, err error)
	// Save persist single workout
	Save(workout *Workout) error
	// Update updates single workout
	Update(workout *Workout) error
	// FindAll retrieve all workouts
	FindAll() ([]Workouts, error)
	// FindOneByEndomondoID finds one workout by endomondo id
	FindOneByEndomondoID(endomondoID string) (*Workout, error)
}
