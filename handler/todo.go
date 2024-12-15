package handler

import (
	"context"
	"encoding/json"
	"errors"
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

// ServeHTTP handles HTTP requests to the /todos endpoint.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// 既存の POST 処理
		req := model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}

		todo, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(model.CreateTODOResponse{TODO: *todo}); err != nil {
			log.Println(err)
			return
		}

	case http.MethodPut:
		// PUT 処理を追加
		//UpdateTODORequestにJSON Decodeを行う
		req := model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}

		//idが0の場合に400 BadRequestを返す処理
		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}

		//DBにあるTODOに変更を行い、対象がない場合はHTTP Responseを返す
		updatedTODO, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil {
			if errors.Is(err, &model.ErrNotFound{}) {
				http.Error(w, "", http.StatusNotFound)
				return
			}
		}

		//更新完了時、更新したTODOをUpdateTODOResponse に代入し、JSON Encodeを行ってHTTP Responsを返す
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(model.UpdateTODOResponse{TODO: *updatedTODO}); err != nil {
			log.Println(err)
			return
		}
	}
}

// handleCreate handles the creation of TODOs.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	_, _ = h.svc.CreateTODO(ctx, "", "")
	return &model.CreateTODOResponse{}, nil
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
