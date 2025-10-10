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

	// Register Health Check Endpoint
	mux.Handle("/healthz", &handler.HealthzHandler{})

	// Initialize TODOService and TODOHandler
	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)

	// Register TODO Endpoint
	mux.Handle("/todos", todoHandler)
	mux.Handle("/todos/", todoHandler)

	return mux
}
