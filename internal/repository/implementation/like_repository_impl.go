package implementation

import (
	"social/internal/repository"
)

type postgresLikeRepository struct {
	db queryRower
}

func NewLikeRepository(db queryRower) repository.LikeRepository {
	return &postgresLikeRepository{db: db}
}
