package handler

import (
	"net/http"
)

type PanicHandler struct{}

// ServeHTTP implements http.Handler.
func (p *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/do-panic" {
		panic("Intentional panic!")
	}
	// 他のパスへのリクエストを処理する
	http.NotFound(w, r)
}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}



