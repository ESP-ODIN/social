package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"social/internal/dto"
	"social/internal/model"
)

var ErrInvalidPost = errors.New("author_id must not be empty, content must not be empty, and media_url must be at most 255 characters")
var ErrInvalidPostID = errors.New("post id must be a non-zero UUID")
var ErrPostNotFound = errors.New("post not found")

type PostService interface {
	Create(context.Context, uuid.UUID, dto.CreatePostRequest) (model.Post, error)
	GetByID(context.Context, string) (model.Post, error)
	Delete(context.Context, string, uuid.UUID) error
}
