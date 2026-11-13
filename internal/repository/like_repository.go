package repository

import (
	"context"
	"errors"

	"social/internal/model"
)

var ErrAlreadyLiked = errors.New("user already liked this post")

type LikeRepository interface {
	Create(context.Context, model.Like) (model.Like, error)
}
