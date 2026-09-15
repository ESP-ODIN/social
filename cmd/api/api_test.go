package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/health", http.StatusOK},
		{http.MethodPost, "/health", http.StatusMethodNotAllowed},
		{http.MethodGet, "/unknown", http.StatusNotFound},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			rec := httptest.NewRecorder()
			routes().ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d", rec.Code, tc.status)
			}
			if tc.status == http.StatusOK {
				if got := rec.Header().Get("Content-Type"); got != "application/json" {
					t.Fatalf("Content-Type = %q", got)
				}
				var body struct {
					Status string `json:"status"`
				}
				if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Status != "ok" {
					t.Fatalf("invalid health response: %s (error: %v)", rec.Body, err)
				}
			}
		})
	}
}
