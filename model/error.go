package model

type ErrNotFound struct {
	Massage string
}

func (e ErrNotFound) Error() string {
	return e.Massage
}
