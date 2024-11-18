package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
// HealthzHandlerはヘルスチェックエンドポイントを実装します。
type HealthzHandler struct{}

// NewHealthzHandler returns HealthzHandler based http.Handler.
// NewHealthzHandlerは、httpベースのHealthzHandlerを返します。Handler。
func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

// ServeHTTP implements http.Handler interface.
// ServeHTTPはhttp.Handlerインターフェースを実装します。
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//ヘッダーを設定して、レスポンス形式をJSONと通知する
	w.Header().Set("Content Type", "application/json")
	//JSONで返す内容を準備する
	response := &model.HealthzResponse{
		Message: "OK", //カンマあり
	}
	err := json.NewEncoder(w).Encode(response)
	//JsonエンコーダーでresponseをJsonにして、直接wに書き込む
	if err != nil {
		//エラーがあれば、ログに記録して処理を中断
		log.Println("Error encoding response:", err) //エラー内容をログに出力
		return                                       //処理をここで終了
	}
}
