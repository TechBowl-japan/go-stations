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
	mux.Handle("/healthz", handler.NewHealthzHandler())
	todoService := service.NewTODOService(todoDB)
	mux.Handle("/todos", handler.NewTODOHandler(todoService))
	return mux
}
