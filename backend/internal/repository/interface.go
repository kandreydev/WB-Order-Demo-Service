package repository

import (
	"context"

	"github.com/GkadyrG/L0/backend/internal/model"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=OrderRepository
type OrderRepository interface {
	Save(ctx context.Context, order *model.Order) error
	GetByID(ctx context.Context, id string) (*model.OrderResponse, error)
	GetAll(ctx context.Context) ([]*model.OrderPreview, error)
	GetAllFull(ctx context.Context, limit int) ([]*model.OrderResponse, error)
}
