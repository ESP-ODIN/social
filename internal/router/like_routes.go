package router

import (
	"net/http"

	"github.com/ESP-ODIN/authkit-go"

	"social/internal/handler"
	"social/internal/service"
)

func registerLikeRoutes(mux *http.ServeMux, likes service.LikeService, authenticator *authkit.Authenticator) {
	mux.Handle("POST /api/v1/posts/{postId}/likes", authenticator.RequireAuth(handler.CreateLike(likes)))
}
