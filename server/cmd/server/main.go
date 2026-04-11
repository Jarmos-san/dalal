// Package `main` is the entry point of the application.
//
// It is responsible for initializing configuration, setting up HTTP routing,
// constructing the application container and managing the server lifecycle including
// graceful shutdown on system signals.
package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/application"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// `main` initialises and runs the HTTP server.
//
// It performs the following steps:
//   - Loads configurations from environment variables.
//   - Sets up HTTP routes.
//   - Constructs the application container.
//   - Starts the server in a seperate goroutine.
//   - Listens for OS signals to trigger graceful shutdown.
//
// The server is gracefully shutdown when an interrupt or termination signal is
// received, allowing in-flight requests to complete wihin a timeout period.
func main() {
	// Load application configuration from environment variables.
	cfg := config.LoadConfig()

	// Initialise the HTTP request multiplexer and register routes.
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Construct the app container with configurations and the handler.
	app := application.New(cfg, mux)

	// Create a context that is cancelled on interrupt or kill signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	defer stop()

	// Start the HTTP server in a seperate goroutine. This allows the `main` goroutine
	// to listen for shutdown signals
	go func() {
		log.Printf("starting server...")
		if err := app.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Block until a shutdown signal is received.
	<-ctx.Done()

	// Create a timeout-bound context for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt a graceful shutdown of the server.
	if err := app.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Printf("server shutdown completed gracefully...")
}
