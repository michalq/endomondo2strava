package dao

import (
	"database/sql"

	"github.com/michalq/endo2strava/internal/modules/workouts"
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
func (w *Workouts) SaveAll(workouts []workouts.Workout) error {
	for _, workout := range workouts {
		existing, err := w.FindOneByEndomondoID(workout.EndomondoID)
		if err != nil {
			return err
		}
		if existing != nil {
			if err := w.Update(&workout); err != nil {
				return err
			}
			continue
		}
		if err := w.Save(&workout); err != nil {
			return err
		}
	}
	return nil
}

// FindAll finds all workouts in db
func (w *Workouts) FindAll() ([]workouts.Workout, error) {
	rows, err := w.db.Query("SELECT endomondo_id, strava_id, path, ext, upload_started, upload_ended, title, description, hashtags, pictures, details_exported FROM workouts")
	if err != nil {
		return nil, err
	}

	allWorkouts := []workouts.Workout{}
	for rows.Next() {
		var workout workouts.Workout
		err = rows.Scan(&workout.EndomondoID, &workout.StravaID, &workout.Path, &workout.Ext, &workout.UploadStarted, &workout.UploadEnded)
		if err != nil {
			return nil, err
		}
		allWorkouts = append(allWorkouts, workout)
	}

	return allWorkouts, nil
}

// Save saves single workout in db
func (w *Workouts) Save(workout *workouts.Workout) error {
	stmt, _ := w.db.Prepare("INSERT INTO workouts (endomondo_id, strava_id, path, ext, upload_started, upload_ended, title, description, hashtags, pictures, details_exported) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	_, err := stmt.Exec(workout.EndomondoID, workout.StravaID, workout.Path, workout.Ext, workout.UploadStarted, workout.UploadEnded, workout.Title, workout.Description, workout.Hashtags, workout.Pictures, workout.DetailsExported)
	return err
}

// Update saves single workout in db
func (w *Workouts) Update(workout *workouts.Workout) error {
	stmt, _ := w.db.Prepare("UPDATE workouts SET strava_id = ?, path = ?, ext = ?, upload_started = ?, upload_ended = ?, title = ?, description = ?, hashtags = ?, pictures = ?, details_exported = ? WHERE endomondo_id = ?")
	_, err := stmt.Exec(workout.StravaID, workout.Path, workout.Ext, workout.UploadStarted, workout.UploadEnded, workout.Title, workout.Description, workout.Hashtags, workout.Pictures, workout.DetailsExported, workout.EndomondoID)
	return err
}

// FindOneByEndomondoID finds one workout
func (w *Workouts) FindOneByEndomondoID(endomondoID string) (*workouts.Workout, error) {
	workout := &workouts.Workout{}
	err := w.db.
		QueryRow("SELECT endomondo_id, strava_id, path, ext, upload_started, upload_ended, title, description, hashtags, pictures, details_exported FROM workouts WHERE endomondo_id=?", endomondoID).
		Scan(&workout.EndomondoID, &workout.StravaID, &workout.Path, &workout.Ext, &workout.UploadStarted, &workout.UploadEnded, &workout.Title, &workout.Description, &workout.Hashtags, &workout.Pictures, &workout.DetailsExported)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}
	return workout, nil
}
