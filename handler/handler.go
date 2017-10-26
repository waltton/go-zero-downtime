package handler

import (
	"net/http"
)

type Handler struct{}

func New() http.Handler {
	h := Handler{}
	mux := http.NewServeMux()

	mux.HandleFunc("/hello", h.hello)
	mux.HandleFunc("/delay", h.delay)

	return mux
}
