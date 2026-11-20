package router

import (
	"net/http"

	"github.com/ESP-ODIN/authkit-go"

	"social/internal/handler"
	"social/internal/service"
)

func registerGetUserPostsRoutes(mux *http.ServeMux, posts service.PostService, authenticator *authkit.Authenticator) {
	mux.Handle("GET /api/v1/users/{userId}/posts", authenticator.RequireAuth(handler.GetUserPosts(posts)))
}
