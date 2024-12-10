package handler

import "net/http"

type PanicHandler struct {
}

func (ph *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("panic!!")
}
