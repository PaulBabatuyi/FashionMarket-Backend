package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PaulBabatuyi/FashionMarket-Backend/internal/data"
	"github.com/PaulBabatuyi/FashionMarket-Backend/pkg/validator"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string     `json:"name"`
		Description string     `json:"description"`
		Price       data.Price `json:"price"`
		ImageUrl    string     `json:"image_url"`
		Stock       int32      `json:"stock"`
		Category    []string   `json:"category"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Get the authenticated user
	user := app.contextGetUser(r)
	if user.IsAnonymous() || user.ID == 0 {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	// Log user_id for debugging
	// app.logger.PrintInfo("Creating product", map[string]string{
	// 	"user_id": fmt.Sprintf("%d", user.ID),
	// 	"input":   fmt.Sprintf("%+v", input),
	// })

	product := &data.Product{
		UserId:      user.ID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		ImageUrl:    input.ImageUrl,
		Stock:       input.Stock,
		Category:    input.Category,
	}

	// app.logger.PrintInfo("Product struct ", map[string]string{
	// 	"user_id": fmt.Sprintf("%d", product.UserId),
	// })

	// Initialize a new Validator instance.
	v := validator.New()

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)

	}

	err = app.models.Products.Insert(product)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/products/{%d}", product.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// Get the authenticated user from the context
	user := app.contextGetUser(r)
	if user == nil {
		app.authenticationRequiredResponse(w, r)
		return
	}

	product, err := app.models.Products.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	if user == nil {
		app.authenticationRequiredResponse(w, r)
		return
	}

	product, err := app.models.Products.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string     `json:"name"`
		Description *string     `json:"description"`
		Price       *data.Price `json:"price"`
		ImageUrl    *string     `json:"image_url"`
		Stock       *int32      `json:"stock"`
		Category    []string    `json:"category"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.ImageUrl != nil {
		product.ImageUrl = *input.ImageUrl
	}
	if input.Stock != nil {
		product.Stock = *input.Stock
	}
	if input.Category != nil {
		product.Category = input.Category
	}

	v := validator.New()

	if data.ValidateProduct(v, product); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Products.Update(product)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Products.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Category []string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Category = app.readCSV(qs, "category", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "price", "created_at", "-id", "-name", "-price", "-created_at", "category"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	products, metadata, err := app.models.Products.GetAll(input.Name, input.Category, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"products": products, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
