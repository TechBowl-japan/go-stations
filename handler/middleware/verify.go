package middleware

import (
	"log"
	"net/http"
	"os"
)

// BasicAuth は Basic 認証を実装するミドルウェア
func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// 環境変数から User ID と Password を取得
		validUser := os.Getenv("BASIC_AUTH_USER_ID")
		validPass := os.Getenv("BASIC_AUTH_PASSWORD")

		// リクエストの認証情報を取得
		user, pass, ok := r.BasicAuth()

		// 認証情報が不正、または未入力の場合
		if !ok || user != validUser || pass != validPass {
			// ログに認証エラーの詳細を出力（コンソールに表示）
			log.Printf("HTTP Status Code:401 - User: %s, Password: %s",
				user, pass)

			// `WWW-Authenticate` ヘッダーを設定して再認証を促す
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized: Invalid credentials", http.StatusUnauthorized)
			return
		}

		// 認証成功時は次の Handler を実行
		log.Printf("Successful authentication - User: %s", user)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
