// Package `main` is the entry point of the application.
//
// It is responsible for initializing configuration, setting up HTTP routing,
// constructing the application container and managing the server lifecycle including
// graceful shutdown on system signals.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/Jarmos-san/arthika/server/internal/application"
	"github.com/Jarmos-san/arthika/server/internal/config"
	"github.com/Jarmos-san/arthika/server/internal/dto"
	"github.com/Jarmos-san/arthika/server/internal/handlers"
	"github.com/Jarmos-san/arthika/server/internal/logger"
	"github.com/Jarmos-san/arthika/server/internal/services"
)

// shutdownTimeout defines the maximum duration allowed for gracefully shutting
// down the HTTP server.
//
// During shutdown, the server stops accepting new connections and waits for in-flight
// requests to complete. If this timeout is exceeded, the shutdown process is aborted
// and any remaining connections may be terminated.
const shutdownTimeout = 5 * time.Second

// `main` initialises and runs the HTTP server.
//
// It performs the following steps:
//   - Loads configurations from environment variables.
//   - Sets up HTTP routes.
//   - Constructs the application container.
//   - Starts the server in a separate goroutine.
//   - Listens for OS signals to trigger graceful shutdown.
//
// The server is gracefully shutdown when an interrupt or termination signal is
// received, allowing in-flight requests to complete wihin a timeout period.
func main() { //nolint:funlen
	// Load application configuration from environment variables.
	cfg := config.LoadConfig()

	logger := logger.New(cfg.LogLevel)

	// Initialise the HTTP request multiplexer and register routes.
	mux := http.NewServeMux()

	// Register the routes and their handlers
	userService := services.NewUserService()
	userHandler := handlers.NewUserHandler(userService, logger)
	mux.HandleFunc("GET /users/", userHandler.GetUser)
	mux.HandleFunc("POST /users/register", userHandler.CreateUser)
	mux.HandleFunc("POST /login", userHandler.LoginUser)

	// Temporary scratch handler for experimental requirements only
	mux.HandleFunc("GET /scratch/", func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/vnd.api+json")
		writer.WriteHeader(http.StatusOK)

		data := dto.ResourceObject{
			Type: "user",
			ID:   uuid.NewString(),
			Attributes: map[string]any{
				"name": "John Doe",
			},
			Relationships: nil,
			Links:         nil,
		}

		newResp := dto.NewSingleDocument(data)

		err := json.NewEncoder(writer).Encode(newResp)
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)

			return
		}
	})

	// Construct the app container with configurations and the handler.
	app := application.New(cfg, mux, logger)

	// Create a context that is cancelled on interrupt or kill signals.
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
	)
	defer stop()

	// Start the HTTP server in a separate goroutine. This allows the `main` goroutine
	// to listen for shutdown signals
	go func() {
		logger.Info("starting server")

		err := app.Run()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server startup failed", "error", err.Error())
		}
	}()

	// Block until a shutdown signal is received.
	<-ctx.Done()

	// Create a timeout-bound context for graceful shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt a graceful shutdown of the server.
	err := app.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("server shutdown failed", "error", err.Error())
	}

	logger.Info("server shutdown completed gracefully")
}
