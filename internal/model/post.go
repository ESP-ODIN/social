package model

import "time"

type Post struct {
	ID        string    `json:"id" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	AuthorID  string    `json:"author_id" format:"uuid" example:"7e733b0f-136e-4cd9-923d-19c733cd466a"`
	Content   string    `json:"content" example:"Hello"`
	MediaURL  *string   `json:"media_url" example:"https://example.com/image.jpg"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
}
