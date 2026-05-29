// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds every tunable the server needs to boot.
type Config struct {
	Port            string
	DatabasePath    string
	JWTSecret       string
	JWTTTL          time.Duration
	AllowedOrigin   string
	StripeSecretKey string
	StripeWebhook   string
	FrontendURL     string
}

// Load reads the environment and applies sensible development defaults.
// It only fails when a value is present but malformed.
func Load() (Config, error) {
	ttl, err := parseHours("JWT_TTL_HOURS", 72)
	if err != nil {
		return Config{}, err
	}
	return Config{
		Port:            env("PORT", "8080"),
		DatabasePath:    env("DATABASE_PATH", "lua_academy.db"),
		JWTSecret:       env("JWT_SECRET", "dev-insecure-secret-change-me"),
		JWTTTL:          ttl,
		AllowedOrigin:   env("ALLOWED_ORIGIN", "http://localhost:5173"),
		StripeSecretKey: env("STRIPE_SECRET_KEY", ""),
		StripeWebhook:   env("STRIPE_WEBHOOK_SECRET", ""),
		FrontendURL:     env("FRONTEND_URL", "http://localhost:5173"),
	}, nil
}

// env returns the variable or a fallback when it is unset/empty.
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// parseHours reads an integer hour count and turns it into a Duration.
func parseHours(key string, fallback int) (time.Duration, error) {
	raw := os.Getenv(key)
	if raw == "" {
		return time.Duration(fallback) * time.Hour, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return time.Duration(n) * time.Hour, nil
}
