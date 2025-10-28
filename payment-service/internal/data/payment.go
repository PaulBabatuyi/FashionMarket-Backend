package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/internal/validator"
	"github.com/lib/pq"
)

type Product struct {
	ID          int64     `json:"id"`
	UserId      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       Price     `json:"price,omitempty,string"`
	ImageUrl    string    `json:"image_url"`
	Stock       int32     `json:"stock"`
	Category    []string  `json:"category,omitempty"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Version     int32     `json:"version"`
}

func ValidateProduct(v *validator.Validator, product *Product) {
	v.Check(product.Name != "", "name", "must be provided")
	v.Check(len(product.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(product.Price != 0, "price", "must be provided")
	v.Check(product.Price >= 50, "price", "must be greater than %50")
	v.Check(product.Stock != 0, "stock", "must be provided")
	v.Check(product.Description != "", "description", "must be provided")
	v.Check(product.ImageUrl != "", "image", "must be provided")
	v.Check(product.Category != nil, "Category", "must be provided")
	v.Check(len(product.Category) >= 1, "Category", "must contain at least 1 Category")
	v.Check(len(product.Category) <= 4, "Category", "available Category cloth are Men, Women, boy and girl ")
	// Note that we're using the Unique helper in the line below to check that all
	// values in the input.Category slice are unique.
	v.Check(validator.Unique(product.Category), "Category", "must not contain duplicate values")
	// Use the Valid() method to see if any of the checks failed. If they did, then use
}

// func ValidateProductName(v *validator.Validator, name string) {
//     v.Check(name != "", "name", "must be provided")
//     v.Check(len(name) <= 255, "name", "must not be more than 255 characters long")
// }

// func ValidateDescription(v *validator.Validator, description string) {
//     v.Check(description != "", "description", "must be provided")
// }

// func ValidatePrice(v *validator.Validator, price float64) {
//     v.Check(price > 0, "price", "must be greater than 0")
//     v.Check(price <= 999999.99, "price", "must not exceed 999999.99")
// }

// func ValidateImageURL(v *validator.Validator, imageURL string) {
//     v.Check(imageURL != "", "image_url", "must be provided")
//     v.Check(validator.Matches(imageURL, validator.URLRX), "image_url", "must be a valid URL")
// }

// func ValidateStock(v *validator.Validator, stock int32) {
//     v.Check(stock >= 0, "stock", "must be non-negative")
// }

// func ValidateCategory(v *validator.Validator, category []string) {
//     v.Check(len(category) > 0, "category", "must include at least one category")
//     for _, cat := range category {
//         v.Check(cat != "", "category", "must not contain empty values")
//     }
// }

type ProductModel struct {
	DB *sql.DB
}

// Add a placeholder method for inserting a new record in the movies table.
func (m ProductModel) Insert(product *Product) error {
	query := `
        INSERT INTO products (user_id, name, description, price, image_url, stock, category)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id, user_id, created_at, updated_at, version`

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
		&product.UserId,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Version,
	)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m ProductModel) Get(id, user_id int64) (*Product, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, user_id, name, description, price, image_url, stock, category, created_at, updated_at, version 
	FROM products
	WHERE id = $1 AND user_id = $2`

	var product Product

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// Importantly, use defer to make sure that we cancel the context before the Get()
	// method returns.
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id, user_id).Scan(
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
	SET name = $1, description = $2, price = $3, image_url = $4, stock = $5, category = $6, updated_at = $7, version = version + 1
	WHERE id = $8 AND user_id = $9 AND version = $10
	RETURNING version`

	args := []interface{}{
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.Stock,
		pq.Array(product.Category),
		product.UpdatedAt,
		product.ID,
		product.UserId,
		product.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Use QueryRowContext() and pass the context as the first argument.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&product.Version)
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
	SELECT count(*) OVER(), id, user_id, name, description, price, image_url, stock, category, created_at, updated_at, version
	FROM products
	WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '') 
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
