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
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodPut:
		h.handlePut(w, r)
	default:
		http.Error(w, "Unsupported request method", http.StatusMethodNotAllowed)
	}
}

func (h *TODOHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// クエリパラメータのprev_idとsizeを取得
	prev_id, _ := strconv.ParseInt(query.Get("prev_id"), 10, 64)
	size, _ := strconv.ParseInt(query.Get("size"), 10, 64)

	// sizeとprev_idの取得
	req := &model.ReadTODORequest{
		PrevID: prev_id,
		Size:   size,
	}

	// ReadTODOメソッドを呼び出してDBからTODOを取得
	res, err := h.Read(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to fetch TODOs", http.StatusInternalServerError)
		return
	}

	// JSONとしてHTTPレスポンスを返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TODOHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into CreateTODORequest
	req := &model.CreateTODORequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// subjectが空文字列かどうかをチェック
	if req.Subject == "" {
		http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
		return
	}

	// Call the Create method
	res, err := h.Create(r.Context(), req)
	if err != nil {
		http.Error(w, "Failed to create TODO", http.StatusInternalServerError)
		return
	}

	// Send the response back as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *TODOHandler) handlePut(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into UpdateTODORequest
	req := &model.UpdateTODORequest{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// idが0かどうかをチェック
	if req.ID == 0 {
		http.Error(w, "ID cannot be 0", http.StatusBadRequest)
		return
	}
	// subjectが空文字列かどうかをチェック
	if req.Subject == "" {
		http.Error(w, "Subject cannot be empty", http.StatusBadRequest)
		return
	}

	// Call the Update method
	res, err := h.Update(r.Context(), req)
	if err != nil {
		// エラーが ErrNotFound の場合、404 Not Found を返す
		if _, ok := err.(*model.ErrNotFound); ok {
			http.Error(w, "TODO not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update TODO", http.StatusInternalServerError)
		}
		return
	}

	// Send the response back as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(res); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	res := &model.CreateTODOResponse{
		TODO: *todo,
	}
	return res, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}

	res := &model.ReadTODOResponse{TODOs: make([]model.TODO, 0)}
	for _, todo := range todos {
		res.TODOs = append(res.TODOs, *todo)
	}

	return res, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	res := &model.UpdateTODOResponse{
		TODO: *todo,
	}
	return res, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
