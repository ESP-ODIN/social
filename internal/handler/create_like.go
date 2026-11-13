package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"social/internal/dto"
	"social/internal/service"
)

// CreateLike godoc
//
//	@Summary		Like a post
//	@Description	Adds a like to a post for the authenticated user.
//	@Tags			likes
//	@Produce		json
//	@Security		BearerAuth
//	@Param			postId	path		string	true	"Post ID"	format(uuid)
//	@Success		201		{object}	model.Like
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		404		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/posts/{postId}/likes [post]
func CreateLike(likes service.LikeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := authkit.IdentityFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		like, err := likes.Create(r.Context(), r.PathValue("postId"), identity.ID)
		switch {
		case errors.Is(err, service.ErrInvalidPostID):
			writeLikeError(w, http.StatusBadRequest, service.ErrInvalidPostID.Error())
		case errors.Is(err, service.ErrPostNotFound):
			writeLikeError(w, http.StatusNotFound, service.ErrPostNotFound.Error())
		case errors.Is(err, service.ErrAlreadyLiked):
			writeLikeError(w, http.StatusConflict, service.ErrAlreadyLiked.Error())
		case err != nil:
			slog.Error("create like failed", "error", err)
			writeLikeError(w, http.StatusInternalServerError, "could not like post")
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(like)
		}
	}
}

func writeLikeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}
