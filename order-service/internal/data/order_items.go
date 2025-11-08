package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/validator"
	_ "github.com/lib/pq"
)

type OrderItem struct {
	ID              int64     `json:"id"`
	OrderID         int64     `json:"-"`
	ProductID       int64     `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductImageURL *string   `json:"product_image_url,omitempty"`
	UnitPrice       float64   `json:"unit_price"`
	Quantity        int       `json:"quantity"`
	Subtotal        float64   `json:"subtotal"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type OrderItemModel struct {
	DB *sql.DB
}

func (o OrderItemModel) Insert(item *OrderItem) error {
	query := `
        INSERT INTO order_items (
		 order_id, product_id, product_name, product_image_url, unit_price, quantity
		 ) VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, subtotal, created_at`

	args := []interface{}{
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.ProductImageURL,
		item.UnitPrice,
		item.Quantity,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return o.DB.QueryRowContext(ctx, query, args...).Scan(
		&item.ID,
		&item.Subtotal,
		&item.CreatedAt,
	)
}

func (m OrderItemModel) Get(id, orderID, userID int64) (*OrderItem, error) {
	if id < 1 || orderID < 1 || userID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT oi.id, oi.order_id, oi.product_id, oi.product_name, oi.product_image_url,
               oi.unit_price, oi.quantity, oi.subtotal, oi.created_at, oi.updated_at
               o.user_id
        FROM order_items oi
        JOIN orders o ON oi.order_id = o.id
        WHERE oi.id = $1 AND oi.order_id = $2 AND o.user_id = $3`

	var item OrderItem
	var dbUserID int64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, orderID, userID).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.ProductName,
		&item.ProductImageURL,
		&item.UnitPrice,
		&item.Quantity,
		&item.Subtotal,
		&item.CreatedAt,
		&item.UpdatedAt,
		&dbUserID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

func (o OrderItemModel) Update(item *OrderItem) error {

	query := `
	UPDATE order_items
    SET quantity = $1, update_at = NOW() 
    WHERE id = $2 AND subtotal = $3
	RETURNING subtotal`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newSubtotal float64
	err := o.DB.QueryRowContext(ctx, query, item.Quantity, item.ID, item.Subtotal).Scan(item.Subtotal)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	item.Subtotal = newSubtotal
	return nil
}

func (o OrderItemModel) Delete(id, orderID, userID int64) error {
	if id < 1 || orderID < 1 || userID < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM order_items oi
        USING orders o
        WHERE oi.id = $1
          AND oi.order_id = $2
          AND oi.order_id = o.id
          AND o.user_id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := o.DB.ExecContext(ctx, query, id, orderID, userID)
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

func ValidateOrderItem(v *validator.Validator, item *OrderItem, index int) {
	prefix := func(field string) string {
		return fmt.Sprintf("items[%d].%s", index, field)
	}

	v.Check(item.ProductID > 0, prefix("product_id"), "must be a positive integer")
	v.Check(item.ProductName != "", prefix("product_name"), "must be provided")
	v.Check(len(item.ProductName) <= 255, prefix("product_name"), "must not exceed 255 characters")

	v.Check(item.UnitPrice > 0, prefix("unit_price"), "must be greater than zero")
	v.Check(item.Quantity > 0, prefix("quantity"), "must be greater than zero")
	v.Check(item.Quantity <= 1000, prefix("quantity"), "cannot exceed 1000 units")

	// Optional image URL
	if item.ProductImageURL != nil {
		v.Check(*item.ProductImageURL != "", prefix("product_image_url"), "must not be empty if provided")
		v.Check(validator.IsURL(*item.ProductImageURL), prefix("product_image_url"), "must be a valid URL")
	}
}
