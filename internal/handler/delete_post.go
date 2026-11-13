package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"social/internal/service"
)

// DeletePost godoc
//
//	@Summary Delete a post
//	@Description Permanently deletes a post. Only its author can delete it.
//	@Tags posts
//	@Produce json
//	@Security BearerAuth
//	@Param id path string true "Post ID" format(uuid)
//	@Success 204 "Post deleted"
//	@Failure 400 {object} dto.ErrorResponse
//	@Failure 401 {string} string "Unauthorized"
//	@Failure 403 {object} dto.ErrorResponse
//	@Failure 404 {object} dto.ErrorResponse
//	@Failure 500 {object} dto.ErrorResponse
//	@Router /api/v1/posts/{id} [delete]
func DeletePost(posts service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := authkit.IdentityFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		err := posts.Delete(r.Context(), r.PathValue("id"), identity.ID)
		switch {
		case errors.Is(err, service.ErrPostUnauthorized):
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		case errors.Is(err, service.ErrInvalidPostID):
			writePostError(w, http.StatusBadRequest, service.ErrInvalidPostID.Error())
		case errors.Is(err, service.ErrPostForbidden):
			writePostError(w, http.StatusForbidden, service.ErrPostForbidden.Error())
		case errors.Is(err, service.ErrPostNotFound):
			writePostError(w, http.StatusNotFound, service.ErrPostNotFound.Error())
		case err != nil:
			slog.Error("delete post failed", "error", err)
			writePostError(w, http.StatusInternalServerError, "could not delete post")
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
