package router

import (
	"github.com/ESP-ODIN/authkit-go"
	"net/http"
	"social/internal/handler"
	"social/internal/service"
)

func registerGetPostRoutes(mux *http.ServeMux, posts service.PostService, authenticator *authkit.Authenticator) {
	mux.Handle("GET /api/v1/posts/{id}", authenticator.RequireAuth(handler.GetPost(posts)))
}
