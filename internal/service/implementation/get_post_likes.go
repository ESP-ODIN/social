package implementation

import (
	"context"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/service"
)

func (s *likeService) GetByPostID(ctx context.Context, postID string) ([]model.Like, error) {
	parsed, err := uuid.Parse(postID)
	if err != nil || parsed == uuid.Nil {
		return nil, service.ErrInvalidPostID
	}
	return s.repo.GetByPostID(ctx, parsed.String())
}
