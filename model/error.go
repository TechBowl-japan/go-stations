package model


type ErrNotFound struct {
	ErrorNotFound error
}

func (e *ErrNotFound) Error() string {
	return e.ErrorNotFound.Error()
}