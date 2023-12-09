package middleware

import "net/http"

func Recovery(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
// TODO: ここに実装をする

	// 前処理をここに実装する

	h.ServeHTTP(w, r)

	//後処理をここに実装する
	
}
return http.HandlerFunc(fn)
}