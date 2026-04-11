// Package `config_test` contains black-box tests for the `config` package.
//
// These tests validate the public behaviour of configuration loading, ensuring correct
// handling of defaults, environment overrides and fallback behaviour for invalid
// inputs.
package config_test

import (
	"testing"
	"time"

	"github.com/Jarmos-san/arthika/server/internal/config"
)

// TestLoadConfig_Defaults() verifies that `config.LoadConfig()` returns the expected
// default configuration with no environment variables are set.
func TestLoadConfig_Defaults(t *testing.T) {
	cfg := config.LoadConfig()

	defaultCfg := config.Config{
		Addr:         ":8000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	if cfg != defaultCfg {
		t.Fatalf("expected default config %+v, got %+v", defaultCfg, cfg)
	}
}

// TestLoadConfig_FromEnv() verifies that environment variables correctly override
// default confguration values.
func TestLoadConfig_FromEnv(t *testing.T) {
	t.Setenv("ADDR", ":9000")
	t.Setenv("READ_TIMEOUT", "5s")
	t.Setenv("WRITE_TIMEOUT", "6s")
	t.Setenv("IDLE_TIMEOUT", "7s")

	cfg := config.LoadConfig()

	if cfg.Addr != ":9000" {
		t.Errorf("expected Addr :9000, got %s", cfg.Addr)
	}

	if cfg.ReadTimeout.String() != "5s" {
		t.Errorf("expected ReadTimeout 5s, got %s", cfg.ReadTimeout)
	}

	if cfg.WriteTimeout.String() != "6s" {
		t.Errorf("expected WriteTimeout 6s, got %s", cfg.WriteTimeout)
	}

	if cfg.IdleTimeout.String() != "7s" {
		t.Errorf("expected IdleTimeout 7s, got %s", cfg.IdleTimeout)
	}
}

// TestLoadConfig_InvalidDurationFallsBack verifies that invalid duration values in
// environment variables do not cause failure and instead fall back to default values.
func TestLoadConfig_InvalidDurationFallsBack(t *testing.T) {
	t.Setenv("READ_TIMEOUT", "invalid")

	cfg := config.LoadConfig()
	def := config.Config{
		ReadTimeout: 10 * time.Second,
	}

	if cfg.ReadTimeout != def.ReadTimeout {
		t.Errorf("expected fallback ReadTimeout %v, got %v",
			def.ReadTimeout, cfg.ReadTimeout)
	}
}

// TestLoadConfig_PartialOverride verifies that when only a subset of environment
// variables are provided, only those values are overridden and the remaining fields
// retain their default values.
func TestLoadConfig_PartialOverride(t *testing.T) {
	t.Setenv("ADDR", ":7000")

	cfg := config.LoadConfig()
	def := config.Config{
		ReadTimeout: 10 * time.Second,
	}

	if cfg.Addr != ":7000" {
		t.Errorf("expected Addr :7000, got %s", cfg.Addr)
	}

	if cfg.ReadTimeout != def.ReadTimeout {
		t.Errorf("expected default ReadTimeout %v, got %v",
			def.ReadTimeout, cfg.ReadTimeout)
	}
}
