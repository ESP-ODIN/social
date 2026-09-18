package implementation

import (
	"context"

	"github.com/jackc/pgx/v5"
	"social/internal/model"
	"social/internal/repository"
)

type queryRower interface {
	QueryRow(context.Context, string, ...any) pgx.Row
}

type postgresPostRepository struct {
	db queryRower
}

func NewPostRepository(db queryRower) repository.PostRepository {
	return &postgresPostRepository{db: db}
}

func (r *postgresPostRepository) Create(ctx context.Context, post model.Post) (model.Post, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO post (author_id, content, media_url)
		VALUES ($1, $2, $3)
		RETURNING id, author_id, content, media_url, created_at, updated_at
	`, post.AuthorID, post.Content, post.MediaURL).Scan(
		&post.ID, &post.AuthorID, &post.Content, &post.MediaURL, &post.CreatedAt, &post.UpdatedAt,
	)
	return post, err
}
