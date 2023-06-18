package model

type ErrNotFound struct{}

func (err ErrNotFound) Error() string {
	return "err"
}
