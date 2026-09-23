package config

import "testing"

func TestLoad(t *testing.T) {
	for _, tc := range []struct {
		value, want string
		invalid     bool
	}{
		{"", "127.0.0.1:8080", false},
		{"127.0.0.1:9090", "127.0.0.1:9090", false},
		{"[::1]:8080", "[::1]:8080", false},
		{"localhost", "", true},
		{"localhost:abc", "", true},
		{"localhost:0", "", true},
		{"localhost:65536", "", true},
	} {
		t.Run(tc.value, func(t *testing.T) {
			t.Setenv("HTTP_ADDR", tc.value)
			t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
			setAuthEnv(t)
			cfg, err := Load()
			if (err != nil) != tc.invalid {
				t.Fatalf("Load error = %v, invalid = %v", err, tc.invalid)
			}
			if !tc.invalid && cfg.HTTPAddr != tc.want {
				t.Fatalf("address = %q, want %q", cfg.HTTPAddr, tc.want)
			}
			if !tc.invalid && (cfg.AuthJWKSURL != "http://localhost:8080/.well-known/jwks.json" || cfg.AuthIssuer != "http://localhost:8080" || cfg.AuthAudience != "odin-api") {
				t.Fatalf("auth configuration was not loaded")
			}
		})
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	t.Setenv("HTTP_ADDR", "127.0.0.1:8080")
	t.Setenv("DATABASE_URL", "")
	setAuthEnv(t)
	if _, err := Load(); err == nil {
		t.Fatal("expected error when DATABASE_URL is missing")
	}
}

func setAuthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("AUTH_JWKS_URL", "http://localhost:8080/.well-known/jwks.json")
	t.Setenv("AUTH_ISSUER", "http://localhost:8080")
	t.Setenv("AUTH_AUDIENCE", "odin-api")
}
func TestLoad_MissingAuthConfig(t *testing.T) {
	for _, name := range []string{"AUTH_JWKS_URL", "AUTH_ISSUER", "AUTH_AUDIENCE"} {
		for _, value := range []string{"", "  "} {
			t.Run(name+value, func(t *testing.T) {
				t.Setenv("HTTP_ADDR", "127.0.0.1:8081")
				t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
				setAuthEnv(t)
				t.Setenv(name, value)
				if _, err := Load(); err == nil || err.Error() != name+" is required" {
					t.Fatalf("unexpected error: %v", err)
				}
			})
		}
	}
}
