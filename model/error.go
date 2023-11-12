package model

import (
	"fmt"
)

type ErrNotFound struct {
	todo_id int64
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("TODO is not found with id: %d", e.todo_id)
}

func NewErrNotFound(todo_id int64) error {
	return &ErrNotFound{todo_id}
}
