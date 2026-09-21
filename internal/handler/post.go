package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"social/internal/dto"
	"social/internal/middleware"
	"social/internal/service"
)

func CreatePost(posts service.PostService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorID, ok := middleware.UserIDFromContext(r.Context())
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

		post, err := posts.Create(r.Context(), authorID, req)
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
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}
