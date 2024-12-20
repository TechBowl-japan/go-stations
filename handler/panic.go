package handler

import (
	"net/http"
)

// A PanicHandler implements health check endpoint.
type PanicHandler struct{}

// NewPanicHandler returns PanicHandler based http.Handler.
func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Panic Happen!")
}
