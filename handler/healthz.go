package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
)

// A HealthzHandler implements health check endpoint.
type HealthzHandler struct{}

// ServeHTTP handles the HTTP request and writes the health check response.
func (h *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// HealthzResponse 構造体を作成し、Message に "OK" をセット
	response := model.HealthzResponse{
		Message: "OK",
	}

// JSONシリアライズを実行
w.Header().Set("Content-Type", "application/json")
err := json.NewEncoder(w).Encode(response)
if err != nil {
	// シリアライズエラーが発生した場合はエラーログを出力
	log.Println("Failed to encode response:", err)
	w.WriteHeader(http.StatusInternalServerError)
	return
}
}
