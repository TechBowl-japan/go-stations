package model

import (
	"fmt"
)

type ErrNotFound struct {
	todo_id int
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("TODO is not found")
}