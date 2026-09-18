package router

import (
	"net/http"

	"social/internal/handler"
	"social/internal/middleware"
	"social/internal/service"
)

func registerPostRoutes(mux *http.ServeMux, posts service.PostService, jwtSecretKey string) {
	mux.Handle("POST /api/v1/posts", middleware.Auth(jwtSecretKey)(handler.CreatePost(posts)))
}
