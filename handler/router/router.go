package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	mux.Handle("/healthz", handler.NewHealthzHandler())
	mux.Handle("/do-panic", middleware.Recovery(handler.NewPanicHandler()))

	todoHandler := handler.NewTODOHandler(service.NewTODOService(todoDB))
	mux.Handle("/todos", todoHandler)
	return mux
}
