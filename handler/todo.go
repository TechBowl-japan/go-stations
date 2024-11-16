package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var todoRequest model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if todoRequest.Subject == "" {
			http.Error(w, "subject should not be empty", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.CreateTODO(r.Context(), todoRequest.Subject, todoRequest.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		todoResponse := model.CreateTODOResponse{TODO: *todo}
		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			log.Println(err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "PUT" {
		var UpdateTODORequest model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&UpdateTODORequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if UpdateTODORequest.ID == 0 {
			http.Error(w, "Invalid ID. ID should not be 0", http.StatusBadRequest)
			return
		}
		if UpdateTODORequest.Subject == "" {
			http.Error(w, "Subject should not be empty", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.UpdateTODO(r.Context(), UpdateTODORequest.ID, UpdateTODORequest.Subject, UpdateTODORequest.Description)
		switch {
		case errors.Is(err, &model.ErrNotFound{}):
			http.NotFound(w, r)
			return
		case err != nil:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		todoResponse := &model.UpdateTODOResponse{TODO: *todo}
		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			http.Error(w, "Failed to encode json", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" {
		prevIDstr := r.URL.Query().Get("prev_id")
		sizeStr := r.URL.Query().Get("size")

		req := model.ReadTODORequest{}

		if prevIDstr != "" {
			prevID, err := strconv.ParseUint(prevIDstr, 10, 64)
			if err != nil {
				http.Error(w, "invalid prev_id parameter", http.StatusBadRequest)
				return
			}
			req.PrevID = int64(prevID)
		}

		if sizeStr != "" {
			size, err := strconv.ParseUint(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid size parameter", http.StatusBadRequest)
				return
			}
			req.Size = int64(size)
		} else {
			req.Size = int64(5)
		}

		todos, err := h.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		todoResponse := &model.ReadTODOResponse{TODOs: todos}
		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			http.Error(w, "Failed to encode json", http.StatusInternalServerError)
			return
		}
	} else if r.Method == "DELETE" {
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		if len(req.IDs) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		err := h.svc.DeleteTODO(r.Context(), req.IDs)
		if errors.Is(err, &model.ErrNotFound{}) {
			http.Error(w, "There is no corresponding data", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := &model.DeleteTODOResponse{}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode json", http.StatusInternalServerError)
			return
		}
	}

}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
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
