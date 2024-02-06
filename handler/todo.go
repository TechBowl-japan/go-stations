package handler

import (
	"context"
	"encoding/json"
	"fmt"
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

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todo, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todo}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{
		TODO: *todo,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, err
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Register the Content-Type as application/json
	w.Header().Set("Content-Type", "application/json")

	// Check the request methods
	// The order is CRUD related methods
	// POST -> Create
	// GET -> Read
	// PUT -> Update
	// DELETE -> Delete

	// CREATE
	if r.Method == "POST" {
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check whether the subject in the request is empty or not
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Subject is empty")
			fmt.Fprintln(w, r)
			return
		}

		resp, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(w, err)
			return
		}
		fmt.Fprintln(w, err, resp)
	}

	// READ
	if r.Method == "GET" {
		q := r.URL.Query()
		prevID, err := strconv.ParseInt(q.Get("prev_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		size, err := strconv.ParseInt(q.Get("size"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		req := model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
		}

		resp, err := h.Read(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}
	}

	// UPDATE
	if r.Method == "PUT" {
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Need to create error message for returning err
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "invalid request: subject is required", http.StatusBadRequest)
			return
		}
		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "invalid request: id is required", http.StatusBadRequest)
			fmt.Fprintln(w, r)
			return
		}
		resp, err := h.Update(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// DELETE
	if r.Method == "DELETE" {
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "invalid request: IDs length must be more than 0", http.StatusBadRequest)
			return
		}

		resp, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
