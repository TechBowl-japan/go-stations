package model

type ErrNotFound struct {
}

func (e *ErrNotFound) Error() string {
	return "Not Found in DB"
}
