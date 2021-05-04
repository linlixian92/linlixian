package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/sync/errgroup"
)

func main() {
	mux := http.NewServeMux()
	serverOut := make(chan struct{})
	mux.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		serverOut <- struct{}{}
	})
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return server.ListenAndServe() })

	g.Go(func() error {
		select {
		case <-ctx.Done():
			log.Println("errgroup exit...")
		case <-serverOut:
			log.Println("server will out...")
		}
		err := server.Shutdown(ctx)
		return err
	})

	g.Go(func() error {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-quit:
			return errors.New("get os signal:" + sig.String())
		}
	})
	log.Printf("errgroup exiting: %+v\n", g.Wait())

}
