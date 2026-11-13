package implementation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"social/internal/repository"
	"social/internal/service"
)

func (s *likeService) Delete(ctx context.Context, postID string, userID uuid.UUID) error {
	parsedPostID, err := uuid.Parse(postID)
	if err != nil || parsedPostID == uuid.Nil {
		return service.ErrInvalidPostID
	}
	if userID == uuid.Nil {
		return service.ErrPostUnauthorized
	}
	err = s.repo.Delete(ctx, parsedPostID.String(), userID.String())
	if errors.Is(err, repository.ErrLikeNotFound) {
		return service.ErrLikeNotFound
	}
	return err
}
