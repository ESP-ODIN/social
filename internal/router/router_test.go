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
	"social/internal/model"
	serviceimpl "social/internal/service/implementation"
	"social/internal/testutil"
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
	a := testutil.NewAuth(t)
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
			New(fakePinger{err: tc.pingErr}, nil, a.Authenticator).ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
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
	a := testutil.NewAuth(t)
	rec := httptest.NewRecorder()
	New(fakePinger{}, nil, a.Authenticator).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/posts", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rec.Code)
	}
}

func TestPostRouteUsesJWTAuthor(t *testing.T) {
	authorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	a := testutil.NewAuth(t)
	token := a.Token(t, a.Claims(authorID))

	posts := &recordingPostService{}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(`{"content":"hello"}`))
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	routes := New(fakePinger{}, posts, a.Authenticator)
	routes.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated || !posts.called || posts.authorID != authorID || posts.request.Content != "hello" {
		t.Fatalf("status = %d, called = %v, author = %s, request = %+v", rec.Code, posts.called, posts.authorID, posts.request)
	}
	if got := rec.Header().Get("Location"); got != "/api/v1/posts/post-id" {
		t.Fatalf("Location = %q", got)
	}
}

func TestSwaggerRoutes(t *testing.T) {
	handler := New(fakePinger{}, nil, "test-secret")
	for _, path := range []string{"/swagger/index.html", "/swagger/doc.json"} {
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status = %d, want 200", path, rec.Code)
		}
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	var spec struct {
		Paths map[string]json.RawMessage `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatalf("invalid spec: %v", err)
	}
	for _, path := range []string{"/health", "/ready", "/api/v1/posts"} {
		if _, ok := spec.Paths[path]; !ok {
			t.Fatalf("spec is missing %s; regenerate it with `go tool swag init -g cmd/api/main.go -o docs --parseInternal`", path)
		}
func TestPostAuthenticationFailures(t *testing.T) {
	a := testutil.NewAuth(t)
	id := uuid.New()
	expired := a.Claims(id)
	expired.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))
	noExpiry := a.Claims(id)
	noExpiry.ExpiresAt = nil
	issuer := a.Claims(id)
	issuer.Issuer = "wrong"
	audience := a.Claims(id)
	audience.Audience = jwt.ClaimStrings{"wrong"}
	wrongSigner := testutil.NewAuth(t)
	hs, err := jwt.NewWithClaims(jwt.SigningMethodHS256, a.Claims(id)).SignedString([]byte("obsolete-secret"))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ name, header string }{
		{"missing", ""}, {"malformed", "Basic abc"}, {"empty bearer", "Bearer "}, {"invalid", "Bearer invalid"},
		{"expired", "Bearer " + a.Token(t, expired)}, {"missing expiration", "Bearer " + a.Token(t, noExpiry)},
		{"wrong issuer", "Bearer " + a.Token(t, issuer)}, {"wrong audience", "Bearer " + a.Token(t, audience)},
		{"wrong signature", "Bearer " + wrongSigner.Token(t, a.Claims(id))},
		{"missing id", "Bearer " + a.Token(t, a.Claims(uuid.Nil))}, {"old HS256", "Bearer " + hs},
	} {
		t.Run(tc.name, func(t *testing.T) {
			posts := &recordingPostService{}
			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader("{}"))
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			New(fakePinger{}, posts, a.Authenticator).ServeHTTP(rec, req)
			if rec.Code != http.StatusUnauthorized || posts.called {
				t.Fatalf("status = %d, called = %v", rec.Code, posts.called)
			}
			if rec.Body.String() != "Unauthorized\n" {
				t.Fatalf("unexpected authentication response: %q", rec.Body.String())
			}
		})
	}
}

type recordingRepository struct{ posts []model.Post }

func (r *recordingRepository) Create(_ context.Context, post model.Post) (model.Post, error) {
	r.posts = append(r.posts, post)
	post.ID = uuid.NewString()
	return post, nil
}

func TestAuthenticatedCreationPersistsIdentityAndReusesJWKS(t *testing.T) {
	a := testutil.NewAuth(t)
	repo := &recordingRepository{}
	routes := New(fakePinger{}, serviceimpl.NewPostService(repo), a.Authenticator)
	for i := 0; i < 2; i++ {
		id := uuid.New()
		claims := a.Claims(id)
		claims.Email = ""
		req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader("{\"content\":\" hello \"}"))
		req.Header.Set("Authorization", "Bearer "+a.Token(t, claims))
		rec := httptest.NewRecorder()
		routes.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("status = %d: %s", rec.Code, rec.Body)
		}
		var post model.Post
		if err := json.Unmarshal(rec.Body.Bytes(), &post); err != nil {
			t.Fatal(err)
		}
		if len(repo.posts) != i+1 || repo.posts[i].AuthorID != id.String() || post.AuthorID != id.String() || post.Content != "hello" {
			t.Fatalf("identity not persisted: %+v", post)
		}
	}
	if got := a.Fetches.Load(); got != 1 {
		t.Fatalf("JWKS requests = %d, want 1", got)
	}
}
