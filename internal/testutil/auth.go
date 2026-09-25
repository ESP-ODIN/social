// Package testutil provides ephemeral authentication fixtures for Social tests.
package testutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ESP-ODIN/authkit-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Auth struct {
	Authenticator *authkit.Authenticator
	Fetches       atomic.Int32
	key           *rsa.PrivateKey
}

func NewAuth(t *testing.T) *Auth {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	a := &Auth{key: key}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Fetches.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]string{{
			"kty": "RSA", "use": "sig", "alg": "RS256", "kid": "social-test",
			"n": base64.RawURLEncoding.EncodeToString(key.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(key.E)).Bytes()),
		}}})
	}))
	t.Cleanup(server.Close)
	a.Authenticator, err = authkit.New(authkit.Config{JWKSURL: server.URL, Issuer: "social-test-auth", Audience: "social-test-api", HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func (a *Auth) Claims(id uuid.UUID) authkit.Claims {
	return authkit.Claims{ID: id, Email: "test@example.com", RegisteredClaims: jwt.RegisteredClaims{
		Issuer: "social-test-auth", Audience: jwt.ClaimStrings{"social-test-api"},
		// Deliberately different: the author must come from id, never sub.
		Subject: uuid.NewString(), IssuedAt: jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}}
}

func (a *Auth) Token(t *testing.T, claims authkit.Claims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "social-test"
	raw, err := token.SignedString(a.key)
	if err != nil {
		t.Fatal(err)
	}
	return raw
}
