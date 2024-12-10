package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKey string

func WithOS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uas := r.UserAgent()
		ua := useragent.Parse(uas)
		os := ua.OS

		k := ctxKey("os")
		ctx := context.WithValue(r.Context(), k, os)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
