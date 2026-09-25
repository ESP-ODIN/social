package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"social/internal/dto"
	"social/internal/model"
	"social/internal/service"
	"social/internal/testutil"
)

type fakePostService struct {
	post     model.Post
	err      error
	authorID uuid.UUID
	called   bool
}

func (s *fakePostService) Create(_ context.Context, authorID uuid.UUID, _ dto.CreatePostRequest) (model.Post, error) {
	s.authorID = authorID
	s.called = true
	return s.post, s.err
}

func TestCreatePost(t *testing.T) {
	const id = "550e8400-e29b-41d4-a716-446655440000"
	a := testutil.NewAuth(t)
	token := a.Token(t, a.Claims(uuid.MustParse(id)))
	for _, tc := range []struct {
		name, body string
		service    fakePostService
		status     int
	}{
		{"created", `{"content":"hello"}`, fakePostService{post: model.Post{ID: id, Content: "hello"}}, http.StatusCreated},
		{"client author rejected", `{"author_id":"7e733b0f-136e-4cd9-923d-19c733cd466a","content":"hello"}`, fakePostService{}, http.StatusBadRequest},
		{"malformed JSON", `{`, fakePostService{}, http.StatusBadRequest},
		{"unknown field", `{"unknown":1}`, fakePostService{}, http.StatusBadRequest},
		{"invalid post", `{}`, fakePostService{err: service.ErrInvalidPost}, http.StatusBadRequest},
		{"storage error", `{}`, fakePostService{err: errors.New("database down")}, http.StatusInternalServerError},
		{"too large", `{"content":"` + strings.Repeat("x", (1<<20)+1) + `"}`, fakePostService{}, http.StatusRequestEntityTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader(tc.body))
			req.Header.Set("Authorization", "Bearer "+token)
			svc := tc.service
			a.Authenticator.RequireAuth(CreatePost(&svc)).ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status = %d, want %d; body = %s", rec.Code, tc.status, rec.Body)
			}
			if tc.status == http.StatusCreated {
				if !svc.called || svc.authorID != uuid.MustParse(id) {
					t.Fatalf("service author = %s, called = %v", svc.authorID, svc.called)
				}
				if got := rec.Header().Get("Location"); got != "/api/v1/posts/"+id {
					t.Fatalf("Location = %q", got)
				}
				var post model.Post
				if err := json.Unmarshal(rec.Body.Bytes(), &post); err != nil || post.ID != id {
					t.Fatalf("invalid post response: %s (%v)", rec.Body, err)
				}
			}
			if tc.name == "client author rejected" && svc.called {
				t.Fatal("service called with client-supplied author_id")
			}
		})
	}
}

func TestCreatePostMissingIdentity(t *testing.T) {
	svc := &fakePostService{}
	rec := httptest.NewRecorder()
	CreatePost(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/posts", strings.NewReader("{")))
	if rec.Code != http.StatusUnauthorized || svc.called {
		t.Fatalf("status = %d, called = %v", rec.Code, svc.called)
	}
}
