package strava

// UploadParameters params needed by upload endpoint
type UploadParameters struct {

	// File The uploaded file.
	File string
	// Name The desired name of the resulting activity.
	Name string
	// Description The desired description of the resulting activity.
	Description string
	// Trainer Whether the resulting activity should be marked as having been performed on a trainer.
	Trainer string
	// Commute Whether the resulting activity should be tagged as a commute.
	Commute string
	// DataType The format of the uploaded file.
	DataType string
	// ExternalID The desired external identifier of the resulting activity.
	ExternalID string
}
