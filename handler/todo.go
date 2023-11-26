package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"log"
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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var reqCreateTodo model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&reqCreateTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if reqCreateTodo.Subject == "" {
			http.Error(w, "Subject was empty", http.StatusBadRequest)
			return
		}
		var resCreateTodo model.CreateTODOResponse
		todo, err := h.svc.CreateTODO(r.Context(), reqCreateTodo.Subject, reqCreateTodo.Description)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		resCreateTodo.TODO = *todo 
		w.WriteHeader(http.StatusOK)
		
		if err := json.NewEncoder(w).Encode(&resCreateTodo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
		
	} else if r.Method == http.MethodPut {
		var reqUpdateTodo model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&reqUpdateTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		if reqUpdateTodo.ID == 0 {
			http.Error(w, "ID was not propriate.", http.StatusBadRequest)
			return 
		}

		if reqUpdateTodo.Subject == "" {
			http.Error(w, "Subject was empty", http.StatusBadRequest)
			return
		}

		var resUpdateTodo model.UpdateTODOResponse
		todo, err := h.svc.UpdateTODO(r.Context(), reqUpdateTodo.ID, reqUpdateTodo.Subject, reqUpdateTodo.Description)
        
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

		if todo == nil {
			http.NotFound(w, r)
			return 
		}

		resUpdateTodo.TODO = *todo 
		w.WriteHeader(http.StatusOK)
		
		if err := json.NewEncoder(w).Encode(&resUpdateTodo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}

	} else if r.Method == http.MethodGet {
		var reqReadTodo model.ReadTODORequest
		prev_idStr := r.URL.Query().Get("prev_id")
		if prev_idStr == "" {
			prev_idStr = "0"
		}
		prev_id, err := strconv.Atoi(prev_idStr)
		if err != nil {
			log.Println(err)
			return 
		}
		reqReadTodo.PrevID = int64(prev_id)

        size_Str := r.URL.Query().Get("size")
		if size_Str == "" {
			size_Str = "0"
		}
		size, err := strconv.Atoi(size_Str)
		if err != nil {
			log.Println(err)
			return
		}
		reqReadTodo.Size = int64(size)

		todos, err := h.svc.ReadTODO(r.Context(), reqReadTodo.PrevID, reqReadTodo.Size)
		if err != nil {
			log.Println(err)
			return 
		}

		var resReadTodo model.ReadTODOResponse
		if todos == nil {
			resReadTodo.TODOs = []*model.TODO{} 
		} else {
			resReadTodo.TODOs = todos
		}
		
		if err := json.NewEncoder(w).Encode(&resReadTodo); err != nil {
			log.Printf("Error encoding response: %v", err)
			return 
		}

		
	} else if r.Method == http.MethodDelete {
		var reqDeleteTodo model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&reqDeleteTodo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		if len(reqDeleteTodo.IDs) == 0 {
			http.Error(w, "Id list was empty", http.StatusBadRequest)
			return 
		}
		err := h.svc.DeleteTODO(r.Context(), reqDeleteTodo.IDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return 
		}
		var resDeleteTodo model.DeleteTODOResponse
		if err := json.NewEncoder(w).Encode(&resDeleteTodo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return 
		}
	}
}
