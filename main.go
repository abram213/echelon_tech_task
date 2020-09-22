package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"task_ws_et/app"
	"task_ws_et/config"
	"task_ws_et/storage"
	"time"
)

func main() {
	conf, err := config.New(".env")
	if err != nil {
		log.Fatalf("init config err: %v\n", err)
	}

	st := storage.New()

	a, err := app.New(st, *conf)
	if err != nil {
		log.Fatalf("new app err: %v\n", err)
	}

	s := &http.Server{
		Addr:        ":" + a.Port,
		Handler:     a.Router,
		ReadTimeout: 1 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)

		signal.Notify(sigCh, os.Interrupt)

		<-sigCh
		fmt.Println("signal caught. shutting down...")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Fatalf("server shutdown err: %v", err)
		}
		close(done)
	}()

	fmt.Printf("serving at http://localhost%s\n", s.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("serving err: %v", err)
		close(done)
	}
}