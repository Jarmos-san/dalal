// Package application tests the behavior of the application container.
//
// These tests validate correct wiring of dependencies, server lifecycle management, and
// basic HTTP request handling.
package application

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/config"
)

// newTestLogger returns a logger that discards all output. This ensures tests remain
// silent and deterministic.
func newTestLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(io.Discard, nil),
	)
}

// `TestNew_InitializesApplication` verifies that New correctly initializes the
// application struct and wires the HTTP server with the provided configuration and
// handler.
func TestNew_InitializesApplication(t *testing.T) {
	cfg := config.Config{
		Addr:         ":0", // use ephemeral port
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}

	handler := http.NewServeMux()
	logger := newTestLogger()

	app := New(cfg, handler, logger)

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

// `TestRunAndShutdown` verifies that the server can start and shut down gracefully
// without returning unexpected errors.
func TestRunAndShutdown(t *testing.T) {
	// Create a listener on an ephemeral port to avoid conflicts
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer ln.Close()

	cfg := config.Config{
		Addr:         ln.Addr().String(),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	logger := newTestLogger()
	app := New(cfg, mux, logger)

	// Replace server listener manually to control lifecycle
	app.Server.Addr = ""
	app.Server.Handler = mux

	go func() {
		_ = app.Server.Serve(ln)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatalf("expected graceful shutdown, got error: %v", err)
	}
}

// `TestServer_HandlesRequest` verifies that the application server correctly routes
// HTTP requests using the configured handler.
func TestServer_HandlesRequest(t *testing.T) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to create listener: %v", err)
	}
	defer ln.Close()

	cfg := config.Config{
		Addr:         ln.Addr().String(),
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		IdleTimeout:  2 * time.Second,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	logger := newTestLogger()
	app := New(cfg, mux, logger)

	go func() {
		_ = app.Server.Serve(ln)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + ln.Addr().String() + "/health")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = app.Shutdown(ctx)
}
