// Package config provides functionality for loading and managing application
// configuration values.
//
// Configuration is primarily sourced from environment variables with sensible defaults
// applied when variables are not set or is invalid.
package config

import (
	"fmt"
	"os"
	"time"
)

// defaultTimeout defines the default duration applied to HTTP server timeouts.
//
// It is used as the fallback value for ReadTimeout, WriteTimeout, and IdleTimeout when
// no explicit configuration is provided. The value is chosen to balance responsiveness
// and tolerance for slow clients.
const defaultTimeout = 10 * time.Second

// Config represents the runtime configuration for the application.
//
// All fields are populated via environment variables in the `LoadConfig()` function
// with fallback defaults when necessary.
type Config struct {
	// Addr is the network address the HTTP server listens on.
	// Example: ":8000"
	Addr string

	// ReadTimeout is the maximum duration for reading the entire request, including
	// the body.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before time out writes to the response.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the next request when
	// keep-alives are enabled.
	IdleTimeout time.Duration

	// LogLevel defines the minimum severity of logs emitted by the application logger.
	// Supported values include "debug", "info", "warn", and "error".
	LogLevel string
}

// getEnv() retrieves the value of the environment variable identified by key. If the
// variable is not set, it returns the provided fallback.
func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}

	return fallback
}

// getEnvDuration() retrieves an environment variable and attempts to parse it as a
// time.Duration.
//
// The expected format is a valid duration string such as "5s", "1m", "500ms". If the
// variable is not set or cannot be parsed, the fallback value is returned.
func getEnvDuration(key string, fallback time.Duration) (time.Duration, error) {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		parsedDuration, err := time.ParseDuration(val)
		if err != nil {
			return fallback, fmt.Errorf("invalid %s: %w", key, err)
		}

		return parsedDuration, nil
	}

	return fallback, nil
}

// LoadConfig constructs a Config by reading environment variables and applying
// default values where necessary.
//
// The following environment variables are supported:
//
//   - ADDR: HTTP server address (e.g., ":8000")
//   - READ_TIMEOUT: request read timeout (e.g., "10s")
//   - WRITE_TIMEOUT: response write timeout (e.g., "10s")
//   - IDLE_TIMEOUT: keep-alive idle timeout (e.g., "60s")
//   - LOG_LEVEL: the log level (e.g., "debug", "info", "warn", "error")
//
// If an environment variable is not set or contains an invalid value, the corresponding
// default value is used instead.
func LoadConfig() Config {
	defaultCfg := Config{
		Addr:         ":8000",
		ReadTimeout:  defaultTimeout,
		WriteTimeout: defaultTimeout,
		IdleTimeout:  defaultTimeout,
		LogLevel:     "info",
	}

	addr := getEnv("ADDR", defaultCfg.Addr)
	readTimeout, _ := getEnvDuration("READ_TIMEOUT", defaultCfg.ReadTimeout)
	writeTimeout, _ := getEnvDuration("WRITE_TIMEOUT", defaultCfg.WriteTimeout)
	idleTimeout, _ := getEnvDuration("IDLE_TIMEOUT", defaultCfg.IdleTimeout)
	logLevel := getEnv("LOG_LEVEL", defaultCfg.LogLevel)

	return Config{
		Addr:         addr,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
		LogLevel:     logLevel,
	}
}
