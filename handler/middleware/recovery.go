package middleware

import (
	"log"
	"net/http"
)

func Recovery(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// TODO: ここに実装をする
		defer func() {
			// panicが発生した時にrecoverする
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}
				w.WriteHeader(http.StatusInternalServerError)

				// エラーメッセージをレスポンスに書き込む
				_, writeErr := w.Write([]byte("Internal Server Error"))
				if writeErr != nil {
					log.Printf("write error: %v", writeErr)
				}
			}
		}()
		// hからServeHTTPを呼び出してhttpリクエストをchainさせる
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func PanicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Panic!")
}
