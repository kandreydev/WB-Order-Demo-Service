package cache

import (
	"context"
	"sync"
	"time"

	"github.com/GkadyrG/L0/backend/config"
	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/GkadyrG/L0/backend/internal/repository"
	"github.com/pkg/errors"
)

const cacheInitLimit = 1000

type wrapOrder struct {
	order     *model.OrderResponse
	updatedAt time.Time
}

type CacheDecorator struct {
	repo repository.OrderRepository

	mu     sync.RWMutex
	orders map[string]wrapOrder
}

func New(ctx context.Context, cfg *config.Config, orderRepo repository.OrderRepository) (*CacheDecorator, error) {
	cache := &CacheDecorator{
		orders: make(map[string]wrapOrder),
		repo:   orderRepo,
	}

	if err := cache.initializeCache(ctx); err != nil {
		return nil, err
	}

	cache.сleanCash(cfg.Cache.CleanupInterval, cfg.Cache.TTL)

	return cache, nil
}

func (c *CacheDecorator) initializeCache(ctx context.Context) error {
	orders, err := c.repo.GetAllFull(ctx, cacheInitLimit)
	if err != nil {
		return errors.Wrap(err, "failed to load all orders for cache initialization")
	}
	for _, order := range orders {
		c.set(order)
	}

	return nil
}

func (c *CacheDecorator) set(order *model.OrderResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = wrapOrder{
		order:     order,
		updatedAt: time.Now(),
	}
}

func (c *CacheDecorator) get(id string) (*model.OrderResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	wrap, ok := c.orders[id]
	return wrap.order, ok
}

func (c *CacheDecorator) сleanCash(cleanupInterval time.Duration, timeToLive time.Duration) {
	go func() {
		ticker := time.NewTicker(cleanupInterval)
		defer ticker.Stop()

		for range ticker.C {
			c.mu.Lock()
			for key, value := range c.orders {
				if time.Now().After(value.updatedAt.Add(timeToLive)) {
					delete(c.orders, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

func (c *CacheDecorator) Save(ctx context.Context, order *model.Order) error {
	if err := c.repo.Save(ctx, order); err != nil {
		return err
	}
	c.set(order.ToResponse())
	return nil
}

func (c *CacheDecorator) GetByID(ctx context.Context, id string) (*model.OrderResponse, error) {
	order, exists := c.get(id)
	if exists {
		return order, nil
	}

	order, err := c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	c.set(order)
	return order, nil
}

func (c *CacheDecorator) GetAll(ctx context.Context) ([]*model.OrderPreview, error) {
	return c.repo.GetAll(ctx)
}

func (c *CacheDecorator) GetAllFull(ctx context.Context, limit int) ([]*model.OrderResponse, error) {
	return c.repo.GetAllFull(ctx, limit)
}
