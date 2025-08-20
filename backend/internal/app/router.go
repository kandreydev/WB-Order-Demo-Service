package app

import (
	"github.com/GkadyrG/L0/backend/config"
	order "github.com/GkadyrG/L0/backend/internal/handler"
	"github.com/GkadyrG/L0/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func GetRouter(cfg *config.Config, h *order.Handler) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.CORS(cfg))
	router.Get("/api/order/{id}", h.GetByID())
	router.Get("/api/orders", h.GetAll())

	return router
}
