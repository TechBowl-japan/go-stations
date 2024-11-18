package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
// TODOHandlerは、TODOに関するREST APIエンドポイントを処理を実装します。
type TODOHandler struct {
	svc *service.TODOService //TODOServiceを使用してデータ操作を行う
}

// NewTODOHandler returns TODOHandler based http.Handler.
// NewTODOHandlerは新しいTODOHandlerを返します。
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc, //TODOServiceを注入
	}
}

// ServeHTTP handles HTTP requests for the TODO API.
// リクエストのHTTPメソッドに基づいて適切なハンドラを呼び出します。
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost: //POSTメソッドの場合
		h.handleCreate(w, r) //TODO作成の処理を呼び出す
	default:
		//他のメソッドは許可されていないため、エラーレスポンスを返す
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// handleCreate handles the POST request to create a niw TODO.
// handleCreateは、新しいTODOを作成するためのPOSTリクエストを処理する。
func (h *TODOHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to CreateTODORequest
	//リクエストボディを解析し、CreateTODORequest構造体にデコードする。
	var req model.CreateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		//JSONのデコードに失敗した場合、400BadRequestを返す
		log.Printf("Error decoding CreateTODORequest: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() //リクエストボディをクローズする
	//必須フィールドであるSubjectが空でないかをチェックする
	if req.Subject == "" {
		//Subjectが空の場合、400BadRequestを返す
		http.Error(w, "Subject is required", http.StatusBadRequest)
		return
	}
	//Contextを取得し、Createメソッドを呼び出してTODOを作成する
	ctx := r.Context()
	res, err := h.Create(ctx, &req)
	if err != nil {
		//TODOの作成時にエラーが発生した場合、500Internal Server Errorを返す
		http.Error(w, "Failed to create TODO", http.StatusInternalServerError)
		return
	}
	//レスポンスヘッダを設定し、成功ステータス(200 OK)を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		//レスポンスヘッダのエンコードに失敗した場合、500 Internal Server Errorを返す
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}

}

// Create handles the endpoint that creates the TODO.
// TODOServiceのCreateTODOメソッドを呼び出し、新しいTODOを作成する
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	//TODOServiceを使用して新しいTODOを作成する
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		//作成中にエラーが発生した場合、そのエラー呼び出し元に返す
		return nil, err
	}
	//作成されたTODOをレスポンスとして返す
	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
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
