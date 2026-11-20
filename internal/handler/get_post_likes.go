package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"social/internal/model"
	"social/internal/service"
)

// GetPostLikes godoc
//
//	@Summary		Get likes on a post
//	@Description	Returns all users who liked a post, newest first.
//	@Tags			likes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postId	path		string	true	"Post ID"	format(uuid)
//	@Success		200		{array}		model.Like
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/posts/{postId}/likes [get]
func GetPostLikes(likes service.LikeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := likes.GetByPostID(r.Context(), r.PathValue("postId"))
		switch {
		case errors.Is(err, service.ErrInvalidPostID):
			writeLikeError(w, http.StatusBadRequest, service.ErrInvalidPostID.Error())
		case err != nil:
			slog.Error("get post likes failed", "error", err)
			writeLikeError(w, http.StatusInternalServerError, "could not get post likes")
		default:
			if result == nil {
				result = []model.Like{}
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(result); err != nil {
				slog.Error("write get post likes response failed", "error", err)
			}
		}
	}
}
