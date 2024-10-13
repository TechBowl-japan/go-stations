package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/mileusna/useragent"
)

func Device(h http.Handler) http.Handler {
	type osContextKey string

	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), osContextKey("os"), ua.OS)
		log.Println(ctx.Value(osContextKey("os")))
		h.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
