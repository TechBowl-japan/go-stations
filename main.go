package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	// 環境変数が設定されているか確認
	if os.Getenv("BASIC_AUTH_USER_ID") == "" || os.Getenv("BASIC_AUTH_PASSWORD") == "" {
		log.Fatal("IDとPASSWORDを設定してください")
	}

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
	securedMux := middleware.BasicAuth(mux)

	// Graceful Shutdown のための `http.Server`
	server := &http.Server{
		Addr:    port,
		Handler: securedMux,
	}

	// シグナルを受け取るためのコンテキスト
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// リクエストの処理が完了するまで待つための WaitGroup
	var wg sync.WaitGroup

	// サーバーを非同期で起動
	go func() {
		log.Println("Server running at http://localhost" + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("サーバー起動エラー: %v", err)
		}
	}()

	// シグナルを受け取るまで待機
	<-ctx.Done()
	log.Println("シグナル受信! サーバーのシャットダウンを開始...")

	// Graceful Shutdown のためのタイムアウト付きコンテキスト
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// HTTP サーバーのシャットダウン
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("サーバーシャットダウン失敗: %v", err)
	}

	// すべてのリクエストが完了するのを待つ
	wg.Wait()
	log.Println("すべてのリクエストが完了しました。サーバーを安全に停止します。")

	return nil
}
