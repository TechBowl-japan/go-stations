package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
// TODO: ここに実装をする

	// 前処理をここに実装する
    defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()
	h.ServeHTTP(w, r)


}
return http.HandlerFunc(fn)
}