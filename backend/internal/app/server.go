package app

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/GkadyrG/L0/backend/internal/kafka/consumer"
)

type Server struct {
	httpsrv  *http.Server
	consumer *consumer.Consumer
	logger   *slog.Logger
}

func NewServer(httpsrv *http.Server, consumer *consumer.Consumer, logger *slog.Logger) *Server {
	return &Server{
		httpsrv:  httpsrv,
		consumer: consumer,
		logger:   logger,
	}
}

func (s *Server) Start(ctx context.Context, topics []string) error {
	go func() {
		s.consumer.Run(ctx, topics)
	}()

	s.logger.Info("starting http server", "addr", s.httpsrv.Addr)
	if err := s.httpsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("stopping server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpsrv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	if err := s.consumer.Close(); err != nil {
		return err
	}

	s.logger.Info("server stopped")

	return nil
}
