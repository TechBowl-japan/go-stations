package middleware

import (
	"context"
	"net/http"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/mileusna/useragent"
)

// request元のOSを判定し、contextに保存する
func IdentifyOS(h http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := r.Context()
		ctx = context.WithValue(ctx, model.OSContextKey("os"), ua.OS)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
