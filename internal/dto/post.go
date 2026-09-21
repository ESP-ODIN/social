package dto

type CreatePostRequest struct {
	Content  string  `json:"content"`
	MediaURL *string `json:"media_url,omitempty"`
}
