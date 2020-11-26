package export

import (
	"github.com/michalq/endo2strava/internal/modules/users"
	"github.com/michalq/endo2strava/pkg/endomondo-client"
)

// Orchestrator orchestrate export
type Orchestrator struct {
	endomondoExporter *Exporter
	usersRepository   users.Users
}

// NewOrchestrator creates instance of orchestrator
func NewOrchestrator(endomondoExporter *Exporter, usersRepository users.Users) *Orchestrator {
	return &Orchestrator{endomondoExporter, usersRepository}
}

// Run run specific stages of export
func (o *Orchestrator) Run(authorizedEndomondoClient *endomondo.Client, user *users.User, format string) error {

	workoutsQuantity, err := o.endomondoExporter.FindWorkoutsQuantity(authorizedEndomondoClient, user)
	if err != nil {
		return err
	}
	user.WorkoutsQuantity = workoutsQuantity
	if err := o.usersRepository.Update(user); err != nil {
		return err
	}
	err = o.endomondoExporter.FindWorkouts(authorizedEndomondoClient, user.WorkoutsQuantity)
	if err != nil {
		return err
	}
	err = o.endomondoExporter.FindWorkoutsDetails(authorizedEndomondoClient)
	if err != nil {
		return err
	}
	err = o.endomondoExporter.DownloadWorkouts(authorizedEndomondoClient, format)
	if err != nil {
		return err
	}
	return nil
}
