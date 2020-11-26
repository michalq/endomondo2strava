package report

import "github.com/michalq/endo2strava/internal/modules/workouts"

// Generator generates report
type Generator struct {
	workoutsRepository workouts.Workouts
}

// NewGenerator creates instance of generator
func NewGenerator(workoutsRepository workouts.Workouts) *Generator {
	return &Generator{workoutsRepository}
}

// Generate calculates summary
func (g *Generator) Generate() (*Report, error) {
	workouts, err := g.workoutsRepository.FindAll()
	if err != nil {
		return nil, err
	}
	report := &Report{}
	for _, workout := range workouts {
		if workout.Path != "" {
			report.Downloaded++
		}
		if workout.Pictures != "" {
			report.FoundPhotos++
		}
		if workout.DetailsExported != 0 {
			report.FoundDetails++
		}
		if workout.StravaID != "" {
			report.Imported++
		}
	}
	report.FoundWorkouts = len(workouts)
	return report, nil
}
