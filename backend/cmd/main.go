package main

import (
	"log/slog"
	"os"

	"github.com/GkadyrG/L0/backend/internal/app"
)

// var (
//	dsn         = "postgresql://orders_user:orders_pass@localhost:5432/orders_db?sslmode=disable"
//	migratePath = "database/migrations"
// )

func main() {
	if err := app.Run(); err != nil {
		slog.Error("app run", slog.Any("err", err))
		os.Exit(1)
	}

	// logg := logger.SetupLogger(logger.EnvLocal)
	// cfg := config.LoadConfig()

	// conn, err := storage.GetConnect(cfg.GetConnStr())
	// if err != nil {
	// 	slog.Error("connection pool", slog.Any("err", err))
	// 	return
	// }

	// repo := repository.New(conn)

	// if err := migrate.RunMigrations(dsn, migratePath, logg); err != nil {
	// 	slog.Error("migrate", slog.Any("err", err))
	// 	return
	// }
	// fmt.Println("✅ Migrations ran successfully")

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// cacheDecorator, err := cache.New(ctx, repo)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	// order := &model.Order{
	// 	OrderUID:          fmt.Sprintf("test-uid-%d", time.Now().UnixNano()),
	// 	TrackNumber:       "track-123",
	// 	Entry:             "entry-1",
	// 	Locale:            "en",
	// 	InternalSignature: "sig-abc",
	// 	CustomerID:        "cust-123",
	// 	DeliveryService:   "DHL",
	// 	ShardKey:          "shard-1",
	// 	SmID:              42,
	// 	DateCreated:       time.Now(),
	// 	OofShard:          "oof-1",
	// 	CreatedAt:         time.Now(),
	// 	Delivery: model.Delivery{
	// 		Name:    "John Doe",
	// 		Phone:   "+123456789",
	// 		Zip:     "12345",
	// 		City:    "Test City",
	// 		Address: "Test Street 1",
	// 		Region:  "Test Region",
	// 		Email:   "john@example.com",
	// 	},
	// 	Payment: model.Payment{
	// 		Transaction:  "txn-123",
	// 		RequestID:    "req-123",
	// 		Currency:     "USD",
	// 		Provider:     "Visa",
	// 		Amount:       1000,
	// 		PaymentDT:    time.Now().Unix(),
	// 		Bank:         "Test Bank",
	// 		DeliveryCost: 50,
	// 		GoodsTotal:   950,
	// 		CustomFee:    0,
	// 	},
	// 	Items: []model.Item{
	// 		{
	// 			ChrtID:      1,
	// 			TrackNumber: "track-123",
	// 			Price:       500,
	// 			RID:         "rid-1",
	// 			Name:        "Item 1",
	// 			Sale:        10,
	// 			Size:        "M",
	// 			TotalPrice:  450,
	// 			NmID:        1001,
	// 			Brand:       "Brand A",
	// 			Status:      1,
	// 		},
	// 		{
	// 			ChrtID:      2,
	// 			TrackNumber: "track-123",
	// 			Price:       500,
	// 			RID:         "rid-2",
	// 			Name:        "Item 2",
	// 			Sale:        0,
	// 			Size:        "L",
	// 			TotalPrice:  500,
	// 			NmID:        1002,
	// 			Brand:       "Brand B",
	// 			Status:      1,
	// 		},
	// 	},
	// }

	// // Сохраняем заказ
	// if err := cacheDecorator.Save(ctx, order); err != nil {
	// 	log.Fatalf("failed to save order: %v", err)
	// }

	// // Получаем заказ по ID
	// gotOrder, err := cacheDecorator.GetByID(ctx, "test-uid-123")
	// if err != nil {
	// 	log.Fatalf("failed to get order by ID: %v", err)
	// }
	// fmt.Printf("Order fetched by ID:\n%+v\n\n", gotOrder)

	// // Получаем все заказы
	// allOrders, err := cacheDecorator.GetAll(ctx)
	// if err != nil {
	// 	log.Fatalf("failed to get all orders: %v", err)
	// }
	// fmt.Printf("All orders previews:\n")
	// for _, o := range allOrders {
	// 	fmt.Printf("- %+v\n", o)
	// }

	// fmt.Println("\n✅ Все операции с репозиторием прошли успешно!")
}
