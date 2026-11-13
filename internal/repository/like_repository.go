package repository

import (
	"context"
	"errors"

	"social/internal/model"
)

var ErrAlreadyLiked = errors.New("user already liked this post")

var ErrLikeNotFound = errors.New("like not found")

type LikeRepository interface {
	Create(context.Context, model.Like) (model.Like, error)
	Delete(ctx context.Context, postID, userID string) error
}
