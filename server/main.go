package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Config struct {
	Addr         string // e.g., ":8000"
	ReadTimeout  int    // seconds
	WriteTimeout int    // seconds
	IdleTimeout  int    // seconds
}

type Application struct {
	Config  Config
	Server  *http.Server
	Handler http.Handler
}

func NewApplication(cfg Config, handler http.Handler) *Application {
	srv := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	return &Application{
		Config:  cfg,
		Server:  srv,
		Handler: handler,
	}
}

func (a Application) Run() error {
	log.Printf("starting server on %s", a.Config.Addr)
	return a.Server.ListenAndServe()
}

func (a Application) Shutdown(ctx context.Context) error {
	return a.Server.Shutdown(ctx)
}

func main() {
	cfg := Config{
		Addr:         ":8000",
		ReadTimeout:  10,
		WriteTimeout: 10,
		IdleTimeout:  10,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	application := NewApplication(cfg, mux)

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
