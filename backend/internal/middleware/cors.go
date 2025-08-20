package middleware

import (
	"net/http"

	"github.com/GkadyrG/L0/backend/config"
	"github.com/go-chi/cors"
)

// CORS builds a chi-compatible middleware configured from app config.
// If cfg.Cors.Enabled is false, it returns a no-op middleware.
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	if cfg == nil || !cfg.Cors.Enabled {
		return func(next http.Handler) http.Handler { return next }
	}
	return cors.Handler(cors.Options{
		AllowedOrigins: cfg.Cors.AllowedOrigins,
		AllowedMethods: cfg.Cors.AllowedMethods,
		AllowedHeaders: cfg.Cors.AllowedHeaders,
	})
}
