package handler

import (
	"fmt"
	"net/http"
)

// A HealthzHandler implements health check endpoint.
type PanicHandler struct {
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Start Panic")
	panic("a problem")
}
