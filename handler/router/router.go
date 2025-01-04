package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

// func NewRouter(todoDB *sql.DB) *http.ServeMux {
func NewRouter(todoDB *sql.DB) http.Handler {
	// HTTPリクエストのURLとそれに対応するハンドラを登録する。HTTPリクエストが来た時に、そのURLにマッチしたハンドラを呼び出す。
	mux := http.NewServeMux()
	//healthzエンドポイント追加
	mux.Handle("/healthz", handler.NewHealthzHandler())
	// TODOエンドポイント追加
	mux.Handle("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)))
	//panicエンドポイント追加
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//リクエストを受け取るたびにpanic!を発生
		panic("panic!")
	})
	//panicHandlerがパニックを起こした場合にリカバリー
	mux.Handle("/do-panic", middleware.Recovery(panicHandler))

	//OSエンドポイント追加
	//OS情報を取得してレスポンスとして返す
	osInfoHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//middleware.GetOSFromContext()でmiddlewareがContextに格納したOS情報を取り出す。
		osName := middleware.GetOSFromContext(r.Context())
		w.Header().Set("Content-Type", "application/json")
		//HTTPステータスコード200 OKを返します。
		w.WriteHeader(http.StatusOK)
		//実際のレスポンスボディ{"os":"Windows"}のようなJSONデータをクライアントに送る
		w.Write([]byte(`{"os": "` + osName + `"}`))
	})
	// mux.Handle("/os-info", middleware.OSContextMiddleware(osInfoHandler))
	mux.Handle("/os-info", middleware.BasicAuthMiddleware(osInfoHandler))

	// BasicedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("You have accessed a protected resource."))
	// })
	// mux.Handle("/Basiced", middleware.BasicAuthMiddleware(BasicedHandler))

	//ルーター(mux)に対してAccexxLoggerMiddlewareを適用。アクセスログ記録後、ルーター(mux)に処理を渡す
	accessLoggerMiddleware := middleware.AccessLoggerMiddleware(mux)
	//リクエストが来たらOSContextMiddlewareを実行、その後accessLoggerMiddlewareに処理を渡す
	return middleware.OSContextMiddleware(accessLoggerMiddleware)
	// OSloggerMux := middleware.OSContextMiddleware(accessLoggerMiddleware)
	// return OSloggerMux.(*http.ServeMux)
}
