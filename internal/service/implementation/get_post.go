package implementation

import (
	"context"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/service"
)

func (s *postService) GetByID(ctx context.Context, id string) (model.Post, error) {
	parsed, err := uuid.Parse(id)
	if err != nil || parsed == uuid.Nil {
		return model.Post{}, service.ErrInvalidPostID
	}
	return s.repo.GetByID(ctx, parsed.String())
}
