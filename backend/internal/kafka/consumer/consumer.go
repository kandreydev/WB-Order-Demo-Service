package consumer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/GkadyrG/L0/backend/internal/usecase"
	"github.com/GkadyrG/L0/backend/internal/validate"
	"github.com/IBM/sarama"
)

type consumerHandler struct {
	ready chan struct{}
	uc    *usecase.UseCase
}

type Consumer struct {
	group   sarama.ConsumerGroup
	handler *consumerHandler
	logger  *slog.Logger
}

func NewConsumer(brokers []string, groupID string, uc *usecase.UseCase, logger *slog.Logger) (*Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_8_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	g, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return nil, err
	}

	h := &consumerHandler{ready: make(chan struct{}), uc: uc}
	return &Consumer{group: g, handler: h, logger: logger}, nil
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("ConsumeClaim started", "topic", claim.Topic(), "partition", claim.Partition())

	for msg := range claim.Messages() {
		var order model.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			slog.Error("Unmarshal failed", "error", err)
			continue
		}

		if err := validate.ValidateOrder(order); err != nil {
			slog.Error("Invalid order", "error", err)
			continue
		}

		if err := h.uc.Save(sess.Context(), &order); err != nil {
			slog.Error("failed to save order", "error", err)
			continue
		}

		slog.Info("order saved")

		sess.MarkMessage(msg, "")
	}

	return nil
}

func (c *Consumer) Run(ctx context.Context, topics []string) {
	go func() {
		for {

			if ctx.Err() != nil {
				return
			}

			if err := c.group.Consume(ctx, topics, c.handler); err != nil {
				c.logger.Error("failed to consume messages", "error", err)
				return
			}

		}

	}()

	c.logger.Info("consumer started for topic", "topic", topics)
	<-c.handler.ready

}

func (c *Consumer) Close() error {
	c.logger.Info("closing consumer group")
	return c.group.Close()
}
