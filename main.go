package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/waltton/go-zero-downtime/config"
	"github.com/waltton/go-zero-downtime/server"
)

func waitSignal(ctx context.Context) func() error {
	return func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigs:
			return fmt.Errorf("signal '%s' received", sig)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func untilFail(wg *sync.WaitGroup, cancel context.CancelFunc, f func() error) {
	err := f()
	if err != nil && err != context.Canceled {
		log.Println(err)
		cancel()
	}
	wg.Done()
}

func main() {
	configFileName := flag.String("c", "./config.json", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadFromJSONFile(*configFileName)
	if err != nil {
		log.Panicf("Could not load configs, error: %v", err)
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go untilFail(&wg, cancel, waitSignal(ctx))

	wg.Add(1)
	srv := server.New(ctx, cfg.Server)
	go untilFail(&wg, cancel, srv.Run)

	wg.Wait()
}
