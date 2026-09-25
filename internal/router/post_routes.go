package router

import (
	"github.com/ESP-ODIN/authkit-go"
	"net/http"

	"social/internal/handler"

	"social/internal/service"
)

func registerPostRoutes(mux *http.ServeMux, posts service.PostService, authenticator *authkit.Authenticator) {
	mux.Handle("POST /api/v1/posts", authenticator.RequireAuth(handler.CreatePost(posts)))
}
