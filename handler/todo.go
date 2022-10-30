package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
// func NewTODOHandler(svc *service.TODOService) *TODOHandler {
// 	return &TODOHandler{
// 		svc: svc,
// 	}
// }

func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
// func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	//	fmt.Fprint(nil, "Hello TODO!")
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	// return &model.CreateTODOResponse{TODO: *todo}, nil
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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Hello TODO! Serve HTTP")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "decode error")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
			return
		}
		fmt.Fprintln(w, err, resp)
	}

	if r.Method == "GET" {
		fmt.Fprintln(w, "HTTP Request is GET")
		http.Error(w, "", http.StatusOK)
	}
	//h.Create(nil, nil)
}
