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

	// GET
	case http.MethodGet:
		var (
			req *model.ReadTODORequest = &model.ReadTODORequest{PrevID: 0, Size: 5}
			res *model.ReadTODOResponse
			err error
		)
		if r.URL.Query().Get("prev_id") != "" {
			req.PrevID, err = strconv.ParseInt(r.URL.Query().Get("prev_id"), 10, 64)
		}
		if r.URL.Query().Get("size") != "" {
			req.Size, err = strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
		}

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err = h.Read(r.Context(), req)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		err = json.NewEncoder(w).Encode(&res)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

	// DELETE
	case http.MethodDelete:
		var (
			req *model.DeleteTODORequest = &model.DeleteTODORequest{}
			res *model.DeleteTODOResponse
			err error
		)

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil || len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err = h.Delete(r.Context(), req)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			err = json.NewEncoder(w).Encode(&res)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
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
	res, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)

	return &model.ReadTODOResponse{TODOs: res}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *res}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)

	if err != nil {
		return nil, err
	} else {
		return &model.DeleteTODOResponse{}, nil
	}

}
