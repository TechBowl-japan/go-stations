package middleware

import (
	"fmt"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Error:", err)
			}

			fmt.Println("Recoverd!")
			h.ServeHTTP(w, r)
		}()
	}
	return http.HandlerFunc(fn)
}
