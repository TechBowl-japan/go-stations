package middleware

import (
	"log"
	"net/http"
	"os"
)

var (
	requiredID       = os.Getenv("BASIC_AUTH_USER_ID")
	requiredPassword = os.Getenv("BASIC_AUTH_PASSWORD")
)

func checkAuth(r *http.Request) bool {
	// requiredID := os.Getenv("BASIC_AUTH_USER_ID")
	// requiredPassword := os.Getenv("BASIC_AUTH_PASSWORD")

	log.Printf("Environment BASIC_AUTH_USER_ID: %s", requiredID)
	log.Printf("Environment BASIC_AUTH_PASSWORD: %s", requiredPassword)
	userID, password, ok := r.BasicAuth()
	if !ok {
		return false
	}
	log.Printf("Required ID: %s, Provided ID: %s", requiredID, userID)
	log.Printf("Required Password: %s, Provided Password: %s", requiredPassword, password)
	return userID == requiredID && password == requiredPassword
}

// MiddlewareでHTTPリクエストがサーバーに届く前に共通処理開始
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//認証失敗
		if !checkAuth(r) {
			w.Header().Set("WWW-Authenticate", `Basic realm= "restricted"`)
			// w.WriteHeader(http.StatusUnauthorized) //401
			// http.Error(w, "Unauthorized", http.StatusUnauthorized)
			http.Error(w, "ユーザー名とパスワードを入力してください。", http.StatusUnauthorized)
			return
		}

		//認証成功時は次の処理へ
		w.Write([]byte("認証成功！\n"))
		next.ServeHTTP(w, r)
	})
}
