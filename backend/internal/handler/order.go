package order

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GkadyrG/L0/backend/internal/apperr"
	"github.com/GkadyrG/L0/backend/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Handler struct {
	us     usecase.OrderProvider
	logger *slog.Logger
}

func New(us usecase.OrderProvider, logger *slog.Logger) *Handler {
	return &Handler{us: us, logger: logger}
}

func (h *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := chi.URLParam(r, "id")
		order, err := h.us.GetByID(ctx, id)

		if err != nil {
			h.logger.Error("failed to get order", "err", err)

			if errors.Is(err, apperr.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "order not found"})
				return
			}

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "internal server error"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, order)

	}
}

func (h *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		ordersPreview, err := h.us.GetAll(ctx)
		if err != nil {
			h.logger.Error("failed to get all orders preview", "err", err)

			if errors.Is(err, apperr.ErrNotFound) {
				render.Status(r, http.StatusNotFound)
				render.JSON(w, r, map[string]string{"error": "orders preview not found"})
				return
			}

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "internal server error"})
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, ordersPreview)
	}
}
