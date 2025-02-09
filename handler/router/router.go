package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	//HealthzHandlerエンドポイントを登録
	mux.Handle("/healthz", handler.NewHealthzHandler())

	//TODOHandlerエンドポイントを登録
	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)
	mux.Handle("/todos", todoHandler)

	//do-panicエンドポイントを登録
	mux.Handle("/do-panic", middleware.Recovery(func(w http.ResponseWriter, r *http.Request) {
		panic("パニックです")
	}))

	// /information エンドポイントを登録
	mux.Handle("/information", middleware.UserOs(func(w http.ResponseWriter, r *http.Request) {
	}))

	return mux
}
