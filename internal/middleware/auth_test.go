package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	sign := func(method jwt.SigningMethod, secret string, id uuid.UUID, expiry time.Time) string {
		t.Helper()
		token, err := jwt.NewWithClaims(method, Claims{
			ID:               id,
			Email:            "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiry)},
		}).SignedString([]byte(secret))
		if err != nil {
			t.Fatal(err)
		}
		return token
	}
	future := time.Now().Add(time.Hour)
	withoutExpiration, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{ID: id}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	withoutEmail, err := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:               id,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(future)},
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name, header string
		want         int
	}{
		{"missing", "", http.StatusUnauthorized},
		{"malformed authorization", "Basic abc", http.StatusUnauthorized},
		{"empty bearer", "Bearer ", http.StatusUnauthorized},
		{"invalid token", "Bearer token_invalide", http.StatusUnauthorized},
		{"wrong secret", "Bearer " + sign(jwt.SigningMethodHS256, "other-secret", id, future), http.StatusUnauthorized},
		{"expired", "Bearer " + sign(jwt.SigningMethodHS256, "test-secret", id, time.Now().Add(-time.Hour)), http.StatusUnauthorized},
		{"missing expiration", "Bearer " + withoutExpiration, http.StatusUnauthorized},
		{"wrong algorithm", "Bearer " + sign(jwt.SigningMethodHS384, "test-secret", id, future), http.StatusUnauthorized},
		{"missing id", "Bearer " + sign(jwt.SigningMethodHS256, "test-secret", uuid.Nil, future), http.StatusUnauthorized},
		{"without email", "Bearer " + withoutEmail, http.StatusNoContent},
		{"valid", "Bearer " + sign(jwt.SigningMethodHS256, "test-secret", id, future), http.StatusNoContent},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			var gotID uuid.UUID
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				gotID, _ = UserIDFromContext(r.Context())
				w.WriteHeader(http.StatusNoContent)
			})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			Auth("test-secret")(next).ServeHTTP(rec, req)
			if rec.Code != tc.want || called != (tc.want == http.StatusNoContent) {
				t.Fatalf("status = %d, handler called = %v", rec.Code, called)
			}
			if called && gotID != id {
				t.Fatalf("context ID = %s, want %s", gotID, id)
			}
		})
	}
}
