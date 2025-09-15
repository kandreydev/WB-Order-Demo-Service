package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/GkadyrG/L0/backend/config"
	migrate "github.com/GkadyrG/L0/backend/database"
	"github.com/GkadyrG/L0/backend/internal/cache"
	order "github.com/GkadyrG/L0/backend/internal/handler"
	"github.com/GkadyrG/L0/backend/internal/kafka/consumer"
	"github.com/GkadyrG/L0/backend/internal/logger"
	"github.com/GkadyrG/L0/backend/internal/repository"
	"github.com/GkadyrG/L0/backend/internal/server"
	"github.com/GkadyrG/L0/backend/internal/storage"
	"github.com/GkadyrG/L0/backend/internal/usecase"
)

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig()

	logger := logger.SetupLogger(cfg.App.LogLevel)

	if err := migrate.RunMigrations(cfg.GetDSN(), cfg.MigratePath, logger); err != nil {
		logger.Error("migrate", slog.Any("err", err))
		return err
	}

	conn, err := storage.GetConnect(cfg.GetConnStr())
	if err != nil {
		logger.Error("connection pool", slog.Any("err", err))
		return err
	}

	repo := repository.New(conn)
	cacheDecorator, err := cache.New(ctx, cfg, repo)
	if err != nil {
		logger.Error("cache.New", slog.Any("err", err))
		return err
	}

	uc := usecase.New(cacheDecorator)
	handler := order.New(uc, logger)
	router := GetRouter(cfg, handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.App.Address, cfg.App.Port),
		Handler: router,
	}

	cons, err := consumer.NewConsumer(cfg.GetKafkaBrokers(), cfg.Kafka.KafkaGroupID, uc, logger)
	if err != nil {
		logger.Error("consumer.NewConsumer", slog.Any("err", err))
		return err
	}

	server := server.NewServer(srv, cons, logger)

	go func() {
		if err := RunEmulator(ctx, cfg, logger, EmulatorOptions{Num: cfg.EmulatorMessages}); err != nil {
			logger.Error("emulator failed", "err", err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.Start(ctx, cfg.GetKafkeTopics()); err != nil {
			logger.Error("failed to start server", slog.Any("err", err))
		}
	}()

	<-done

	if err := server.Stop(); err != nil {
		logger.Error("failed to stop server", slog.Any("err", err))
		return err
	}

	return nil
}
