package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	HTTPAddr     string
	DatabaseURL  string
	AuthJWKSURL  string
	AuthIssuer   string
	AuthAudience string
}

func Load() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid HTTP_ADDR %q: %w", addr, err)
	}
	n, err := strconv.Atoi(port)
	if err != nil || n < 1 || n > 65535 {
		return Config{}, fmt.Errorf("HTTP_ADDR port must be between 1 and 65535")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}
	cfg := Config{HTTPAddr: addr, DatabaseURL: dbURL, AuthJWKSURL: os.Getenv("AUTH_JWKS_URL"), AuthIssuer: os.Getenv("AUTH_ISSUER"), AuthAudience: os.Getenv("AUTH_AUDIENCE")}
	for _, entry := range []struct{ name, value string }{
		{"AUTH_JWKS_URL", cfg.AuthJWKSURL}, {"AUTH_ISSUER", cfg.AuthIssuer}, {"AUTH_AUDIENCE", cfg.AuthAudience},
	} {
		if strings.TrimSpace(entry.value) == "" {
			return Config{}, fmt.Errorf("%s is required", entry.name)
		}
	}
	return cfg, nil
}
