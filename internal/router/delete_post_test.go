package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/repository"
	serviceimpl "social/internal/service/implementation"
	"social/internal/testutil"
)

func (s *recordingPostService) Delete(context.Context, string, uuid.UUID) error {
	return errors.New("unexpected Delete call")
}

func (r *recordingRepository) Delete(context.Context, string, string) error {
	return errors.New("unexpected Delete call")
}

type deleteRepository struct{ recordingRepository }

func (r *deleteRepository) Delete(_ context.Context, id, authorID string) error {
	for i, post := range r.posts {
		if post.ID == id && post.AuthorID == authorID {
			r.posts = append(r.posts[:i], r.posts[i+1:]...)
			return nil
		}
	}
	return repository.ErrPostNotFound
}

func TestDeletePostRoute(t *testing.T) {
	a := testutil.NewAuth(t)
	authorID, id := uuid.New(), uuid.NewString()
	ownerToken := a.Token(t, a.Claims(authorID))
	otherToken := a.Token(t, a.Claims(uuid.New()))
	for _, tc := range []struct {
		name, id, token string
		status          int
	}{
		{"owner", id, ownerToken, 204},
		{"other user", id, otherToken, 403},
		{"missing token", id, "", 401},
		{"invalid token", id, "invalid", 401},
		{"invalid id", "bad-id", ownerToken, 400},
		{"missing post", uuid.NewString(), ownerToken, 404},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &deleteRepository{recordingRepository{posts: []model.Post{{ID: id, AuthorID: authorID.String()}}}}
			routes := New(fakePinger{}, serviceimpl.NewPostService(repo), a.Authenticator)
			request := func(method, path, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(method, path, nil)
				if token != "" {
					req.Header.Set("Authorization", "Bearer "+token)
				}
				rec := httptest.NewRecorder()
				routes.ServeHTTP(rec, req)
				return rec
			}
			rec := request(http.MethodDelete, "/api/v1/posts/"+tc.id, tc.token)
			if rec.Code != tc.status {
				t.Fatalf("status=%d, body=%s", rec.Code, rec.Body)
			}
			if tc.status == 204 {
				if len(repo.posts) != 0 || rec.Body.Len() != 0 {
					t.Fatal("post was not deleted or response was not empty")
				}
				if got := request(http.MethodGet, "/api/v1/posts/"+id, ownerToken); got.Code != 404 {
					t.Fatalf("GET after delete: %d", got.Code)
				}
				if got := request(http.MethodDelete, "/api/v1/posts/"+id, ownerToken); got.Code != 404 {
					t.Fatalf("second DELETE: %d", got.Code)
				}
			} else if len(repo.posts) != 1 {
				t.Fatal("post deleted despite rejected request")
			}
		})
	}
}

func TestDeletePostSwagger(t *testing.T) {
	a := testutil.NewAuth(t)
	rec := httptest.NewRecorder()
	New(fakePinger{}, nil, a.Authenticator).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil))
	var spec struct {
		Paths map[string]map[string]struct {
			Responses map[string]json.RawMessage `json:"responses"`
			Security  []map[string][]string      `json:"security"`
		} `json:"paths"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &spec); err != nil {
		t.Fatal(err)
	}
	op, ok := spec.Paths["/api/v1/posts/{id}"]["delete"]
	if !ok {
		t.Fatal("DELETE missing from Swagger")
	}
	for _, status := range []string{"204", "400", "401", "403", "404", "500"} {
		if _, ok := op.Responses[status]; !ok {
			t.Fatalf("missing response %s", status)
		}
	}
	if len(op.Security) != 1 {
		t.Fatal("missing security definition")
	}
	if _, ok := op.Security[0]["BearerAuth"]; !ok {
		t.Fatal("missing BearerAuth")
	}
}
