package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/waltton/go-zero-downtime/config"
	"github.com/waltton/go-zero-downtime/handler"
)

type Server struct {
	ctx context.Context
	cfg *config.Server
	srv http.Server
}

func New(ctx context.Context, cfg *config.Server) *Server {
	handler := handler.New()

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: handler,
	}

	return &Server{ctx, cfg, srv}
}

func (s Server) Shutdown() error {
	<-s.ctx.Done()

	shutdownInterval := time.Second * time.Duration(s.cfg.IntervalToShutdown)
	log.Printf("Shutting down the server with timeout of %s...\n", shutdownInterval)

	ctx, cancel := context.WithTimeout(s.ctx, shutdownInterval)
	defer cancel()

	err := s.srv.Shutdown(ctx)
	log.Println("Shutdown err", err)
	return err
}

func (s Server) Run() error {
	var wg sync.WaitGroup
	var err error
	wg.Add(2)

	exit := make(chan bool)

	go func() {
		defer wg.Done()

		select {
		case <-exit:
			return
		case <-s.ctx.Done():
		}

		shutdownInterval := time.Second * time.Duration(s.cfg.IntervalToShutdown)
		log.Printf("Shutting down the server with timeout of %s...\n", shutdownInterval)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownInterval)
		defer cancel()

		if err := s.srv.Shutdown(ctx); err != nil {
			log.Printf("Error while shuting down the server, error %s\n", err)
		}
	}()

	go func() {
		err = s.srv.ListenAndServe()
		if err != http.ErrServerClosed {
			exit <- true
		}
		wg.Done()
	}()

	wg.Wait()
	return err
}
