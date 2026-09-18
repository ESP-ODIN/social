package repository

import (
	"context"

	"social/internal/model"
)

type PostRepository interface {
	Create(context.Context, model.Post) (model.Post, error)
}
