package router

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"social/internal/model"
	"social/internal/repository"
	serviceimpl "social/internal/service/implementation"
	"social/internal/testutil"
	"testing"
)

func TestGetPostRoute(t *testing.T) {
	a := testutil.NewAuth(t)
	id := uuid.NewString()
	repo := &recordingRepository{posts: []model.Post{{ID: id, AuthorID: uuid.NewString(), Content: "hello"}}}
	routes := New(fakePinger{}, serviceimpl.NewPostService(repo), a.Authenticator)
	token := a.Token(t, a.Claims(uuid.New()))
	for _, tc := range []struct {
		name, id, token string
		status          int
	}{
		{"other user can read", id, token, 200},
		{"missing token", id, "", 401},
		{"invalid token", id, "invalid", 401},
		{"invalid id", "bad-id", token, 400},
		{"zero id", uuid.Nil.String(), token, 400},
		{"missing post", uuid.NewString(), token, 404},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/posts/"+tc.id, nil)
			if tc.token != "" {
				req.Header.Set("Authorization", "Bearer "+tc.token)
			}
			rec := httptest.NewRecorder()
			routes.ServeHTTP(rec, req)
			if rec.Code != tc.status {
				t.Fatalf("status=%d, body=%s", rec.Code, rec.Body)
			}
			if tc.status == 200 {
				var post model.Post
				if err := json.Unmarshal(rec.Body.Bytes(), &post); err != nil || post != repo.posts[0] {
					t.Fatalf("unexpected post: %s", rec.Body)
				}
			}
		})
	}
}

func (s *recordingPostService) GetByID(_ context.Context, id string) (model.Post, error) {
	s.called = true
	return model.Post{ID: id}, nil
}

func (r *recordingRepository) GetByID(_ context.Context, id string) (model.Post, error) {
	for _, post := range r.posts {
		if post.ID == id {
			return post, nil
		}
	}
	return model.Post{}, repository.ErrPostNotFound
}
