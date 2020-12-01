package upload

import (
	"fmt"
)

// ErrorCollection collects errors
type ErrorCollection struct {
	collection []error
}

// NewErrorCollection creates new instance of ErrorCollection
func NewErrorCollection() *ErrorCollection {
	return &ErrorCollection{}
}

// Error returns
func (e *ErrorCollection) Error() string {
	return fmt.Sprintf("found %d errors", len(e.collection))
}

// All returns all collected errors
func (e *ErrorCollection) All() []error {
	return e.collection
}

// Append adds another error to collection
func (e *ErrorCollection) Append(err error) {
	e.collection = append(e.collection, err)
}
