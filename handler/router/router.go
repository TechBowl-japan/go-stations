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
	svc := service.NewTODOService(todoDB)
	mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	//mux.HandleFunc("/do-panic", handler.NewPanicHandler().ServeHTTP)
	// healthzHandler := handler.NewHealthzHandler()
	// mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	healthzHandler := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	panicHandler := handler.NewPanicHandler()
	mux.HandleFunc("/do-panic1", panicHandler.ServeHTTP)
	mux.Handle("/do-panic", middleware.Recovery(panicHandler))
	//mux.Handle("/admin", middleware.Recovery(http.HandlerFunc(handleAdmin)))
	//trial := http.HandlerFunc(panicHandler.handleDoPanic)
	//mux.Handle("/do-panic", middleware.Recovery(panicHandler))
	//mux.Handle("/do-panic", middleware.Recovery(trial))
	return mux
}
