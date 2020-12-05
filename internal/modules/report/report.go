package report

// Report represents stats of process
type Report struct {
	FoundWorkouts int
	FoundDetails  int
	FoundPhotos   int
	Downloaded    int
	Imported      int
	ImportStarted int
	Verified      int
}
