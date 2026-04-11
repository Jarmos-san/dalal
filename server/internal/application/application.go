// Package `application` provides the core application container responsible for wiring
// configuration, HTTP handlers and server lifecycle management.
//
// It acts as the composition root of the service, coordinating dependencies and
// exposing methods to start and gracefully shutdown the HTTP server.
package application

import (
	"context"
	"log"
	"net/http"

	"github.com/Jarmos-san/arthika/server/internal/config"
)

// `application` represents the main application container.
//
// It holds the runtime configuration, the HTTP server instance and the root HTTP
// handler. This struct is responsible for managing the lifecycle of the HTTP server.
type application struct {
	// Config contains the application configuration values.
	Config config.Config

	// Server is the HTTP server responsible for handling incoming requests.
	Server *http.Server

	// Handler is hte root HTTP handler used by the server to route requests.
	Handler http.Handler
}

// `New` constructs and returns a new application instance.
//
// It initialises an `http.Server` using the provided configuration and handler. The
// returned application is ready to be started via the `application.Run` method.
func New(cfg config.Config, handler http.Handler) *application {
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      handler,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &application{
		Config:  cfg,
		Server:  server,
		Handler: handler,
	}
}

// `Run` starts the HTTP server and begins listening for incoming requests.
//
// It logs the server start event and blocks until the server stops. Any error returned
// by `ListenAndServer` is propagated to the caller.
func (a *application) Run() error {
	log.Printf("server is ready at: http://localhost%s", a.Config.Addr)
	return a.Server.ListenAndServe()
}

// `Shutdown` gracefully stops the HTTP server.
//
// It attempts to shutdown the server using the provided context, allowing in-flight
// requests to complete before termination. If the context expires before shutdown
// completes, an error is returned.
func (a *application) Shutdown(ctx context.Context) error {
	log.Print("shutting down...")
	return a.Server.Shutdown(ctx)
}
