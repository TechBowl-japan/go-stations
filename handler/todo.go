package handler

import (
	"context"
	"encoding/json"
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
	if r.Method == http.MethodPost {
		todoRequest := &model.CreateTODORequest{}
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			log.Println(err)
		}

		if todoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todo := &model.TODO{}
		todo, err = h.svc.CreateTODO(r.Context(), todoRequest.Subject, todoRequest.Description)
		if err != nil {
			log.Println(err)
		}

		err = json.NewEncoder(w).Encode(&model.CreateTODOResponse{TODO: *todo})
		if err != nil {
			log.Println(err)
		}
	}

	if r.Method == http.MethodPut {
		todoRequest := &model.UpdateTODORequest{}
		err := json.NewDecoder(r.Body).Decode(&todoRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err)
		}

		if todoRequest.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if todoRequest.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		todo := &model.TODO{}
		todo, err = h.svc.UpdateTODO(r.Context(), todoRequest.ID, todoRequest.Subject, todoRequest.Description)
		if err != nil {
			log.Println(err)
		}

		err = json.NewEncoder(w).Encode(&model.CreateTODOResponse{TODO: *todo})
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			log.Println(err)
		}
	}

	if r.Method == http.MethodGet {
		todoRequest := &model.ReadTODORequest{PrevID: 0, Size: 5}
		if r.URL.Query().Get("prev_id") != "" {
			prevID, err := strconv.Atoi(r.URL.Query().Get("prev_id"))
			if err != nil {
				log.Println(err)
			}
			todoRequest.PrevID = int64(prevID)
		}

		if r.URL.Query().Get("size") != "" {
			size, err := strconv.Atoi(r.URL.Query().Get("size"))
			if err != nil {
				log.Println(err)
			}
			todoRequest.Size = int64(size)
		}

		todos, err := h.svc.ReadTODO(r.Context(), todoRequest.PrevID, todoRequest.Size)
		if err != nil {
			log.Println(err)
		}

		convertedTodos := make([]model.TODO, len(todos))
		for i, todo := range todos {
			convertedTodos[i] = *todo
		}
		err = json.NewEncoder(w).Encode(&model.ReadTODOResponse{TODOs: convertedTodos})
		if err != nil {
			log.Println(err)
		}
	}
}

// Create handles the endpoint that creates the TODO.
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
