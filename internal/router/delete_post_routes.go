package router

import (
	"net/http"

	"github.com/ESP-ODIN/authkit-go"
	"social/internal/handler"
	"social/internal/service"
)

func registerDeletePostRoutes(mux *http.ServeMux, posts service.PostService, authenticator *authkit.Authenticator) {
	mux.Handle("DELETE /api/v1/posts/{id}", authenticator.RequireAuth(handler.DeletePost(posts)))
}
