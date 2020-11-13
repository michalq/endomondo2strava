package dao

import (
	"database/sql"

	"github.com/michalq/endo2strava/pkg/synchronizer"
)

// Workouts is an repository for workouts export/import data
type Workouts struct {
	db *sql.DB
}

// NewWorkouts creates instance of workouts repository
func NewWorkouts(db *sql.DB) *Workouts {
	return &Workouts{db}
}

// SaveAll save all workouts in database
func (w *Workouts) SaveAll(workouts []synchronizer.Workout) error {
	for _, workout := range workouts {
		if workout, err := w.FindOneByEndomondoID(workout.EndomondoID); err != nil {
			return err
		} else if workout != nil {
			continue
		}
		if err := w.Save(&workout); err != nil {
			return err
		}
	}
	return nil
}

// FindAll finds all workouts in db
func (w *Workouts) FindAll() ([]synchronizer.Workout, error) {
	rows, err := w.db.Query("SELECT endomondo_id, strava_id, path, ext, upload_started, upload_ended FROM workouts")
	if err != nil {
		return nil, err
	}

	workouts := []synchronizer.Workout{}
	for rows.Next() {
		var workout synchronizer.Workout
		err = rows.Scan(&workout.EndomondoID, &workout.StravaID, &workout.Path, &workout.Ext, &workout.UploadStarted, &workout.UploadEnded)
		if err != nil {
			return nil, err
		}
		workouts = append(workouts, workout)
	}

	return workouts, nil
}

// Save saves single workout in db
func (w *Workouts) Save(workout *synchronizer.Workout) error {
	stmt, _ := w.db.Prepare("INSERT INTO workouts (endomondo_id, strava_id, path, ext, upload_started, upload_ended) VALUES (?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(workout.EndomondoID, workout.StravaID, workout.Path, workout.Ext, workout.UploadStarted, workout.UploadEnded)
	return err
}

// Update saves single workout in db
func (w *Workouts) Update(workout *synchronizer.Workout) error {
	stmt, _ := w.db.Prepare("UPDATE workouts SET strava_id = ?, path = ?, ext = ?, upload_started = ?, upload_ended = ? WHERE endomondo_id = ?")
	_, err := stmt.Exec(workout.StravaID, workout.Path, workout.Ext, 0, 0, workout.EndomondoID)
	return err
}

// FindOneByEndomondoID finds one workout
func (w *Workouts) FindOneByEndomondoID(endomondoID string) (*synchronizer.Workout, error) {
	workout := &synchronizer.Workout{}
	err := w.db.
		QueryRow("SELECT endomondo_id, strava_id, path, ext, upload_started, upload_ended FROM workouts WHERE endomondo_id=?", endomondoID).
		Scan(&workout.EndomondoID, &workout.StravaID, &workout.Path, &workout.Ext, &workout.UploadStarted, &workout.UploadEnded)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return workout, nil
}
