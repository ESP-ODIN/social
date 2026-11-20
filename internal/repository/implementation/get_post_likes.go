package implementation

import (
	"context"

	"social/internal/model"
)

func (r *postgresLikeRepository) GetByPostID(ctx context.Context, postID string) ([]model.Like, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, post_id, user_id, created_at
		FROM "like"
		WHERE post_id = $1
		ORDER BY created_at DESC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []model.Like
	for rows.Next() {
		var l model.Like
		if err := rows.Scan(&l.ID, &l.PostID, &l.UserID, &l.CreatedAt); err != nil {
			return nil, err
		}
		likes = append(likes, l)
	}
	return likes, rows.Err()
}
