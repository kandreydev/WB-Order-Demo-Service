package usecase

import (
	"context"

	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/GkadyrG/L0/backend/internal/repository"
)

type UseCase struct {
	repo repository.OrderRepository
}

func New(repo repository.OrderRepository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) Save(ctx context.Context, order *model.Order) error {
	return u.repo.Save(ctx, order)
}

func (u *UseCase) GetByID(ctx context.Context, id string) (*model.OrderResponse, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *UseCase) GetAll(ctx context.Context) ([]*model.OrderPreview, error) {
	return u.repo.GetAll(ctx)
}

func (u *UseCase) GetAllFull(ctx context.Context, limit int) ([]*model.OrderResponse, error) {
	return u.repo.GetAllFull(ctx, limit)
}
