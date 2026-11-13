package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"social/internal/service"
)

// DeleteLike godoc
//
//	@Summary		Unlike a post
//	@Description	Removes the authenticated user's like from a post.
//	@Tags			likes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postId	path	string	true	"Post ID"	format(uuid)
//	@Success		204		"Like removed"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/posts/{postId}/likes [delete]
func DeleteLike(likes service.LikeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := authkit.IdentityFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		err := likes.Delete(r.Context(), r.PathValue("postId"), identity.ID)
		switch {
		case errors.Is(err, service.ErrInvalidPostID):
			writeLikeError(w, http.StatusBadRequest, service.ErrInvalidPostID.Error())
		case errors.Is(err, service.ErrLikeNotFound):
			writeLikeError(w, http.StatusNotFound, service.ErrLikeNotFound.Error())
		case err != nil:
			slog.Error("delete like failed", "error", err)
			writeLikeError(w, http.StatusInternalServerError, "could not unlike post")
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
