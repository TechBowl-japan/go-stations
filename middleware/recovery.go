package middleware

import "net/http"

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
	}
	// TODO: hを呼び出す必要がある
	// http.HandlerFuncは、func(w http.ResponseWriter, r *http.Request)のような関数をhttp.Handlerに変換する
	return http.HandlerFunc(fn)
}
