package model

import "fmt"

// ErrNotFound is an error struct for representing a "not found" error.
type ErrNotFound struct {
	Resource string
	ID int64
}

// Error implements the error interface for ErrNotFound.
// This provides a descriptive error message.
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}