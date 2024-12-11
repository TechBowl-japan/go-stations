package middleware

import (
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		expexted_user := os.Getenv("BASIC_AUTH_USER_ID")
		expexted_pass := os.Getenv("BASIC_AUTH_PASSWORD")

		user, pass, ok := r.BasicAuth()

		if !ok || user != expexted_user || pass != expexted_pass {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)

	}
	return http.HandlerFunc(fn)
}
