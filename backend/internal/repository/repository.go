package repository

import (
	"context"

	"github.com/GkadyrG/L0/backend/internal/apperr"
	"github.com/GkadyrG/L0/backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type Repo struct {
	conn *pgxpool.Pool
}

func New(conn *pgxpool.Pool) *Repo {
	return &Repo{conn: conn}
}

func (r *Repo) Save(ctx context.Context, order *model.Order) error {
	tx, err := r.conn.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	const ordersQuery = `
        INSERT INTO orders (
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created,
            oof_shard, created_at
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
    `

	_, err = tx.Exec(ctx, ordersQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
		order.CreatedAt,
	)
	if err != nil {
		return errors.Wrap(err, "insert order")
	}

	const deliveryQuery = `
	INSERT INTO delivery (
		order_uid, name, phone, zip, city, address, region, email
	) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
`
	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return errors.Wrap(err, "insert delivery")
	}

	const paymentQuery = `
        INSERT INTO payment (
            order_uid, transaction, request_id, currency, provider, amount,
            payment_dt, bank, delivery_cost, goods_total, custom_fee
        ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `
	_, err = tx.Exec(ctx, paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return errors.Wrap(err, "insert payment")
	}

	const itemsQuery = `
            INSERT INTO items (
                order_uid, chrt_id, track_number, price, rid, name, sale,
                size, total_price, nm_id, brand, status
            ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
        `
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemsQuery,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.RID,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status,
		)
		if err != nil {
			return errors.Wrap(err, "insert item")
		}
	}

	return nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*model.OrderResponse, error) {
	o := model.OrderResponse{}

	const orderdQuery = `
        SELECT 
            order_uid, track_number, customer_id, date_created
        FROM orders 
        WHERE order_uid = $1
    `
	err := r.conn.QueryRow(ctx, orderdQuery, id).Scan(
		&o.OrderUID,
		&o.TrackNumber,
		&o.CustomerID,
		&o.DateCreated,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(apperr.ErrNotFound, "not found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "get order")
	}

	const deliveryQuery = `
        SELECT name, phone, city, address, email
        FROM delivery WHERE order_uid = $1
    `
	err = r.conn.QueryRow(ctx, deliveryQuery, id).Scan(
		&o.Delivery.Name,
		&o.Delivery.Phone,
		&o.Delivery.City,
		&o.Delivery.Address,
		&o.Delivery.Email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(apperr.ErrNotFound, "not found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "get delivery")
	}

	const paymentQuery = `
        SELECT transaction, currency, amount
        FROM payment WHERE order_uid = $1
    `
	err = r.conn.QueryRow(ctx, paymentQuery, id).Scan(
		&o.Payment.Transaction,
		&o.Payment.Currency,
		&o.Payment.Amount,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errors.Wrap(apperr.ErrNotFound, "not found")
	}

	if err != nil {
		return nil, errors.Wrap(err, "get payment")
	}

	const itemsQuery = `
	SELECT price, name, brand, status
	FROM items WHERE order_uid = $1
	`
	rows, err := r.conn.Query(ctx, itemsQuery, id)
	if err != nil {
		return nil, errors.Wrap(err, "get items")
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ItemResponse
		if err := rows.Scan(
			&item.Price,
			&item.Name,
			&item.Brand,
			&item.Status,
		); err != nil {
			return nil, errors.Wrap(err, "scan item")
		}
		o.Items = append(o.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows iteration")
	}

	if len(o.Items) == 0 {
		return nil, errors.Wrap(apperr.ErrNotFound, "orders not found")
	}

	return &o, nil
}

func (r *Repo) GetAll(ctx context.Context) ([]*model.OrderPreview, error) {
	const orderQuery = `
	SELECT order_uid, track_number, customer_id, date_created
	FROM orders ORDER BY date_created DESC
	`
	rows, err := r.conn.Query(ctx, orderQuery)

	if err != nil {
		return nil, errors.Wrap(err, "get all orders")
	}
	defer rows.Close()

	previews := []*model.OrderPreview{}
	for rows.Next() {
		var p model.OrderPreview
		if err := rows.Scan(&p.OrderUID, &p.TrackNumber, &p.CustomerID, &p.DateCreated); err != nil {
			return nil, errors.Wrap(err, "scan preview")
		}
		previews = append(previews, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows iteration")
	}

	if len(previews) == 0 {
		return nil, errors.Wrap(apperr.ErrNotFound, "orders not found")
	}

	return previews, nil
}

func (r *Repo) GetAllFull(ctx context.Context, limit int) ([]*model.OrderResponse, error) {
	const mainQuery = `
        SELECT 
            o.order_uid, 
            o.track_number, 
            o.customer_id, 
            o.date_created,
            d.name, 
            d.phone, 
            d.city, 
            d.address, 
            d.email,
            p.transaction, 
            p.currency, 
            p.amount
        FROM orders o
        LEFT JOIN delivery d ON d.order_uid = o.order_uid
        LEFT JOIN payment p ON p.order_uid = o.order_uid
        ORDER BY o.date_created DESC
		LIMIT $1
    `

	mainRows, err := r.conn.Query(ctx, mainQuery, limit)
	if err != nil {
		return nil, errors.Wrap(err, "query main orders data")
	}
	defer mainRows.Close()

	ordersMap := make(map[string]*model.OrderResponse)
	orderUIDs := make([]string, 0)

	for mainRows.Next() {
		var o model.OrderResponse
		var d model.DeliveryResponse
		var p model.PaymentResponse

		err := mainRows.Scan(
			&o.OrderUID,
			&o.TrackNumber,
			&o.CustomerID,
			&o.DateCreated,
			&d.Name,
			&d.Phone,
			&d.City,
			&d.Address,
			&d.Email,
			&p.Transaction,
			&p.Currency,
			&p.Amount,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scan main row")
		}

		o.Delivery = d
		o.Payment = p
		o.Items = make([]model.ItemResponse, 0)

		ordersMap[o.OrderUID] = &o
		orderUIDs = append(orderUIDs, o.OrderUID)
	}

	if err := mainRows.Err(); err != nil {
		return nil, errors.Wrap(err, "main rows iteration")
	}

	if len(ordersMap) == 0 {
		return []*model.OrderResponse{}, nil
	}

	itemsMap, err := r.getItemsByOrderUIDs(ctx, orderUIDs)
	if err != nil {
		return nil, err
	}

	orders := make([]*model.OrderResponse, 0, len(ordersMap))
	for _, uid := range orderUIDs {
		if order, exists := ordersMap[uid]; exists {
			if items, ok := itemsMap[uid]; ok {
				order.Items = items
			}
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (r *Repo) getItemsByOrderUIDs(ctx context.Context, orderUIDs []string) (map[string][]model.ItemResponse, error) {
	const itemsQuery = `
        SELECT 
            order_uid, 
            price, 
            name, 
            brand, 
            status
        FROM items 
        WHERE order_uid = ANY($1)
    `

	rows, err := r.conn.Query(ctx, itemsQuery, orderUIDs)
	if err != nil {
		return nil, errors.Wrap(err, "query items")
	}
	defer rows.Close()

	itemsMap := make(map[string][]model.ItemResponse)

	for rows.Next() {
		var orderUID string
		var item model.ItemResponse

		err := rows.Scan(
			&orderUID,
			&item.Price,
			&item.Name,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scan item row")
		}

		itemsMap[orderUID] = append(itemsMap[orderUID], item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "items rows iteration")
	}

	return itemsMap, nil
}
