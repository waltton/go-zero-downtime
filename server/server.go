package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/waltton/go-zero-downtime/config"
	"github.com/waltton/go-zero-downtime/handler"
)

type Server struct {
	ctx     context.Context
	cfg     *config.Server
	handler http.Handler
}

func New(ctx context.Context, cfg *config.Server) *Server {
	handler := handler.New()
	return &Server{ctx, cfg, handler}
}

func (s Server) Run() error {
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Host, s.cfg.Port),
		Handler: s.handler,
	}
	return srv.ListenAndServe()
}
