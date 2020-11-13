package synchronizer

// SynchronizationStep represents single step of synchronization
type SynchronizationStep string

// SynchronizationSteps represent all synchronization steps to make
type SynchronizationSteps []SynchronizationStep

// Has checks if given step is in stack
func (s SynchronizationSteps) Has(idle SynchronizationStep) bool {
	for _, step := range s {
		if idle == step {
			return true
		}
	}
	return false
}

const (
	// StepExport export
	StepExport SynchronizationStep = "export"
	// StepImport import
	StepImport SynchronizationStep = "import"
)
