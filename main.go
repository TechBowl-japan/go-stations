package main

import (
	"log"
	"net/http"
	"os"
	"time"

	// errors パッケージをインポート
	"github.com/TechBowl-japan/go-stations/db"
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
		log.Println("Failed to load time zone:", err)
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		log.Println("Failed to initialize database:", dbPath, err)
		return err
	}
	defer todoDB.Close()

	// NOTE: 新しいエンドポイントの登録はrouter.NewRouterの内部で行うようにする
	mux := router.NewRouter(todoDB)

	// TODO: サーバーをlistenする
	//サーバーを非同期で起動
	log.Printf("Starting server on port %s\n", port)
	//ListenAndServe(ネットワークアドレス,ハンドラ)はHTTPサーバーを起動
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Printf("Server stopped unexpectedly: %v\n", err)
	}

	// // サーバー完全起動後に /do-panic に自動リクエストを送信
	// time.Sleep(1 * time.Second) // サーバー起動待機
	// _, err = http.Get("http://localhost" + port + "/do-panic")
	// if err != nil {
	// 	log.Printf("Failed to send request to /do-panic: %v\n", err)
	// }

	// // return nil

	// time.Sleep(1 * time.Second) // サーバーが起動するのを待つ
	// resp, err := http.Get("http://localhost" + port + "/os-info")
	// if err != nil {
	// 	log.Printf("Failed to send request to /os-info: %v\n", err)
	// 	return err
	// }
	// //接続が開いたままだとリソースリーク（メモリやネットワーク接続の浪費）を引き起こすため、defer()で確実に解放する
	// defer resp.Body.Close()

	// log.Printf("Response from /os-info: %v\n", resp.Status)
	return nil

	// 	log.Printf("Starting server on port %s\n", port)
	// 	err = http.ListenAndServe(port, mux)
	// 	if err != nil {
	// 		log.Printf("Failed to start server on port %s: %v\n", port, err)
	// 		return fmt.Errorf("server failed to start on %s: %w", port, err)

	// 	}
	// 	return nil
}
