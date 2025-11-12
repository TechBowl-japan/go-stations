package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	switch {
	case r.Method == http.MethodPost:
		h.handleCreate(w, r)
	case r.Method == http.MethodPut:
		h.handleUpdateByBody(w, r)
	case r.Method == http.MethodPatch:
		h.HandleUpdate(w, r)
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

// handleUpdateByBody handles PUT /todos with JSON body {id, subject, description}.
func (h *TODOHandler) handleUpdateByBody(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", "PUT")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Sta13: id が 0、または subject が空文字列なら 400
	if req.ID == 0 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Subject) == "" {
		http.Error(w, "subject is required", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
	if err != nil {
		var nf *model.ErrNotFound
		if errors.As(err, &nf) {
			http.Error(w, nf.Error(), http.StatusNotFound)
			return
		}
		log.Println("failed to update todo:", err)
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.UpdateTODOResponse{TODO: *todo})
}

// handleUpdate handles PATCH/PUT requests to update an existing TODO.
func (h *TODOHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Accept PATCH or PUT
    if r.Method != http.MethodPatch && r.Method != http.MethodPut {
        w.Header().Set("Allow", "PATCH, PUT")
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// Extract id from path: /todos/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	if idStr == "" || idStr == r.URL.Path {
		http.Error(w, "id is required in path (/todos/{id})", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "id must be a positive integer", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call service
	ctx := r.Context()
	todo, err := h.svc.UpdateTODO(ctx, id, req.Subject, req.Description)
	if err != nil {
		var nf *model.ErrNotFound
		if errors.As(err, &nf) {
			http.Error(w, nf.Error(), http.StatusNotFound)
			return
		}
		log.Println("Failed to update TODO:", err)
		http.Error(w, "Failed to update TODO", http.StatusInternalServerError)
		return
	}

	// Build response
	resp := model.UpdateTODOResponse{
		TODO: *todo,
	}

	// Encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
