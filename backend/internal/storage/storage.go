package storage

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

func GetConnect(connStr string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "config parse")
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "newWithCongig")
	}
	if err = pool.Ping(context.Background()); err != nil {
		slog.Error("Failed to ping database", slog.Any("error", err))
		return nil, errors.Wrap(err, "connfig dbPing")
	}
	slog.Info("Successfully connected to database")
	return pool, nil
}
