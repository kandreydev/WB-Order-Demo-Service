package model

import "time"

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required"`
	Locale            string    `json:"locale" validate:"required"`
	InternalSignature string    `json:"internal_signature"`
	CustomerID        string    `json:"customer_id" validate:"required"`
	DeliveryService   string    `json:"delivery_service" validate:"required"`
	ShardKey          string    `json:"shardkey" validate:"required"`
	SmID              int       `json:"sm_id" validate:"required,min=1"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required"`
	CreatedAt         time.Time `json:"created_at" validate:"required"`

	Delivery Delivery `json:"delivery" validate:"required"`
	Payment  Payment  `json:"payment" validate:"required"`
	Items    []Item   `json:"items" validate:"required,dive"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string `json:"transaction" validate:"required"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency" validate:"required,len=3"`
	Provider     string `json:"provider" validate:"required"`
	Amount       int64  `json:"amount" validate:"required,gt=0"`
	PaymentDT    int64  `json:"payment_dt" validate:"required"`
	Bank         string `json:"bank"`
	DeliveryCost int64  `json:"delivery_cost" validate:"gte=0"`
	GoodsTotal   int64  `json:"goods_total" validate:"gte=0"`
	CustomFee    int64  `json:"custom_fee" validate:"gte=0"`
}

type Item struct {
	ChrtID      int64  `json:"chrt_id" validate:"required"`
	TrackNumber string `json:"track_number" validate:"required"`
	Price       int64  `json:"price" validate:"required,gt=0"`
	RID         string `json:"rid" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Sale        int    `json:"sale" validate:"gte=0"`
	Size        string `json:"size"`
	TotalPrice  int64  `json:"total_price" validate:"required,gt=0"`
	NmID        int64  `json:"nm_id"`
	Brand       string `json:"brand" validate:"required"`
	Status      int    `json:"status" validate:"required,gt=0"`
}

type OrderPreview struct {
	OrderUID    string    `json:"order_uid"`
	TrackNumber string    `json:"track_number"`
	CustomerID  string    `json:"customer_id"`
	DateCreated time.Time `json:"date_created"`
}

type OrderResponse struct {
	OrderUID    string           `json:"order_uid"`
	TrackNumber string           `json:"track_number"`
	CustomerID  string           `json:"customer_id"`
	DateCreated time.Time        `json:"date_created"`
	Delivery    DeliveryResponse `json:"delivery"`
	Payment     PaymentResponse  `json:"payment"`
	Items       []ItemResponse   `json:"items"`
}

type DeliveryResponse struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Address string `json:"address"`
	Email   string `json:"email"`
}

type PaymentResponse struct {
	Transaction string `json:"transaction"`
	Currency    string `json:"currency"`
	Amount      int64  `json:"amount"`
}

type ItemResponse struct {
	Name   string `json:"name"`
	Price  int64  `json:"price"`
	Brand  string `json:"brand"`
	Status int    `json:"status"`
}

func (o Order) ToResponse() *OrderResponse {
	delivery := DeliveryResponse{
		Name:    o.Delivery.Name,
		Phone:   o.Delivery.Phone,
		City:    o.Delivery.City,
		Address: o.Delivery.Address,
		Email:   o.Delivery.Email,
	}

	payment := PaymentResponse{
		Transaction: o.Payment.Transaction,
		Currency:    o.Payment.Currency,
		Amount:      o.Payment.Amount,
	}

	items := make([]ItemResponse, len(o.Items))
	for i, item := range o.Items {
		items[i] = ItemResponse{
			Name:  item.Name,
			Price: item.Price,
			Brand: item.Brand,
		}
	}

	return &OrderResponse{
		OrderUID:    o.OrderUID,
		TrackNumber: o.TrackNumber,
		CustomerID:  o.CustomerID,
		DateCreated: o.DateCreated,
		Delivery:    delivery,
		Payment:     payment,
		Items:       items,
	}
}
