package implementation

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/repository"
	"social/internal/service"
	"testing"
)

func TestGetPostValidationAndErrors(t *testing.T) {
	for _, id := range []string{"", "bad-id", uuid.Nil.String()} {
		repo := &fakeGetPostRepository{}
		_, err := NewPostService(repo).GetByID(context.Background(), id)
		if !errors.Is(err, service.ErrInvalidPostID) || repo.called {
			t.Fatalf("id=%q, err=%v, called=%v", id, err, repo.called)
		}
	}
	id := uuid.NewString()
	storageErr := errors.New("database unavailable")
	for _, tc := range []struct{ input, want error }{
		{nil, nil}, {repository.ErrPostNotFound, service.ErrPostNotFound}, {storageErr, storageErr},
	} {
		repo := &fakeGetPostRepository{post: model.Post{ID: id}, err: tc.input}
		post, err := NewPostService(repo).GetByID(context.Background(), id)
		if !errors.Is(err, tc.want) || !repo.called || repo.id != id || (err == nil && post.ID != id) {
			t.Fatalf("post=%+v, err=%v, repo=%+v", post, err, repo)
		}
	}
}

func (r *fakePostRepository) GetByID(context.Context, string) (model.Post, error) {
	return model.Post{}, errors.New("unexpected GetByID call in creation test")
}

type fakeGetPostRepository struct {
	repository.PostRepository
	id     string
	err    error
	called bool
	post   model.Post
}

func (r *fakeGetPostRepository) GetByID(_ context.Context, id string) (model.Post, error) {
	r.called = true
	r.id = id
	return r.post, r.err
}
