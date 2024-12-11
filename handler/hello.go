package handler

import (
	"fmt"
	"net/http"
	"time"
)

type HelloHandler struct {
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World")
}
