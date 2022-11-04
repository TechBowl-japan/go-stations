package handler

import (
	"context"
	"encoding/json"
	"fmt"
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
	todo, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todo}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, _ := h.svc.UpdateTODO(ctx, int64(req.ID), req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, err
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

	if r.Method == "PUT" {
		var req model.UpdateTODORequest
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
		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "ID is INVALID")
			fmt.Fprintln(w, r)
			return
		}
		resp, err := h.Update(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if r.Method == "GET" {
		q := r.URL.Query()
		prevID, _ := strconv.ParseInt(q.Get("prev_id"), 10, 64)
		size, _ := strconv.ParseInt(q.Get("size"), 10, 64)
		req := model.ReadTODORequest{
			PrevID: prevID,
			Size:   size,
		}
		resp, err := h.Read(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, err)
			return
		}
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, err)
			return
		}
	}
	if r.Method == "DELETE" {
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "decode error")
			//http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// fmt.Fprintln(w, "Decoded")
		// fmt.Fprintln(w, req)
		// fmt.Fprintln(w, len(req.IDs))
		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, req)
			fmt.Fprintln(w, r.Body)
			return
		}
		resp, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w, "encode error")

			return
		}
		w.WriteHeader(http.StatusOK)
	}
	//h.Create(nil, nil)
}
