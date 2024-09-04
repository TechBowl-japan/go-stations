package model

import "fmt"

type ErrNotFound struct {
	Resource string `json:"resource"`
	ID       int64  `json:"id"`
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s not found: id=%d", e.Resource, e.ID)
}
