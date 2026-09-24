package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"social/internal/dto"
)

// Pinger is implemented by the database pool.
type Pinger interface {
	Ping(context.Context) error
}

// Health godoc
//
//	@Summary	Liveness check
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	dto.StatusResponse
//	@Router		/health [get]
func Health(w http.ResponseWriter, _ *http.Request) {
	writeStatus(w, http.StatusOK, "ok")
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Pings the database on every request.
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	dto.StatusResponse
//	@Failure		503	{object}	dto.StatusResponse
//	@Router			/ready [get]
func Ready(db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		w.Header().Set("Cache-Control", "no-store")
		if err := db.Ping(ctx); err != nil {
			slog.Warn("readiness check failed", "error", err)
			writeStatus(w, http.StatusServiceUnavailable, "unavailable")
			return
		}
		writeStatus(w, http.StatusOK, "ok")
	}
}

func writeStatus(w http.ResponseWriter, status int, value string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(dto.StatusResponse{Status: value})
}
