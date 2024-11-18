package model

type ErrNotFound struct {
	Resource string `json:"resource"`
}

func (e *ErrNotFound) Error() string {
	return e.Resource + " not found"
}
