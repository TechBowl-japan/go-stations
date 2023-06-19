package handler

import (
	"context"
	"encoding/json"
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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// POST
	case http.MethodPost:
		var req *model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)

		// Subjectが空もしくは、デコード時にエラーがある場合
		if req.Subject == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			res, err := h.Create(r.Context(), req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				err = json.NewEncoder(w).Encode(res)
			}
		}

	// PUT
	case http.MethodPut:
		var req *model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&req)

		// IDが0、Subjectが空もしくは、デコード時にエラーがある場合
		if req.ID == 0 || req.Subject == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			res, err := h.Update(r.Context(), req)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				err = json.NewEncoder(w).Encode(res)
			}
		}
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	res, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *res}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *res}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
