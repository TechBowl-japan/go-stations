package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// ServeHTTP handles the HTTP request and dispatches to appropriate methods.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost:
		h.handleCreate(w, r)
    case r.Method == http.MethodGet:
        h.handleRead(w, r)
	case r.Method == http.MethodPut:
		h.handleUpdateByBody(w, r)
	case r.Method == http.MethodPatch:
		h.HandleUpdate(w, r)
	case r.Method == http.MethodDelete:
		h.handleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleCreate handles POST requests to create a new TODO.
func (h* TODOHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into CreateTODORequest
	var req model.CreateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Validate subject
	if req.Subject == "" {
		http.Error(w, "Subject is required", http.StatusBadRequest)
		return
	}

	// Call the service to create the TODO
	ctx := r.Context()
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		log.Println("Failed to create TODO:", err)
		http.Error(w, "Failed to create TODO", http.StatusInternalServerError)
		return
	}

	// Build the response
	resp := model.CreateTODOResponse{
		TODO: *todo,
	}

	// Encode and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// Read handles the endpoint that reads the TODOs.(非HTTP版, 使わないのでコメントアウト)
// func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
// 	_, _ = h.svc.ReadTODO(ctx, 0, 0)
// 	return &model.ReadTODOResponse{}, nil
// }

// handleRead handles GET /todos?prev_id=&size=
func (h *TODOHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()

	// prev_id
	var prevID int64
	if v := q.Get("prev_id"); v != "" {
		if x, err := strconv.ParseInt(v, 10, 64); err == nil && x >= 0 {
			prevID = x
		}
	}

	// size（0 の場合は“全部”取得したいので大きめに置き換える）
	var size int64
	if v := q.Get("size"); v != "" {
		if x, err := strconv.ParseInt(v, 10, 64); err == nil && x >= 0 {
			size = x
		}
	}
	effSize := size
	if effSize == 0 {
		effSize = 1<<31 - 1 // 十分大きい数（テストでは3件なのでこれで全件返る）
	}

	todosPtr, err := h.svc.ReadTODO(r.Context(), prevID, effSize)
	if err != nil {
		log.Println("failed to read todos:", err)
		http.Error(w, "failed to read todos", http.StatusInternalServerError)
		return
	}

	// model.ReadTODOResponse は []TODO なのでデリファレンスして詰め替える
	todos := make([]model.TODO, 0, len(todosPtr))
	for _, t := range todosPtr {
		if t != nil {
			todos = append(todos, *t)
		}
	}

	resp := model.ReadTODOResponse{TODOs: todos}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// handleUpdateByBody handles PUT /todos with JSON body {id, subject, description}.
func (h *TODOHandler) handleUpdateByBody(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.Header().Set("Allow", "PUT")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Sta13: id が 0、または subject が空文字列なら 400
	if req.ID == 0 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Subject) == "" {
		http.Error(w, "subject is required", http.StatusBadRequest)
		return
	}

	todo, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
	if err != nil {
		var nf *model.ErrNotFound
		if errors.As(err, &nf) {
			http.Error(w, nf.Error(), http.StatusNotFound)
			return
		}
		log.Println("failed to update todo:", err)
		http.Error(w, "failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(model.UpdateTODOResponse{TODO: *todo})
}

// handleUpdate handles PATCH/PUT requests to update an existing TODO.
func (h *TODOHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	// Accept PATCH or PUT
    if r.Method != http.MethodPatch && r.Method != http.MethodPut {
        w.Header().Set("Allow", "PATCH, PUT")
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

	// Extract id from path: /todos/{id}
	idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
	if idStr == "" || idStr == r.URL.Path {
		http.Error(w, "id is required in path (/todos/{id})", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "id must be a positive integer", http.StatusBadRequest)
		return
	}

	// Decode request body
	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Failed to decode request:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Call service
	ctx := r.Context()
	todo, err := h.svc.UpdateTODO(ctx, id, req.Subject, req.Description)
	if err != nil {
		var nf *model.ErrNotFound
		if errors.As(err, &nf) {
			http.Error(w, nf.Error(), http.StatusNotFound)
			return
		}
		log.Println("Failed to update TODO:", err)
		http.Error(w, "Failed to update TODO", http.StatusInternalServerError)
		return
	}

	// Build response
	resp := model.UpdateTODOResponse{
		TODO: *todo,
	}

	// Encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Failed to encode response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// handleDelete handles DELETE /todos with JSON body {"ids":[...]}.
func (h *TODOHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", "DELETE")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// JSON Body を Decode: {"ids":[1,2,3]}
	var req model.DeleteTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("failed to decode delete request:", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// id のリストが空なら 400 Bad Request（"Empty Ids" ケース）
	if len(req.IDs) == 0 {
		http.Error(w, "ids is required", http.StatusBadRequest)
		return
	}

	// Service 呼び出し
	if err := h.svc.DeleteTODO(r.Context(), req.IDs); err != nil {
		// Station 19: すべての TODO が存在しなかったとき ErrNotFound → 404
		var nf *model.ErrNotFound
		if errors.As(err, &nf) {
			http.Error(w, nf.Error(), http.StatusNotFound)
			return
		}

		log.Println("failed to delete todo:", err)
		http.Error(w, "failed to delete todo", http.StatusInternalServerError)
		return
	}

	// 成功したら DeleteTODOResponse を JSON で返す（中身は {}）
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(model.DeleteTODOResponse{}); err != nil {
		log.Println("failed to encode delete response:", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

