package implementation

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/repository"
	"social/internal/service"
)

func (s *postService) GetByID(ctx context.Context, id string) (model.Post, error) {
	parsed, err := uuid.Parse(id)
	if err != nil || parsed == uuid.Nil {
		return model.Post{}, service.ErrInvalidPostID
	}
	post, err := s.repo.GetByID(ctx, parsed.String())
	if errors.Is(err, repository.ErrPostNotFound) {
		return model.Post{}, service.ErrPostNotFound
	}
	return post, err
}
