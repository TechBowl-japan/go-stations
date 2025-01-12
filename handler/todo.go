package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP handles the HTTP request and dispatches to appropriate methods.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handleCreate(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCreate handles POST requests to create a new TODO.
func (h* TODOHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into CreateTODORequest
	var req model.CreateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate subject
	if req.Subject == "" {
		http.Error(w, "Subject is required", http.StatusBadRequest)
		return
	}

	// Call the service to create the TODO
	ctx := r.Context()
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		log.Println("Failed to create TODO:", err)
		http.Error(w, "Failed to create TODO", http.StatusInternalServerError)
		return
	}

	// Build the response
	resp := model.CreateTODOResponse{
		TODO: *todo,
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
