package router

import (
	"github.com/ESP-ODIN/authkit-go"
	"net/http"

	"social/internal/handler"
	"social/internal/service"
)

func New(db handler.Pinger, posts service.PostService, authenticator *authkit.Authenticator) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready(db))
	registerPostRoutes(mux, posts, authenticator)
	return mux
}
