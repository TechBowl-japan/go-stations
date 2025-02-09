package handler

import (
	"context"
	"encoding/json"
	"errors"
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

	case http.MethodGet:
		//クエリパラメータを取得、TODOの構造体に沿うように構築
		prevIDInt := r.URL.Query().Get("prev_id")
		sizeInt := r.URL.Query().Get("size")

		//クエリパラメータの型をint64に変換
		//prevIDが0でない場合
		var req model.ReadTODORequest
		if prevIDInt != "" {
			prevID, err := strconv.ParseInt(prevIDInt, 10, 64)
			if err == nil {
				req.PrevID = prevID
			}
		}
		//sizeが0でない場合
		if sizeInt != "" {
			size, err := strconv.ParseInt(sizeInt, 10, 64)
			if err == nil {
				req.Size = size
			}
		}
		//両方が０である場合
		if prevIDInt == "" && sizeInt == "" {
			req.PrevID = 0
			req.Size = -1
		}

		//DBにあるTODOを取得
		todo, err := h.svc.ReadTODO(r.Context(), req.PrevID, req.Size)
		if err != nil {
			log.Println(err)
			return
		}
		//responseに代入
		response := model.ReadTODOResponse{TODOs: todo}

		//JSON Encodeを行なって、HTTP ResPonseに返す
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Println(err)
			return
		}

	case http.MethodDelete:
		//削除リクエストにデコードかける
		var req model.DeleteTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}
		//idのリストが0だった場合、404を返す。lenはスライス、配列、文字列、マップ、チャンネルの長さを取得する組み込み関数である
		if len(req.IDs) == 0 {
			http.Error(w, "BadRequest", http.StatusBadRequest)
			return
		}

		//DBの中にあるTODOを取得し、対象のTODOがなかった場合はErrNotFoundが返され、その場合は404を返す
		err := h.svc.DeleteTODO(r.Context(), req.IDs)
		if err != nil {
			if errors.Is(err, &model.ErrNotFound{}) {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
			//ここでも404を返しているが、実際は場合によって200や500を設定する
			http.Error(w, "", http.StatusNotFound)
			return
		}

		//削除処理が無事成功した場合は、削除レスポンスを作成し、返信を行う。
		response := model.DeleteTODOResponse{}
		w.Header().Set("Content-Tyoe", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
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
