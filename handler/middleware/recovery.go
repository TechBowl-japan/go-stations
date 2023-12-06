package middleware

import "net/http"

func Recovery(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
// TODO: ここに実装をする
}
return http.HandlerFunc(fn)
}