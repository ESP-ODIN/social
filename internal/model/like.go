package model

import "time"

type Like struct {
	ID        string    `json:"id" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	PostID    string    `json:"post_id" format:"uuid" example:"97ba14e4-612b-4a6b-8d0f-ce5d289dfc80"`
	UserID    string    `json:"user_id" format:"uuid" example:"c3c4fb8b-08f6-49fe-931e-e7b1fba62651"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}
