package implementation

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"

	"social/internal/dto"
	"social/internal/model"
	"social/internal/service"
)

type fakePostRepository struct {
	called bool
	post   model.Post
}

func (r *fakePostRepository) Create(_ context.Context, post model.Post) (model.Post, error) {
	r.called = true
	r.post = post
	return post, nil
}

func TestCreatePostValidation(t *testing.T) {
	authorID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	tooLongURL := strings.Repeat("a", 256)
	for _, tc := range []struct {
		name     string
		authorID uuid.UUID
		req      dto.CreatePostRequest
	}{
		{"missing author", uuid.Nil, dto.CreatePostRequest{Content: "hello"}},
		{"empty content", authorID, dto.CreatePostRequest{Content: "  "}},
		{"long media url", authorID, dto.CreatePostRequest{Content: "hello", MediaURL: &tooLongURL}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &fakePostRepository{}
			_, err := NewPostService(repo).Create(context.Background(), tc.authorID, tc.req)
			if !errors.Is(err, service.ErrInvalidPost) || repo.called {
				t.Fatalf("err = %v, repository called = %v", err, repo.called)
			}
		})
	}

	repo := &fakePostRepository{}
	_, err := NewPostService(repo).Create(context.Background(), authorID, dto.CreatePostRequest{Content: " hello "})
	if err != nil || !repo.called || repo.post.Content != "hello" || repo.post.AuthorID != authorID.String() {
		t.Fatalf("err = %v, repository called = %v, post = %+v", err, repo.called, repo.post)
	}
}
