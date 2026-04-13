// Package application provides the core application container responsible for wiring
// configuration, HTTP handlers and server lifecycle management.
//
// It acts as the composition root of the service, coordinating dependencies and
// exposing methods to start and gracefully shutdown the HTTP server.
package application

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/config"
)

// Application represents the main Application container.
//
// It holds the runtime configuration, the HTTP server instance and the root HTTP
// handler, and the structured logger. This type is responsible for managing the
// lifecycle of the HTTP server and coordinating cross-cutting concerns such as logging.
type Application struct {
	// Config contains the application configuration values.
	Config config.Config

	// Server is the HTTP server responsible for handling incoming requests.
	Server *http.Server

	// Handler is the root HTTP handler used by the server to route requests.
	Handler http.Handler

	// Logger is the structured logger used for emitting application-level logs.
	//
	// It is expected to be initialised by the caller and injected into the application.
	// The logger is used for lifecycle events such as startup and shutdown, as well as
	// other operational logging.
	Logger *slog.Logger
}

// New constructs and returns a new application instance.
//
// It initialises an `http.Server` using the provided configuration and handler, and
// associates a structured logger for application-level logging.
//
// The returned application is ready to be started via the `Run` method.
func New(cfg config.Config, handler http.Handler, logger *slog.Logger) *Application {
	server := &http.Server{ //nolint:exhaustruct
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Application{
		Config:  cfg,
		Server:  server,
		Handler: handler,
		Logger:  logger,
	}
}

// Run starts the HTTP server and begins listening for incoming requests.
//
// It logs the server start event and blocks until the server stops. Any error returned
// by `ListenAndServer` is propagated to the caller.
func (a *Application) Run() error {
	a.Logger.Info("server started", "addr", a.Config.Addr)

	err := a.Server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

// Shutdown gracefully stops the HTTP server.
//
// It attempts to shutdown the server using the provided context, allowing in-flight
// requests to complete before termination. If the context expires before shutdown
// completes, an error is returned.
func (a *Application) Shutdown(_ context.Context) error {
	a.Logger.Info("server shutdown")

	err := a.Server.ListenAndServe()
	if err != nil &&
		errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}
