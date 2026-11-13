package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"social/internal/dto"
	"social/internal/service"
	"social/internal/testutil"
)

type fakeDeletePostService struct {
	service.PostService
	id       string
	authorID uuid.UUID
	called   bool
	err      error
}

func (s *fakeDeletePostService) Delete(_ context.Context, id string, authorID uuid.UUID) error {
	s.id, s.authorID, s.called = id, authorID, true
	return s.err
}

func (s *fakePostService) Delete(context.Context, string, uuid.UUID) error {
	return errors.New("unexpected Delete call in creation test")
}

func TestDeletePost(t *testing.T) {
	a := testutil.NewAuth(t)
	authorID, id := uuid.New(), uuid.NewString()
	token := a.Token(t, a.Claims(authorID))
	for _, tc := range []struct {
		name    string
		err     error
		status  int
		message string
	}{
		{"deleted", nil, 204, ""},
		{"invalid id", service.ErrInvalidPostID, 400, service.ErrInvalidPostID.Error()},
		{"other author", service.ErrPostForbidden, 403, service.ErrPostForbidden.Error()},
		{"missing post", service.ErrPostNotFound, 404, service.ErrPostNotFound.Error()},
		{"storage failure", errors.New("private SQL details"), 500, "could not delete post"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := &fakeDeletePostService{err: tc.err}
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/posts/"+id, nil)
			req.SetPathValue("id", id)
			req.Header.Set("Authorization", "Bearer "+token)
			rec := httptest.NewRecorder()
			a.Authenticator.RequireAuth(DeletePost(svc)).ServeHTTP(rec, req)
			if rec.Code != tc.status || !svc.called || svc.id != id || svc.authorID != authorID {
				t.Fatalf("status=%d, service=%+v, body=%s", rec.Code, svc, rec.Body)
			}
			if tc.status == 204 {
				if rec.Body.Len() != 0 {
					t.Fatalf("unexpected response body: %s", rec.Body)
				}
			} else {
				var response dto.ErrorResponse
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil || response.Error != tc.message {
					t.Fatalf("unexpected error: %s", rec.Body)
				}
			}
		})
	}
}

func TestDeletePostMissingIdentity(t *testing.T) {
	svc := &fakeDeletePostService{}
	rec := httptest.NewRecorder()
	DeletePost(svc).ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/api/v1/posts/id", nil))
	if rec.Code != 401 || svc.called {
		t.Fatalf("status=%d, called=%v", rec.Code, svc.called)
	}
}
