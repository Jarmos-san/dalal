// Package logger provides utilities for constructing structured loggers used across
// the application.
//
// It exposes a factory function for creating a configured `slog.Logger` instance with a
// specified log level and output format.
package logger

import (
	"log/slog"
	"os"
)

// New constructs and returns a new `slog.Logger` configured with the provided log
// level.
//
// The level parameter controls the minimum severity of log records that will be
// emitted. Supported values are:
//
//   - "debug": enable debug, info, warn and error logs.
//   - "info": enable info, warn and error logs (default).
//   - "warn": enable warn and error logs.
//   - "error": enables only error logs.
//
// Any unrecognised value defaults to "info".
//
// The logger writes output to standard output (os.Stdout) using a human-readable text
// format via slog.NewTextHandler. This is suitable for development and basic production
// use. For structured logging in production systems, slog.NewJSONHandler will be used
// in the future instead.
//
// The returned logger is safe for concurrent usage.
func New(level string) *slog.Logger {
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

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:       lvl,
			AddSource:   false,
			ReplaceAttr: nil,
		},
	)

	return slog.New(handler)
}
