package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)
	healthzHandler := handler.NewHealthzHandler()
	mux.HandleFunc("/healthz", healthzHandler.ServeHTTP)
	//	mux.Handle("/healthz", HealthzHandler)
	//todoHandler := handler.NewTODOHandler(todoDB)
	// mux.HandleFunc("/todos", todoHandler.Create)
	//var svc *service.TODOService = service.NewTODOService(todoDB)
	//mux.HandleFunc("/todos", todoHandler.ServeHTTP)

	http.HandleFunc("/", hello)
	//http.ListenAndServe(":8000", nil)
	// TODO: サーバーをlistenする
	// サーバーをポート8080で起動
	http.ListenAndServe(defaultPort, mux)
	//http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "Hello World")
	// })

	return nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, Gopher!")
}
