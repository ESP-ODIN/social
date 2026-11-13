package implementation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"social/internal/repository"
)

func (r *postgresLikeRepository) Delete(ctx context.Context, postID, userID string) error {
	var deletedID string
	err := r.db.QueryRow(ctx, `
		DELETE FROM "like"
		WHERE post_id = $1 AND user_id = $2
		RETURNING id
	`, postID, userID).Scan(&deletedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrLikeNotFound
	}
	return err
}
