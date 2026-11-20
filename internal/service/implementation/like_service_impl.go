package implementation

import (
	"social/internal/repository"
	"social/internal/service"
)

type likeService struct {
	repo repository.LikeRepository
}

func NewLikeService(repo repository.LikeRepository) service.LikeService {
	return &likeService{repo: repo}
}
