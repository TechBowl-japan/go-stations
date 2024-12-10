package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type ctxKey string

const osKey = ctxKey("os")

func WithOS(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		uas := r.UserAgent()
		ua := useragent.Parse(uas)
		os := ua.OS

		ctx := context.WithValue(r.Context(), osKey, os)
		r = r.WithContext(ctx)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func GetOSFromContext(ctx context.Context) string {
	if os, ok := ctx.Value(osKey).(string); ok {
		return os
	}

	return ""
}
