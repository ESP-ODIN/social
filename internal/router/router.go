package router

import (
	"net/http"

	"social/internal/handler"
	"social/internal/service"
)

func New(db handler.Pinger, posts service.PostService, jwtSecretKey string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	mux.HandleFunc("GET /ready", handler.Ready(db))
	registerPostRoutes(mux, posts, jwtSecretKey)
	registerSwaggerRoutes(mux)
	return mux
}
