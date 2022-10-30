package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct {
}

// NewHealthzHandler returns HealthzHandler based http.Handler.
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintln("ServeHTTP Starts")
	// w, err := json.MarShal("ok")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	_ = &model.HealthzResponse{}
	response := &model.HealthzResponse{
		Message: "OK",
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(response); err != nil {
		log.Println(err)
	}

}

func Handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello Handler!")
}
