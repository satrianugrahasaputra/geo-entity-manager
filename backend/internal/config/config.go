// Package config loads runtime configuration from environment variables.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config holds all runtime settings for the API server.
type Config struct {
	HTTPPort        int
	DatabaseURL     string
	CORSOrigins     []string
	LogLevel        slog.Level
	SeedData        bool
	ShutdownTimeout time.Duration
}

// LookupFunc matches the signature of os.LookupEnv so tests can inject values.
type LookupFunc func(key string) (string, bool)

// Load reads configuration from the process environment.
func Load() (Config, error) {
	return LoadFrom(os.LookupEnv)
}

// LoadFrom reads configuration using the given lookup function.
// All problems are collected and returned together so misconfiguration is fixed in one pass.
func LoadFrom(lookup LookupFunc) (Config, error) {
	cfg := Config{
		HTTPPort:        8080,
		CORSOrigins:     []string{"http://localhost:5173"},
		LogLevel:        slog.LevelInfo,
		SeedData:        false,
		ShutdownTimeout: 10 * time.Second,
	}

	var errs []error

	if v, ok := nonEmpty(lookup, "HTTP_PORT"); ok {
		port, err := strconv.Atoi(v)
		if err != nil || port < 1 || port > 65535 {
			errs = append(errs, fmt.Errorf("HTTP_PORT must be an integer between 1 and 65535, got %q", v))
		} else {
			cfg.HTTPPort = port
		}
	}

	if v, ok := nonEmpty(lookup, "DATABASE_URL"); ok {
		cfg.DatabaseURL = v
	}

	if v, ok := lookup("CORS_ORIGINS"); ok {
		cfg.CORSOrigins = splitList(v)
	}

	if v, ok := nonEmpty(lookup, "LOG_LEVEL"); ok {
		var level slog.Level
		if err := level.UnmarshalText([]byte(v)); err != nil {
			errs = append(errs, fmt.Errorf("LOG_LEVEL must be one of debug|info|warn|error, got %q", v))
		} else {
			cfg.LogLevel = level
		}
	}

	if v, ok := nonEmpty(lookup, "SEED_DATA"); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("SEED_DATA must be a boolean, got %q", v))
		} else {
			cfg.SeedData = b
		}
	}

	if v, ok := nonEmpty(lookup, "SHUTDOWN_TIMEOUT"); ok {
		d, err := time.ParseDuration(v)
		if err != nil || d <= 0 {
			errs = append(errs, fmt.Errorf("SHUTDOWN_TIMEOUT must be a positive duration (e.g. 10s), got %q", v))
		} else {
			cfg.ShutdownTimeout = d
		}
	}

	if len(errs) > 0 {
		return Config{}, fmt.Errorf("invalid configuration: %w", errors.Join(errs...))
	}
	return cfg, nil
}

func nonEmpty(lookup LookupFunc, key string) (string, bool) {
	v, ok := lookup(key)
	v = strings.TrimSpace(v)
	return v, ok && v != ""
}

func splitList(v string) []string {
	out := []string{}
	for _, part := range strings.Split(v, ",") {
		if p := strings.TrimSpace(part); p != "" {
			out = append(out, p)
		}
	}
	return out
}
