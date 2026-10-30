package implementation

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"social/internal/repository"
)

func (r *postgresPostRepository) Delete(ctx context.Context, id, authorID string) error {
	var deletedID string
	err := r.db.QueryRow(ctx, `
		DELETE FROM post
		WHERE id = $1 AND author_id = $2
		RETURNING id
	`, id, authorID).Scan(&deletedID)
	if errors.Is(err, pgx.ErrNoRows) {
		return repository.ErrPostNotFound
	}
	return err
}
