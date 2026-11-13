package implementation

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"social/internal/model"
	"social/internal/repository"
	"social/internal/service"
)

func (s *likeService) Create(ctx context.Context, postID string, userID uuid.UUID) (model.Like, error) {
	parsedPostID, err := uuid.Parse(postID)
	if err != nil || parsedPostID == uuid.Nil {
		return model.Like{}, service.ErrInvalidPostID
	}
	if userID == uuid.Nil {
		return model.Like{}, service.ErrPostUnauthorized
	}
	like, err := s.repo.Create(ctx, model.Like{
		PostID: parsedPostID.String(),
		UserID: userID.String(),
	})
	if errors.Is(err, repository.ErrAlreadyLiked) {
		return model.Like{}, service.ErrAlreadyLiked
	}
	if errors.Is(err, repository.ErrPostNotFound) {
		return model.Like{}, service.ErrPostNotFound
	}
	return like, err
}
