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
	healthzHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", healthzHandler)

	todoHandler := handler.NewTODOHandler(service.NewTODOService(todoDB))
	mux.Handle("/todos", todoHandler)
	return mux
}
