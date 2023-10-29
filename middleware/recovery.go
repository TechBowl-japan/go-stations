package middleware

import "net/http"

// h.ServerHTTPの実行中にパニックが発生した場合に、500 Internal Server Errorを返す
func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			// recover(): panic()を捕捉する
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
	// http.HandlerFuncは、func(w http.ResponseWriter, r *http.Request)のような関数をhttp.Handlerに変換する
	return http.HandlerFunc(fn)
}
