package implementation

import (
	"context"

	"social/internal/model"
)

func (r *postgresPostRepository) GetByAuthorID(ctx context.Context, authorID string) ([]model.Post, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, author_id, content, media_url, created_at, updated_at
		FROM post
		WHERE author_id = $1
		ORDER BY created_at DESC
	`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []model.Post
	for rows.Next() {
		var p model.Post
		if err := rows.Scan(&p.ID, &p.AuthorID, &p.Content, &p.MediaURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
