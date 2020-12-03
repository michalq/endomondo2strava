package controllers

import (
	"github.com/fatih/color"
)

// StdOutLogger logger that prints simple in std out
type StdOutLogger struct{}

// Info log information
func (StdOutLogger) Info(l string) {
	color.Blue(l)
}

// Warning log warning
func (StdOutLogger) Warning(l string) {
	color.Red(l)
}
