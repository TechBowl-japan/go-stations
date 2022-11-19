package middleware

import (
	"fmt"
	"net/http"
	"os"
)

func BasicAuth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start Basic Authentication")
		if checkAuth(r) == false {
			http.Error(w, "Unauthorized", 401)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func checkAuth(r *http.Request) bool {
	BASIC_AUTH_USER_ID, BASIC_AUTH_PASSWORD, ok := r.BasicAuth()
	// fmt.Println(BASIC_AUTH_USER_ID)
	// fmt.Println(BASIC_AUTH_PASSWORD)
	if ok == false {
		return false
	}
	return BASIC_AUTH_PASSWORD == os.Getenv("BASIC_AUTH_PASSWORD") && BASIC_AUTH_USER_ID == os.Getenv("BASIC_AUTH_USER_ID")
}
