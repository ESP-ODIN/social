package implementation

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"social/internal/model"
	"social/internal/repository"
)

func (r *postgresPostRepository) GetByID(ctx context.Context, id string) (model.Post, error) {
	var post model.Post
	err := r.db.QueryRow(ctx, `
		SELECT id, author_id, content, media_url, created_at, updated_at
		FROM post
		WHERE id = $1
	`, id).Scan(&post.ID, &post.AuthorID, &post.Content, &post.MediaURL, &post.CreatedAt, &post.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Post{}, repository.ErrPostNotFound
	}
	return post, err
}
