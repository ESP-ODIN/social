package repository

import (
	"context"

	"social/internal/model"
)

var ErrPostNotFound = model.ErrPostNotFound

type PostRepository interface {
	Create(context.Context, model.Post) (model.Post, error)
	GetByID(context.Context, string) (model.Post, error)
	Delete(context.Context, string, string) error
}
