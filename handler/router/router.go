package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	
	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)

	mux := http.NewServeMux()
	mux.Handle("/healthz", &handler.HealthzHandler{})
	mux.Handle("/todos", todoHandler)
	return mux
}
