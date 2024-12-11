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

	helloHandler := &handler.HelloHandler{}
	mux.Handle("/", helloHandler)

	healthHandler := handler.NewHealthzHandler()
	mux.Handle("/healthz", middleware.BasicAuth(middleware.WithOS(middleware.DisplayAccessLog(healthHandler))))

	todoService := service.NewTODOService(todoDB)
	todoHandler := handler.NewTODOHandler(todoService)
	mux.Handle("/todos", todoHandler)

	panicHandler := &handler.PanicHandler{}
	mux.Handle("/do-panic", middleware.Recovery(panicHandler))

	return mux
}
