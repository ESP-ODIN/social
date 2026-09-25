package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"social/internal/dto"
	"social/internal/service"
)

// CreatePost godoc
//
//	@Summary		Create a post
//	@Description	Creates a post authored by the user identified by the JWT `id` claim.
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			post	body		dto.CreatePostRequest	true	"Post to create"
//	@Success		201		{object}	model.Post
//	@Header			201		{string}	Location	"URL of the created post"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		413		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/api/v1/posts [post]
func CreatePost(posts service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity, ok := authkit.IdentityFromContext(r.Context())
		if !ok {
			writePostError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		var req dto.CreatePostRequest
		if err := decoder.Decode(&req); err != nil {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				writePostError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writePostError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}
		if err := decoder.Decode(new(any)); err != io.EOF {
			var tooLarge *http.MaxBytesError
			if errors.As(err, &tooLarge) {
				writePostError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writePostError(w, http.StatusBadRequest, "request body must contain one JSON object")
			return
		}

		post, err := posts.Create(r.Context(), identity.ID, req)
		if errors.Is(err, service.ErrInvalidPost) {
			writePostError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err != nil {
			slog.Error("create post failed", "error", err)
			writePostError(w, http.StatusInternalServerError, "could not create post")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Location", "/api/v1/posts/"+post.ID)
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(post)
	}
}

func writePostError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}
