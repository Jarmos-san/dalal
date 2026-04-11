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
		Level: lvl,
	})

	return slog.New(handler)
}

// TestLogger_LevelFiltering verifies that the logger correctly filters log messages
// according to the configured severity level.
//
// The test suite is table-driven and focuses on behavioral validation:
//
//   - Ensures that log messages at or above the configured level are emitted.
//   - Ensures that log messages below the configured level are suppressed.
//   - Confirms that invalid level inputs fall back to the default ("info").
//
// Instead of inspecting internal logger state, this test captures the logger output
// and checks whether any log entry was written. This approach reflects real-world
// usage and aligns with slog's design, where handlers determine emission behavior.
//
// Each test case defines:
//
//   - level: the configured log level
//   - logFn: the logging operation to execute
//   - expectedOutput: whether output is expected to be emitted
func TestLogger_LevelFiltering(t *testing.T) {
	tests := []struct {
		name           string
		level          string
		logFn          func(*slog.Logger)
		expectedOutput bool
	}{
		{
			name:           "debug level allows debug",
			level:          "debug",
			logFn:          func(l *slog.Logger) { l.Debug("debug message") },
			expectedOutput: true,
		},
		{
			name:           "info level supporesses debug",
			level:          "info",
			logFn:          func(l *slog.Logger) { l.Debug("debug message") },
			expectedOutput: false,
		},
		{
			name:  "error level suppresses warn",
			level: "error",
			logFn: func(l *slog.Logger) {
				l.Warn("warn message")
			},
			expectedOutput: false,
		},
		{
			name:  "error level allows error",
			level: "error",
			logFn: func(l *slog.Logger) {
				l.Error("error message")
			},
			expectedOutput: true,
		},
		{
			name:  "default level is info",
			level: "invalid",
			logFn: func(l *slog.Logger) {
				l.Info("info message")
			},
			expectedOutput: true,
		},
		{
			name:  "default level suppresses debug",
			level: "invalid",
			logFn: func(l *slog.Logger) {
				l.Debug("debug message")
			},
			expectedOutput: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newTestLogger(tt.level, &buf)

			tt.logFn(logger)

			output := buf.String()
			gotOutput := strings.TrimSpace(output) != ""

			if gotOutput != tt.expectedOutput {
				t.Fatalf(
					"expected output: %v, got %v, output: %q",
					tt.expectedOutput,
					gotOutput,
					output,
				)
			}
		})
	}
}
