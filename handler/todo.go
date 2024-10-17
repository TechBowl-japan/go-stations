package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	SVC *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		SVC: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.SVC.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.SVC.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.SVC.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.SVC.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req model.CreateTODORequest

		// Decode JSON body into the CreateTODORequest struct
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
			return
		}

		// Check if the subject is not empty
		if req.Subject == "" {
			http.Error(w, "Bad Request: Subject cannot be empty", http.StatusBadRequest)
			return
		}

		// Call the service to create the TODO item
		todo, err := h.SVC.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Prepare the response
		response := &model.CreateTODOResponse{
			TODO: todo,
		}

		// Set content-type to application/json
		w.Header().Set("Content-Type", "application/json")

		// Encode the response as JSON
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
			return
		}
		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request: ID and Subject cannot be empty", http.StatusBadRequest)
			return
		}
		todo, err := h.SVC.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := &model.UpdateTODOResponse{
			TODO: todo,
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodGet {
		query := r.URL.Query()
		prevIDStr := query.Get("prev_id")
		sizeStr := query.Get("size")

		req := model.ReadTODORequest{}
		var err error

		if prevIDStr != "" {
			req.PrevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request: Invalid prev_id", http.StatusBadRequest)
				return
			}
		}
		if sizeStr != "" {
			req.Size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request: Invalid size", http.StatusBadRequest)
				return
			}
		} else {
			req.Size = 5
		}

		todos, err := h.SVC.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := &model.ReadTODOResponse{
			TODOs: todos,
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)

		if err := encoder.Encode(response); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

	} else if r.Method == http.MethodDelete {
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
			return
		}
		if len(req.IDs) == 0 {
			http.Error(w, "Bad Request: IDs cannot be empty", http.StatusBadRequest)
			return
		}

		err := h.SVC.DeleteTODO(r.Context(), req.IDs)
		if err != nil {
			if _, ok := err.(*model.ErrNotFound); ok {
				http.Error(w, "Not Found: TODOs not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		response := &model.DeleteTODOResponse{
			DeletedIDs: req.IDs,
		}
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(response); err != nil {
			http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
