package implementation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"social/internal/model"
	"social/internal/repository"
)

func (r *postgresLikeRepository) Create(ctx context.Context, like model.Like) (model.Like, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO "like" (post_id, user_id)
		VALUES ($1, $2)
		RETURNING id, post_id, user_id, created_at
	`, like.PostID, like.UserID).Scan(
		&like.ID, &like.PostID, &like.UserID, &like.CreatedAt,
	)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return model.Like{}, repository.ErrAlreadyLiked
	}
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return model.Like{}, repository.ErrPostNotFound
	}
	return like, err
}
