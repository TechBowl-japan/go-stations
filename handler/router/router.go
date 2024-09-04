package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	// /healthzエンドポイントを登録
	healthzHnadler := &handler.HealthzHandler{}
	mux.Handle("/healthz", healthzHnadler)

	// TODOService インスタンスを作成
	todoService := &service.TODOService{DB: todoDB}

	// /todos エンドポイントを登録
	todoHandler := &handler.TODOHandler{SVC: todoService}
	mux.Handle("/todos", todoHandler)

	return mux
}
