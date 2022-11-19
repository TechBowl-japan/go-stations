package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"github.com/joho/godotenv"
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

	errEnv := godotenv.Load("env/auth.env")
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}

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
	// TODO: サーバーをlistenする
	// サーバーをポート8080で起動
	//http.ListenAndServe(defaultPort, mux)
	srv := &http.Server{
		Addr:    defaultPort,
		Handler: mux,
	}

	//http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
	//     fmt.Fprintf(w, "Hello World")
	// })\
	//mux := http.NewServeMux()
	//mux.HandleFunc("/hello", hello)
	// err = srv.ListenAndServe()
	// if err != http.ErrServerClosed {
	// 	log.Fatalln("Server closed with error:", err)
	// }

	go func() {
		log.Println("Go Routine")
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Fatalln("Server closed with error:", err)
		}
	}()

	quit := make(chan os.Signal, 1)                    // Create Signal monitoring
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt) //Monitoring signals
	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		log.Println("Failed to gracefully shutdown:", err)
	}
	log.Println("Server shutdown")

	return nil
}
