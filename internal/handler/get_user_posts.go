package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"social/internal/model"
	"social/internal/service"
)

// GetUserPosts godoc
//
//	@Summary		Get posts by user
//	@Description	Returns all posts authored by the given user, newest first.
//	@Tags			posts
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userId	path		string	true	"User ID"	format(uuid)
//	@Success		200		{array}		model.Post
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/users/{userId}/posts [get]
func GetUserPosts(posts service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := posts.GetByUserID(r.Context(), r.PathValue("userId"))
		switch {
		case errors.Is(err, service.ErrInvalidUserID):
			writePostError(w, http.StatusBadRequest, service.ErrInvalidUserID.Error())
		case err != nil:
			slog.Error("get user posts failed", "error", err)
			writePostError(w, http.StatusInternalServerError, "could not get user posts")
		default:
			if result == nil {
				result = []model.Post{}
			}
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(result); err != nil {
				slog.Error("write get user posts response failed", "error", err)
			}
		}
	}
}
