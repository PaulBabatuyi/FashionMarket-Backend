package data

import (
	"time"
)

const (
	StatusPending    = "pending"
	StatusProcessing = "processing"
	StatusShipped    = "shipped"
	StatusDelivered  = "delivered"
	StatusCanceled   = "canceled"
)

type OrderModel struct {
	ID           int64     `json:"id"`
	UserId       int64     `json:"user_id"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	ShippingAddr []string  `json:"shipping_addr"`
	Category     []string  `json:"category,omitempty"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	Version      int32     `json:"version"`
}
