package repository

import (
	"context"
	"errors"

	"social/internal/model"
)

var ErrPostNotFound = errors.New("post not found")

type PostRepository interface {
	Create(context.Context, model.Post) (model.Post, error)
	GetByID(context.Context, string) (model.Post, error)
	Delete(context.Context, string, string) error
}
