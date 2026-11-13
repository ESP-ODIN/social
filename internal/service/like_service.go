package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"social/internal/model"
)

var ErrAlreadyLiked = errors.New("user already liked this post")
var ErrLikeNotFound = errors.New("like not found")

type LikeService interface {
	Create(ctx context.Context, postID string, userID uuid.UUID) (model.Like, error)
	Delete(ctx context.Context, postID string, userID uuid.UUID) error
}
