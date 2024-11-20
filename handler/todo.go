package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements HTTP handler for TODO operations.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{svc: svc}
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Incoming Request:", req.Subject, req.Description)
		if req.Subject == "" {
			http.Error(w, "subject is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Created TODO:", todo)

		resp := &model.CreateTODOResponse{
			TODO: *todo,
		}
		log.Println("Response Data:", resp)

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println(w)
	} else if r.Method == http.MethodPut {
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("Incoming Update Request:", req.ID, req.Subject, req.Description)
		if req.ID == 0 {
			http.Error(w, "id is required", http.StatusBadRequest)
			return
		}
		if req.Subject == "" {
			http.Error(w, "subject is required", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		todo, err := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Updated TODO:", todo)
		resp := &model.UpdateTODOResponse{
			TODO: *todo,
		}
		log.Println("Response Data:", resp)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodGet {
		prevIDStr := r.URL.Query().Get("prev_id")
		sizeStr := r.URL.Query().Get("size")
		ctx := r.Context()
		var prevID int64
		var size int64

		var err error
		if prevIDStr == "" {
			prevID = 0
		} else {
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid prev_id", http.StatusBadRequest)
				return
			}
		}

		if sizeStr == "" {
			size = 0
		} else {
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid prev_id", http.StatusBadRequest)
				return
			}
		}

		todos, err := h.svc.ReadTODO(ctx, prevID, size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todosList := make([]model.TODO, len(todos))
		for i, todo := range todos {
			todosList[i] = *todo
		}

		log.Print("Read TODOs:", todosList)

		resp := &model.ReadTODOResponse{
			TODOs: todosList,
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
