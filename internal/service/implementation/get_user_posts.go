package implementation

import (
	"context"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/service"
)

func (s *postService) GetByUserID(ctx context.Context, userID string) ([]model.Post, error) {
	parsed, err := uuid.Parse(userID)
	if err != nil || parsed == uuid.Nil {
		return nil, service.ErrInvalidUserID
	}
	return s.repo.GetByAuthorID(ctx, parsed.String())
}
