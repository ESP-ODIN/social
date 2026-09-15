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
			cfg, err := Load()
			if (err != nil) != tc.invalid {
				t.Fatalf("Load error = %v, invalid = %v", err, tc.invalid)
			}
			if !tc.invalid && cfg.HTTPAddr != tc.want {
				t.Fatalf("address = %q, want %q", cfg.HTTPAddr, tc.want)
			}
		})
	}
}
