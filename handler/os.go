package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

type OSHandler struct{}

func NewOSHandler() *OSHandler {
	return &OSHandler{}
}

func (h *OSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	os, ok := r.Context().Value(model.OSContextKey("os")).(string)
	if !ok {
		os = "unknown"
	}

	res := &model.OSResponse{
		OS: os,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		log.Println(err)
	}
}
