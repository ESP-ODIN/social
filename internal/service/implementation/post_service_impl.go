package implementation

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"social/internal/dto"
	"social/internal/model"
	"social/internal/repository"
	"social/internal/service"
)

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) service.PostService {
	return &postService{repo: repo}
}

func (s *postService) Create(ctx context.Context, authorID uuid.UUID, req dto.CreatePostRequest) (model.Post, error) {
	if authorID == uuid.Nil || strings.TrimSpace(req.Content) == "" || (req.MediaURL != nil && len(*req.MediaURL) > 255) {
		return model.Post{}, service.ErrInvalidPost
	}
	post := model.Post{
		AuthorID: authorID.String(),
		Content:  strings.TrimSpace(req.Content),
		MediaURL: req.MediaURL,
	}
	return s.repo.Create(ctx, post)
}
