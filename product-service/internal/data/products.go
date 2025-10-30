package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/validator"
	"github.com/lib/pq"
)

type Product struct {
	ID          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       Price     `json:"price,string"`
	ImageUrl    string    `json:"image_url"`
	Stock       int32     `json:"stock"`
	Category    []string  `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int32     `json:"version"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(len(product.Name) <= 255, "name", "must not exceed 255 characters")

	v.Check(product.Description != "", "description", "must be provided")
	v.Check(len(product.Description) <= 5000, "description", "must not exceed 5000 characters")

	v.Check(product.Price > 0, "price", "must be greater than zero")
	v.Check(product.Price < 1000000, "price", "must be less than 1,000,000")

	v.Check(product.ImageUrl != "", "image_url", "must be provided")
	v.Check(len(product.ImageUrl) <= 1000, "image_url", "must not exceed 1000 characters")

	v.Check(product.Stock >= 0, "stock", "must be zero or greater")

	v.Check(len(product.Category) > 0, "category", "must have at least one category")
	v.Check(len(product.Category) <= 5, "category", "must not exceed 5 categories")
	v.Check(validator.Unique(product.Category), "category", "must not contain duplicate values")
}

type ProductModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (user_id, name, description, price, image_url, stock, category)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, created_at, updated_at, version`

	args := []interface{}{
		product.UserId,
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.Stock,
		pq.Array(product.Category),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// fmt.Printf("Query: %s\nArgs: %v\n", query, args) // Debug logging

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Version,
	)
}

// if product.Stock <= 0 { return errOutOfStock }
// Add a placeholder method for fetching a specific record from the movies table.
func (m ProductModel) Get(id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, user_id, name, description, price, image_url, stock, category, created_at, updated_at, version 
	FROM products
	WHERE id = $1`

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.UserId,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.ImageUrl,
		&product.Stock,
		pq.Array(&product.Category),
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &product, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m ProductModel) Update(product *Product) error {

	query := `
	UPDATE products
    SET name = $1, description = $2, price = $3, image_url = $4, 
        stock = $5, category = $6, updated_at = NOW(), version = version + 1
    WHERE id = $7 AND version = $8
    RETURNING version`

	args := []interface{}{
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.Stock,
		pq.Array(product.Category),
		product.ID,
		product.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&product.Version, &product.UpdatedAt)
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

// Add a placeholder method for deleting a specific record from the movies table.
func (m ProductModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM products
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use ExecContext() and pass the context as the first argument.
	result, err := m.DB.ExecContext(ctx, query, id)
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

func (m ProductModel) GetAll(name string, category []string, filters Filters) ([]*Product, Metadata, error) {
	query := fmt.Sprintf(`
	SELECT count(*) OVER(), id, name, description, price, image_url, stock, category, created_at, updated_at, version
	FROM products
	WHERE (to_tsvector('english', name) @@ plainto_tsquery('english', $1) OR $1 = '') 
	AND (array_length($2::text[], 1) IS NULL OR category && $2)
	ORDER BY %s %s, id ASC
	LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{name, pq.Array(category), filters.limit(), filters.offset()}
	fmt.Printf("Query: %s\nArgs: %v\n", query, args) // Debug logging

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	products := []*Product{}

	for rows.Next() {
		var product Product

		err := rows.Scan(
			&totalRecords,
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageUrl,
			&product.Stock,
			pq.Array(&product.Category),
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return products, metadata, nil
}
