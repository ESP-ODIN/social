package dto

type CreatePostRequest struct {
	Content  string  `json:"content" validate:"required" example:"Hello"`
	MediaURL *string `json:"media_url,omitempty" maxLength:"255" example:"https://example.com/image.jpg"`
}
