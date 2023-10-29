package handler

import (
	"net/http"
)

type PanicHandler struct{}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("do-panic")
}
