package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
	"golang.org/x/sync/errgroup"
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

	httpLn, err := net.Listen("tcp", defaultPort)
	if err != nil {
		return err
	}
	mux := router.NewRouter(todoDB)
	srv := &http.Server{
		Addr:    defaultPort,
		Handler: mux,
	}

	log.Println("server is running on", port)
	eg, ectx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		return srv.Serve(httpLn)
	})

	// Handle signals from the OS
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, os.Interrupt)
	select {
	case <-sigCh:
		fmt.Println("received signal, start exiting server gracefully")
	case <-ectx.Done():
	}

	tctx, cancel := context.WithTimeout(ectx, 1*time.Second)
	defer cancel()

	if err := srv.Shutdown(tctx); err != nil {
		fmt.Fprintf(os.Stderr, "Server shutdown: %s\n", err)
	}

	if err := eg.Wait(); err != nil {
		if http.ErrServerClosed == err {
			fmt.Fprint(os.Stderr, "server closed gracefully\n")
		} else {
			fmt.Fprintf(os.Stderr, "unhandled error: %s\n", err)
		}
	}

	return nil
}
