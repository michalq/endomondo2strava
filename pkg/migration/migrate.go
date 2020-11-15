package migration

import "database/sql"

var migrations []string = []string{
	`CREATE TABLE IF NOT EXISTS workouts (
		endomondo_id TEXT PRIMARY KEY, 
		strava_id TEXT, 
		path TEXT, 
		ext TEXT, 
		upload_started INTEGER NOT NULL DEFAULT 0, 
		upload_ended INTEGER NOT NULL DEFAULT 0
	)`,
	`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY, 
		strava_access_token TEXT, 
		strava_refresh_token TEXT, 
		strava_access_expires_at INTEGER
	)`,
}

// Migrate runs all migrations
func Migrate(db *sql.DB) error {
	for _, migration := range migrations {
		createWorkoutsTableStmt, _ := db.Prepare(migration)
		_, err := createWorkoutsTableStmt.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}
