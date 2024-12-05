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
	case http.MethodPut: //PUTメソッドの場合
		h.handleUpdate(w, r) //TODO編集の処理を呼び出す
	case http.MethodGet: //GETメソッドの場合
		h.handleRead(w, r) //TODO取得の処理を呼び出す
	case http.MethodDelete: //DELETEメソッドの場合
		h.handleDelete(w, r) //TODO削除の処理を呼び出す
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

func (h *TODOHandler) handleRead(w http.ResponseWriter, r *http.Request) {
	//ReadTODORequest構造体のインスタンスを作成
	req := &model.ReadTODORequest{}
	//クエリパラメータを取得
	query := r.URL.Query()

	//"prev_id"パラメータを解析
	if prevIDStr := query.Get("prev_id"); prevIDStr != "" {
		var err error
		//文字列をint64に変換
		req.PrevID, err = strconv.ParseInt(prevIDStr, 10, 64)
		if err != nil {
			//エラーが発生した場合、400BadRequestを返す
			log.Printf("Error parsing prev_id: %v", err)
			http.Error(w, "Invalid prev_id", http.StatusBadRequest)
			return
		}
	}

	//"size"パラメータを解析
	if sizeStr := query.Get("size"); sizeStr != "" {
		var err error
		//文字列をint64に変換
		req.Size, err = strconv.ParseInt(sizeStr, 10, 64)
		if err != nil {
			//エラーが発生した場合、400BadRequestを返す
			log.Printf("Error parsing size: %v", err)
			http.Error(w, "Invalid size", http.StatusBadRequest)
			return
		}
	} else {
		//"size"が指定されていない場合、デフォルト値を設定
		req.Size = 5
	}

	// TODOの取得処理を呼び出す
	ctx := r.Context()
	res, err := h.Read(ctx, req)
	if err != nil {
		//エラーが発生した場合、500Internal Server Errorを返す
		log.Printf("Error reading TODOs: %v", err)
		http.Error(w, "Failed to read TODOs", http.StatusInternalServerError)
		return
	}

	//レスポンスヘッダを設定して成功ステータス(200 OK)を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//レスポンスをJSONとしてエンコード
	if err := json.NewEncoder(w).Encode(res); err != nil {
		//エンコード中にエラーが発生した場合、500 Internal Server Errorを返す
		log.Printf("Error reading response: %v", err)
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
	}
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	// TODOを取得するためにサービス層を呼び出す。
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		//エラーが発生した場合は呼び出し元に返す
		return nil, err
	}

	//サービス層から取得したTODOを変換
	//[]*model.TODO型のスライスを[]model.TODO型のスライスに変換
	convertedTodos := make([]model.TODO, len(todos))
	for i, todo := range todos {
		if todo != nil { //nilチェック
			convertedTodos[i] = *todo
		}
	}

	//変換されたTODOを含むレスポンスを返す
	return &model.ReadTODOResponse{
		TODOs: convertedTodos,
	}, nil
}

// handleUpdate handles the PUT request to update an existing TODO.
// handleUpdateは、既存のTODOを変更するためのPUTリクエストを処理する。
func (h *TODOHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	//リクエストボディを解析し、UpdateTODORequest構造体にデコードする。
	var req model.UpdateTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding UpdateTODORequest: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close() //リクエストボディをクローズする

	//必須フィールドが正しいかをチェックをする。
	if req.ID == 0 || req.Subject == "" {
		//IDが0かSubjectが空の場合、400BadRequestを返す
		http.Error(w, "Invalid ID or Subject", http.StatusBadRequest)
		return
	}

	//Contextを取得し、Updateメソッドを呼び出してTODOを更新する。
	ctx := r.Context()
	res, err := h.Update(ctx, &req)
	if err != nil {
		//TODOが見つからなかった場合
		if _, ok := err.(*model.ErrNotFound); ok {
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}

		//その他のエラーが発生した場合、500Internal Server Errorを返す
		log.Printf("Error updating TODO: %v", err)
		http.Error(w, "Failed to update TODO", http.StatusInternalServerError)
		return
	}

	//レスポンスヘッダを設定し、成功ステータス(200 OK)を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		//レスポンスヘッダのエンコードに失敗した場合、500Internal Server Errorを返す
		http.Error(w, "Faild to encode JSON", http.StatusInternalServerError)
	}
}

// Update handles the endpoint that updates the TODO.
// UpdateはTODOの更新を行うエンドポイントを処理します。
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	//TODOServiceのUpdateTODOメソッドを呼び出してTODOを更新する
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		//更新中にエラーが発生した場合、そのエラーを呼び出し元に返す。
		return nil, err
	}

	//更新されたTODOをレスポンスとして返す
	return &model.UpdateTODOResponse{
		TODO: *todo,
	}, nil
}

func (h *TODOHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	//リクエストボディを解析し、DeleteTODORequest構造体にデコードする。
	var req model.DeleteTODORequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding DeleteTODORequest : %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	//IDsが空かどうかを確認
	if len(req.IDs) == 0 {
		http.Error(w, "IDs are required", http.StatusBadRequest)
		return
	}

	//コンテキストを取得し、削除処理を呼び出す
	ctx := r.Context()
	res, err := h.Delete(ctx, &req) //正しく2つの戻り値を処理
	if err != nil {
		if _, ok := err.(*model.ErrNotFound); ok {
			http.Error(w, "TODO not found", http.StatusNotFound)
			return
		}

		log.Printf("Error deleting TODORequest: %v", err)
		http.Error(w, "Failed to delete TODO", http.StatusInternalServerError)
		return
	}

	//レスポンスヘッダを設定し、成功ステータス（200 OK）を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil { //resをエンコード
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

// Delete handles the endpoint that deletes the TODOs.
// TODOServiceのDeleteTODOメソッドを呼び出し、TODOを削除
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	if err := h.svc.DeleteTODO(ctx, req.IDs); err != nil {
		return nil, err
	}
	return &model.DeleteTODOResponse{}, nil
}
