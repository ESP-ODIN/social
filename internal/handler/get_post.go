package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"social/internal/service"
)

// GetPost godoc
//
//	@Summary Get a post
//	@Description Returns a post by ID. Any authenticated user can read it.
//	@Tags posts
//	@Produce json
//	@Security BearerAuth
//	@Param id path string true "Post ID" format(uuid)
//	@Success 200 {object} model.Post
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 401 {string} string "Unauthorized"
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /api/v1/posts/{id} [get]
func GetPost(posts service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		post, err := posts.GetByID(r.Context(), r.PathValue("id"))
		switch {
		case errors.Is(err, service.ErrInvalidPostID):
			writePostError(w, http.StatusBadRequest, service.ErrInvalidPostID.Error())
		case errors.Is(err, service.ErrPostNotFound):
			writePostError(w, http.StatusNotFound, service.ErrPostNotFound.Error())
		case err != nil:
			slog.Error("get post failed", "error", err)
			writePostError(w, http.StatusInternalServerError, "could not get post")
		default:
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(post); err != nil {
				slog.Error("write get post response failed", "error", err)
			}
		}
	}
}
