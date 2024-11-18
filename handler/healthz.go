package handler

import (
	"encoding/json"
	"log"
	"net/http"

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
	healthzResponse := &model.HealthzResponse{}
	healthzResponse.Message = "OK"
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(healthzResponse); err != nil {
		log.Println(err)
	}
}
