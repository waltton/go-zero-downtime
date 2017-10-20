package handler

import "net/http"
import "io"

func (h Handler) hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello")
}
