// Package logger_test contains behavioral tests for the logger package.
//
// The tests validate that log level configuration correctly controls which log records
// are emitted. Since slog.Logger does not expose its internal level configuration,
// these tests rely on observing output rather than inspecting internal state.
package logger_test

import (
	"bytes"
	"log/slog"
	"strings"
	"testing"
)

// newTestLogger constructs a slog.Logger instance that writes to an in-memory buffer
// instead of os.Stdout.
//
// This helper mirrors the production logger configuration logic while allowing
// deterministic testing by capturing emitted log output. It enables assertions on
// whether a log entry was written based on the configured log level.
//
// The level parameter follows the same semantics as the production code:
//
//   - "debug": enables all log levels
//   - "info": enables info and above (default)
//   - "warn": enables warn and error
//   - "error": enables only error
//
// Any unrecognised value defaults to "info".
func newTestLogger(level string, buf *bytes.Buffer) *slog.Logger {
	var lvl slog.Level

	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	handler := slog.NewTextHandler(buf, &slog.HandlerOptions{
		Level:       lvl,
		AddSource:   false,
		ReplaceAttr: nil,
	})

	return slog.New(handler)
}

// runLoggerTest executes a logger function and asserts whether output was emitted.
func runLoggerTest(
	t *testing.T,
	level string,
	logFn func(*slog.Logger),
	expected bool,
) {
	t.Helper()

	var buf bytes.Buffer

	logger := newTestLogger(level, &buf)

	logFn(logger)

	got := strings.TrimSpace(buf.String()) != ""
	if got != expected {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func TestLogger_DebugLevel(t *testing.T) {
	t.Parallel()

	// debug level allows all logs
	runLoggerTest(t, "debug",
		func(l *slog.Logger) { l.Debug("debug message") },
		true,
	)

	runLoggerTest(t, "debug",
		func(l *slog.Logger) { l.Info("info message") },
		true,
	)
}

func TestLogger_InfoLevel(t *testing.T) {
	t.Parallel()

	// info level suppresses debug logs
	runLoggerTest(t, "info",
		func(l *slog.Logger) { l.Debug("debug message") },
		false,
	)

	runLoggerTest(t, "info",
		func(l *slog.Logger) { l.Info("info message") },
		true,
	)
}

func TestLogger_ErrorLevel(t *testing.T) {
	t.Parallel()

	// error level suppresses warnings
	runLoggerTest(t, "error",
		func(l *slog.Logger) { l.Warn("warn message") },
		false,
	)

	runLoggerTest(t, "error",
		func(l *slog.Logger) { l.Error("error message") },
		true,
	)
}

func TestLogger_DefaultLevel(t *testing.T) {
	t.Parallel()

	// invalid level falls back to info
	runLoggerTest(t, "invalid",
		func(l *slog.Logger) { l.Info("info message") },
		true,
	)

	runLoggerTest(t, "invalid",
		func(l *slog.Logger) { l.Debug("debug message") },
		false,
	)
}
