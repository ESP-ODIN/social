package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"social/internal/dto"
	"social/internal/middleware"
	"social/internal/model"
)

type fakePinger struct{ err error }

func (p fakePinger) Ping(context.Context) error { return p.err }

type recordingPostService struct {
	authorID uuid.UUID
	request  dto.CreatePostRequest
	called   bool
}

func (s *recordingPostService) Create(_ context.Context, authorID uuid.UUID, request dto.CreatePostRequest) (model.Post, error) {
	s.authorID = authorID
	s.request = request
	s.called = true
	return model.Post{ID: "post-id", AuthorID: authorID.String(), Content: request.Content}, nil
}

func TestRoutes(t *testing.T) {
	for _, tc := range []struct {
		name, method, path string
		pingErr            error
		status             int
		bodyStatus         string
	}{
		{"health", http.MethodGet, "/health", nil, http.StatusOK, "ok"},
		{"health method", http.MethodPost, "/health", nil, http.StatusMethodNotAllowed, ""},
		{"ready", http.MethodGet, "/ready", nil, http.StatusOK, "ok"},
		{"database unavailable", http.MethodGet, "/ready", errors.New("database unavailable"), http.StatusServiceUnavailable, "unavailable"},
		{"ready method", http.MethodPost, "/ready", nil, http.StatusMethodNotAllowed, ""},
		{"unknown", http.MethodGet, "/unknown", nil, http.StatusNotFound, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			New(fakePinger{err: tc.pingErr}, nil, "test-secret").ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d", rec.Code, tc.status)
			}
			if tc.bodyStatus == "" {
				return
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q", got)
			}
			var body struct {
				Status string `json:"status"`
			}
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Status != tc.bodyStatus {
				t.Fatalf("invalid response: %s (error: %v)", rec.Body, err)
			}
		})
	}
}

func TestPostRequiresAuthentication(t *testing.T) {
	rec := httptest.NewRecorder()
	New(fakePinger{}, nil, "test-secret").ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/posts", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestPostRouteUsesJWTAuthor(t *testing.T) {
	authorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.Claims{
		ID:    authorID,
		Email: "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}).SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatal(err)
	}

	posts := &recordingPostService{}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(`{"content":"hello"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	New(fakePinger{}, posts, "test-secret").ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated || !posts.called || posts.authorID != authorID || posts.request.Content != "hello" {
		t.Fatalf("status = %d, called = %v, author = %s, request = %+v", rec.Code, posts.called, posts.authorID, posts.request)
	}
	if got := rec.Header().Get("Location"); got != "/api/v1/posts/post-id" {
		t.Fatalf("Location = %q", got)
	}
}
