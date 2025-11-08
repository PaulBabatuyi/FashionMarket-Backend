package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Orders     OrderModel
	OrderItems OrderItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Orders:     OrderModel{DB: db},
		OrderItems: OrderItemModel{DB: db},
	}
}
