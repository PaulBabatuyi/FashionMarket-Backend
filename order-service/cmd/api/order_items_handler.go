package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/data"
	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/validator"
)

func (app *application) createOrderItemHandler(w http.ResponseWriter, r *http.Request) {

	orderID, err := app.readIDParam(r, "order_id")
	if err != nil || orderID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		ProductID       int64   `json:"product_id"`
		ProductName     string  `json:"product_name"`
		ProductImageURL *string `json:"product_image_url,omitempty"`
		UnitPrice       float64 `json:"unit_price"`
		Quantity        int     `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
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

	item := &data.OrderItem{
		OrderID:         orderID,
		ProductID:       input.ProductID,
		ProductName:     input.ProductName,
		ProductImageURL: input.ProductImageURL,
		UnitPrice:       input.UnitPrice,
		Quantity:        input.Quantity,
	}

	// app.logger.PrintInfo("Order struct ", map[string]string{
	// 	"user_id": fmt.Sprintf("%d", order.UserID),
	// }

	// Initialize a new Validator instance.
	// Validate  index = 0 – only one item) ----
	v := validator.New()
	data.ValidateOrderItem(v, item, 0)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.OrderItems.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/orders/%d/items/%d", orderID, item.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"order_item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := app.readIDParam(r, "id")
	if err != nil || itemID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	orderID, err := app.readIDParam(r, "order_id")
	if err != nil || orderID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() || user.ID == 0 {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	item, err := app.models.OrderItems.Get(itemID, orderID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"order_item": item}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := app.readIDParam(r, "id")
	if err != nil || itemID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	orderID, err := app.readIDParam(r, "order_id")
	if err != nil || orderID < 1 {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	if user.IsAnonymous() || user.ID == 0 {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	item, err := app.models.OrderItems.Get(itemID, orderID, user.ID)
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
		Quantity *int `json:"quantity,omitempty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Quantity != nil {
		item.Quantity = *input.Quantity
	}

	v := validator.New()
	// index doesn't matter for single update
	if data.ValidateOrderItem(v, item, 0); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.OrderItems.Update(item)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order_item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	itemID, err := app.readIDParam(r, "id")
	if err != nil || itemID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	orderID, err := app.readIDParam(r, "order_id")
	if err != nil || orderID < 1 {
		app.notFoundResponse(w, r)
		return
	}
	user := app.contextGetUser(r)
	if user.IsAnonymous() || user.ID == 0 {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	err = app.models.OrderItems.Delete(itemID, orderID, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, envelope{"message": "item deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
