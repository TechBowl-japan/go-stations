package router

import (
	"database/sql"
	"net/http"
   
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/service"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	
	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)
	panicHandler := handler.NewPanicHandler()

	mux := http.NewServeMux()
	mux.Handle("/healthz", &handler.HealthzHandler{})
	mux.Handle("/todos", todoHandler)
	mux.Handle("/do-panic", middleware.Recovery(panicHandler))
	return mux
}
