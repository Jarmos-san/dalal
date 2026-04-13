// Package application tests the behavior of the application container.
//
// These tests validate correct wiring of dependencies, server lifecycle management, and
// basic HTTP request handling.
package application_test

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/application"
	"github.com/Jarmos-san/arthika/server/internal/config"
)

// newTestLogger returns a logger that discards all output. This ensures tests remain
// silent and deterministic.
func newTestLogger() *slog.Logger {
	return slog.New(
		slog.DiscardHandler,
	)
}

// `TestNew_InitializesApplication` verifies that New correctly initializes the
// application struct and wires the HTTP server with the provided configuration and
// handler.
func TestNew_InitializesApplication(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		Addr:         ":0", // use ephemeral port
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		LogLevel:     "info",
	}

	handler := http.NewServeMux()
	logger := newTestLogger()

	app := application.New(cfg, handler, logger)

	if app.Config != cfg {
		t.Errorf("expected config %+v, got %+v", cfg, app.Config)
	}

	if app.Server == nil {
		t.Fatal("expected server to be initialized, got nil")
	}

	if app.Handler != handler {
		t.Errorf("expected handler to be set")
	}

	if app.Logger != logger {
		t.Errorf("expected logger to be set")
	}

	if app.Server.Addr != cfg.Addr {
		t.Errorf("expected server Addr %s, got %s", cfg.Addr, app.Server.Addr)
	}
}

// mustListen creates a TCP listener on an available port and fails the test if it
// cannot.
func mustListen(t *testing.T) net.Listener {
	t.Helper()

	listenConfig := net.ListenConfig{
		Control:         nil,
		KeepAlive:       0,
		KeepAliveConfig: net.KeepAliveConfig{}, //nolint:exhaustruct
	}

	ln, err := listenConfig.Listen(context.Background(), "tcp", ":0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}

	return ln
}

// `TestRunAndShutdown` verifies that the server can start and shut down gracefully
// without returning unexpected errors.
func TestRunAndShutdown(t *testing.T) {
	t.Parallel()

	// Create a listener on an ephemeral port to avoid conflicts
	listener := mustListen(t)

	defer func() { _ = listener.Close() }()

	cfg := config.Config{
		Addr:         listener.Addr().String(),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
		LogLevel:     "info",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := newTestLogger()
	app := application.New(cfg, mux, logger)

	// Replace server listener manually to control lifecycle
	app.Server.Addr = ""
	app.Server.Handler = mux

	go func() {
		_ = app.Server.Serve(listener)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	shutdownErr := app.Shutdown(ctx)
	if shutdownErr != nil {
		t.Fatalf("expected graceful shutdown, got error: %v", shutdownErr)
	}
}

// doHealthCheck sends a GET request to the given URL and fails the test on error.
func doHealthCheck(t *testing.T, url string) *http.Response {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	client := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}

	return resp
}

// `TestServer_HandlesRequest` verifies that the application server correctly routes
// HTTP requests using the configured handler.
func TestServer_HandlesRequest(t *testing.T) {
	t.Parallel()

	listener := mustListen(t)

	defer func() { _ = listener.Close() }()

	cfg := config.Config{
		Addr:         listener.Addr().String(),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
		LogLevel:     "info",
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	logger := newTestLogger()
	app := application.New(cfg, mux, logger)

	go func() {
		_ = app.Server.Serve(listener)
	}()

	time.Sleep(100 * time.Millisecond)

	resp := doHealthCheck(t, "http://"+listener.Addr().String()+"/health")

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)
	defer shutdownCancel()

	_ = app.Shutdown(shutdownCtx)
}
