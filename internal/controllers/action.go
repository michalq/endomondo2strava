package controllers

// SynchronizationAction represents single step of synchronization
type SynchronizationAction string

// SynchronizationActions represent all synchronization steps to make
type SynchronizationActions []SynchronizationAction

// Has checks if given step is in stack
func (s SynchronizationActions) Has(idle SynchronizationAction) bool {
	for _, action := range s {
		if idle == action {
			return true
		}
	}
	return false
}

const (
	// ActionExport export
	ActionExport SynchronizationAction = "export"
	// ActionImport import
	ActionImport SynchronizationAction = "import"
	// ActionVerifyImport verification
	ActionVerifyImport SynchronizationAction = "verify"
	// ActionReport shows report
	ActionReport SynchronizationAction = "report"
)
