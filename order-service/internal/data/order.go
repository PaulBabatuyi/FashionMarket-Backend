package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/validator"
	_ "github.com/lib/pq"
)

const (
	StatusPending    = "pending"
	StatusPaid       = "paid"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusDelivered  = "delivered"
	StatusCancelled  = "cancelled"

	PaymentStatusUnpaid   = "unpaid"
	PaymentStatusPaid     = "paid"
	PaymentStatusRefunded = "refunded"
)

type ShippingAddress struct {
	Address string `json:"address"`
	Country string `json:"country"`
}

type Order struct {
	ID              int64           `json:"id"`
	UserID          int64           `json:"user_id"`
	TotalAmount     float64         `json:"total_amount"`
	Currency        string          `json:"currency"`
	Status          string          `json:"status"`
	PaymentStatus   string          `json:"payment_status"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
	Items           []OrderItem     `json:"items"`
	CreatedAt       time.Time       `json:"-"`
	UpdatedAt       time.Time       `json:"-"`
	Version         int32           `json:"version"`
}

type OrderModel struct {
	DB *sql.DB
}

func (o OrderModel) Insert(order *Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := o.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO orders
		 (user_id, total_amount, currency, status, payment_status, shipping_address)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at, version`

	args := []interface{}{
		order.UserID,
		order.TotalAmount,
		order.Currency,
		order.Status,
		order.PaymentStatus,
		order.ShippingAddress,
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(
		&order.ID,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Version,
	)
	if err != nil {
		return err
	}

	// 2. Insert items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		if err := o.insertItemTx(ctx, tx, &order.Items[i]); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (o OrderModel) insertItemTx(ctx context.Context, tx *sql.Tx, item *OrderItem) error {
	query := `
    INSERT INTO order_items (order_id, product_id, product_name, product_image_url,
	 unit_price, quantity)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, subtotal, created_at`

	return tx.QueryRowContext(ctx, query,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.ProductImageURL,
		item.UnitPrice,
		item.Quantity,
	).Scan(&item.ID, &item.Subtotal, &item.CreatedAt)
}

func (o OrderModel) Get(id, user_id int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, user_id, total_amount, currency, status, payment_status, 
	  shipping_address, created_at, updated_at, version 
	FROM orders
	WHERE id = $1 AND user_id = $2`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Currency,
		&order.Status,
		&order.PaymentStatus,
		&order.ShippingAddress,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	items, err := o.GetItems(order.ID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return &order, nil
}

func (o OrderModel) GetItems(orderID int64) ([]OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, product_name, product_image_url, unit_price, quantity, subtotal, created_at		FROM order_items
		WHERE order_id = $1
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductImageURL,
			&item.UnitPrice,
			&item.Quantity,
			&item.Subtotal,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (o OrderModel) Update(order *Order) error {

	query := `
	UPDATE orders
    SET total_amount = $1, currency = $2, status = $3, payment_status = $4, 
        shipping_address = $5, updated_at = NOW(), version = version + 1
    WHERE id = $6 AND version = $7
    RETURNING version, update_at`

	args := []interface{}{
		order.TotalAmount,
		order.Currency,
		order.Status,
		order.PaymentStatus,
		order.ShippingAddress,
		order.ID,
		order.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	err := o.DB.QueryRowContext(ctx, query, args...).Scan(&order.Version, &order.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (o OrderModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM orders
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use ExecContext() and pass the context as the first argument.
	result, err := o.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateOrder(v *validator.Validator, order *Order) {
	// Total Amount
	v.Check(order.TotalAmount > 0, "total_amount", "must be greater than zero")

	// Currency
	v.Check(order.Currency != "", "currency", "must be provided")
	v.Check(len(order.Currency) == 3, "currency", "must be a valid 3-letter ISO code")

	// Status
	v.Check(order.Status != "", "status", "must be provided")
	v.Check(validator.In(order.Status,
		StatusPending, StatusPaid, StatusProcessing,
		StatusShipped, StatusDelivered, StatusCancelled,
	), "status", "must be a valid status")

	// Payment Status
	v.Check(order.PaymentStatus != "", "payment_status", "must be provided")
	v.Check(validator.In(order.PaymentStatus,
		PaymentStatusUnpaid, PaymentStatusPaid, PaymentStatusRefunded,
	), "payment_status", "must be a valid payment status")

	// Shipping Address
	v.Check(order.ShippingAddress.Address != "", "shipping_address.address", "must be provided")
	v.Check(len(order.ShippingAddress.Address) <= 500, "shipping_address.address", "must not exceed 500 characters")

	v.Check(order.ShippingAddress.Country != "", "shipping_address.country", "must be provided")
	v.Check(len(order.ShippingAddress.Country) <= 100, "shipping_address.country", "must not exceed 100 characters")

	// Items
	v.Check(order.Items != nil, "items", "must be provided")
	v.Check(len(order.Items) >= 1, "items", "must contain at least 1 item")
	v.Check(len(order.Items) <= 100, "items", "cannot contain more than 100 items")

	// Validate each item
	for i, item := range order.Items {
		ValidateOrderItem(v, &item, i)
	}
}
