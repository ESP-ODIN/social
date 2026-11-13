package implementation

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/repository"
	"social/internal/service"
)

type fakeDeletePostRepository struct {
	repository.PostRepository
	post              model.Post
	getErr, deleteErr error
	read, deleted     bool
	id, authorID      string
}

func (r *fakeDeletePostRepository) GetByID(context.Context, string) (model.Post, error) {
	r.read = true
	return r.post, r.getErr
}

func (r *fakeDeletePostRepository) Delete(_ context.Context, id, authorID string) error {
	r.deleted, r.id, r.authorID = true, id, authorID
	return r.deleteErr
}

func (r *fakePostRepository) Delete(context.Context, string, string) error {
	return errors.New("unexpected Delete call in creation test")
}

func TestDeletePost(t *testing.T) {
	authorID, id := uuid.New(), uuid.NewString()
	storageErr := errors.New("database unavailable")
	for _, tc := range []struct {
		name, id                string
		actor                   uuid.UUID
		getErr, deleteErr, want error
		read, deleted           bool
	}{
		{"deleted", id, authorID, nil, nil, nil, false, true},
		{"invalid id", "bad-id", authorID, nil, nil, service.ErrInvalidPostID, false, false},
		{"invalid id and missing identity", "bad-id", uuid.Nil, nil, nil, service.ErrInvalidPostID, false, false},
		{"zero id", uuid.Nil.String(), authorID, nil, nil, service.ErrInvalidPostID, false, false},
		{"missing identity", id, uuid.Nil, nil, nil, service.ErrPostUnauthorized, false, false},
		{"other author", id, uuid.New(), nil, repository.ErrPostNotFound, service.ErrPostForbidden, true, true},
		{"missing post", id, authorID, repository.ErrPostNotFound, repository.ErrPostNotFound, service.ErrPostNotFound, true, true},
		{"read failure", id, authorID, storageErr, repository.ErrPostNotFound, storageErr, true, true},
		{"delete failure", id, authorID, nil, storageErr, storageErr, false, true},
		{"already deleted", id, authorID, repository.ErrPostNotFound, fmt.Errorf("delete: %w", repository.ErrPostNotFound), service.ErrPostNotFound, true, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakeDeletePostRepository{post: model.Post{ID: id, AuthorID: authorID.String()}, getErr: tc.getErr, deleteErr: tc.deleteErr}
			err := NewPostService(repo).Delete(context.Background(), tc.id, tc.actor)
			if !errors.Is(err, tc.want) || repo.read != tc.read || repo.deleted != tc.deleted {
				t.Fatalf("err=%v, repo=%+v", err, repo)
			}
			if tc.deleted && (repo.id != id || repo.authorID != tc.actor.String()) {
				t.Fatalf("unexpected delete arguments: %+v", repo)
			}
		})
	}
}
