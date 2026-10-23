package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"social/internal/dto"
	"social/internal/model"
	"social/internal/service"
	"social/internal/testutil"
	"testing"
)

func TestGetPost(t *testing.T) {
	const id = "550e8400-e29b-41d4-a716-446655440000"
	a := testutil.NewAuth(t)
	token := a.Token(t, a.Claims(uuid.New()))
	for _, tc := range []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{"found", nil, 200, ""},
		{"invalid id", service.ErrInvalidPostID, 400, service.ErrInvalidPostID.Error()},
		{"not found", service.ErrPostNotFound, 404, service.ErrPostNotFound.Error()},
		{"storage failure", errors.New("private SQL details"), 500, "could not get post"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := &fakeGetPostService{post: model.Post{ID: id, Content: "hello"}, err: tc.err}
			req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+id, nil)
			req.SetPathValue("id", id)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			a.Authenticator.RequireAuth(GetPost(svc)).ServeHTTP(rec, req)
			if rec.Code != tc.status || !svc.called || svc.id != id || rec.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("status=%d, service=%+v, body=%s", rec.Code, svc, rec.Body)
			}
			if tc.status == 200 {
				var post model.Post
				if err := json.Unmarshal(rec.Body.Bytes(), &post); err != nil || post.ID != id || post.Content != "hello" {
					t.Fatalf("invalid response: %s", rec.Body)
				}
			} else {
				var response dto.ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil || response.Error != tc.message {
					t.Fatalf("invalid error: %s", rec.Body)
				}
			}
		})
	}
}

func TestGetPostMissingIdentity(t *testing.T) {
	svc := &fakeGetPostService{}
	rec := httptest.NewRecorder()
	GetPost(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/posts/id", nil))
	if rec.Code != http.StatusUnauthorized || svc.called {
		t.Fatalf("status=%d, called=%v", rec.Code, svc.called)
	}
}

func (s *fakePostService) GetByID(context.Context, string) (model.Post, error) {
	return model.Post{}, errors.New("unexpected GetByID call in creation test")
}

type fakeGetPostService struct {
	service.PostService
	id     string
	post   model.Post
	err    error
	called bool
}

func (s *fakeGetPostService) GetByID(_ context.Context, id string) (model.Post, error) {
	s.called = true
	s.id = id
	return s.post, s.err
}
