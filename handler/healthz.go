package handler

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var  reqHealthz model.HealthzResponse
	reqHealthz.Message = "OK"
	

	if err := json.NewEncoder(w).Encode(&reqHealthz); err != nil {
		log.Println(err)
	}
}
