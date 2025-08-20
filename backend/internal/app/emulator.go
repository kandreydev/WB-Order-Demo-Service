package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/GkadyrG/L0/backend/config"
	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type EmulatorOptions struct {
	FilePath string
	Num      int
	Delay    time.Duration
}

func RunEmulator(ctx context.Context, cfg *config.Config, log *slog.Logger, opts EmulatorOptions) error {
	producer, err := newSyncProducer(cfg.GetKafkaBrokers())
	if err != nil {
		return errors.Wrap(err, "create producer")
	}
	defer producer.Close()

	baseOrder, err := loadBaseOrder(opts.FilePath)
	if err != nil {
		return err
	}

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < opts.Num; i++ {
		order := randomizeOrder(*baseOrder, rnd, i)

		if err := sendOrder(producer, cfg.Kafka.KafkaTopic, &order, log); err != nil {
			return errors.Wrap(err, "send order")
		}

		if opts.Delay > 0 && i < opts.Num-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(opts.Delay):
			}
		}
	}

	return nil
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Idempotent = true
	cfg.Net.MaxOpenRequests = 1
	cfg.Producer.Return.Successes = true
	return sarama.NewSyncProducer(brokers, cfg)
}

func sendOrder(p sarama.SyncProducer, topic string, order *model.Order, log *slog.Logger) error {
	b, err := json.Marshal(order)
	if err != nil {
		return errors.Wrap(err, "marshal order")
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(order.OrderUID),
		Value: sarama.ByteEncoder(b),
	}
	part, off, err := p.SendMessage(msg)
	if err != nil {
		return errors.Wrap(err, "send")
	}
	log.Info("sent", "partition", part, "offset", off, "order_uid", order.OrderUID)
	return nil
}

func loadBaseOrder(path string) (*model.Order, error) {
	if path == "" {
		return defaultOrder(), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var o model.Order
	if err := json.Unmarshal(data, &o); err != nil {
		return nil, err
	}
	return &o, nil
}

func defaultOrder() *model.Order {
	return &model.Order{
		OrderUID:          "b563feb7b2b84b6test",
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Date(2021, 11, 26, 6, 22, 19, 0, time.UTC),
		OofShard:          "1",
		CreatedAt:         time.Now(),
		Delivery: model.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				RID:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
	}
}

func randomizeOrder(base model.Order, rnd *rand.Rand, i int) model.Order {
	order := base
	stamp := time.Now().UnixNano()
	order.OrderUID = fmt.Sprintf("%s-%d-%d", base.OrderUID, i, stamp)
	order.TrackNumber = randAlphaNum(rnd, 12)
	order.CustomerID = fmt.Sprintf("cust-%d", rnd.Intn(100000))
	order.DateCreated = base.DateCreated.Add(time.Duration(i) * time.Second)
	order.CreatedAt = time.Now()
	order.Delivery.Name = fmt.Sprintf("Test User %d", rnd.Intn(1000))
	order.Delivery.Phone = "+" + fmt.Sprintf("%d", 100000000+rnd.Intn(900000000))
	order.Delivery.Email = fmt.Sprintf("user%d@example.com", rnd.Intn(100000))
	order.Payment.Transaction = fmt.Sprintf("%s-%d-%d", base.Payment.Transaction, i, stamp)
	order.Payment.Amount = int64(100 + rnd.Intn(5000))
	order.Payment.DeliveryCost = int64(50 + rnd.Intn(500))
	order.Payment.GoodsTotal = int64(50 + rnd.Intn(4000))
	if len(order.Items) > 0 {
		order.Items[0].ChrtID = order.Items[0].ChrtID + int64(i)
		order.Items[0].Price = int64(100 + rnd.Intn(1000))
		order.Items[0].TotalPrice = order.Items[0].Price - int64(order.Items[0].Sale)
		order.Items[0].TrackNumber = order.TrackNumber
		order.Items[0].RID = fmt.Sprintf("rid-%s", randAlphaNum(rnd, 8))
		order.Items[0].Name = fmt.Sprintf("Item-%d", rnd.Intn(1000))
		order.Items[0].Brand = []string{"Vivienne Sabo", "Acme", "Umbrella", "Globex"}[rnd.Intn(4)]
	}
	return order
}

func randAlphaNum(r *rand.Rand, n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
