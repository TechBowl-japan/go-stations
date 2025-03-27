package handler

import (
	"context"
	"encoding/json"
	"errors"
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
	case http.MethodGet:
		var readReq model.ReadTODORequest
		const defaultPrevIDStr = "0"
		const defaultSizeStr = "5"
		prevIDStr := r.URL.Query().Get("prev_id")
		if prevIDStr == "" {
			prevIDStr = defaultPrevIDStr
		}
		sizeStr := r.URL.Query().Get("size")
		if sizeStr == "" {
			sizeStr = defaultSizeStr
		}

		if prevID, err := strconv.Atoi(prevIDStr); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			readReq.PrevID = int64(prevID)
		}

		if size, err := strconv.ParseInt(sizeStr, 10, 64); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			readReq.Size = int64(size)
		}

		readResponse, err := h.Read(r.Context(), &readReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(readResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var createReq model.CreateTODORequest

		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if createReq.Subject == "" {
			http.Error(w, "subject is empty", http.StatusBadRequest)
			return
		}

		todoResponse, err := h.Create(r.Context(), &createReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPut:
		var updateReq model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if updateReq.Subject == "" || updateReq.ID == 0 {
			http.Error(w, "subject or id is empty", http.StatusBadRequest)
			return
		}

		todoResponse, err := h.Update(r.Context(), &updateReq)
		if err != nil {
			// TODOが無い場合は404を返却
			var notFoundErr *model.ErrNotFound
			if errors.As(err, &notFoundErr) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		var deleteReq model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if len(deleteReq.IDs) == 0 {
			http.Error(w, "ids is empty", http.StatusBadRequest)
			return
		}
		todoResponse, err := h.Delete(r.Context(), &deleteReq)
		if err != nil {
			// idが存在しない場合は404を返却
			var notFoundErr *model.ErrNotFound
			if errors.As(err, &notFoundErr) {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(todoResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
