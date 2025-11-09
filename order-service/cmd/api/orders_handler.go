package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/data"
	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/validator"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		TotalAmount     float64              `json:"total_amount"`
		Currency        string               `json:"currency"`
		Status          string               `json:"status"`
		PaymentStatus   string               `json:"payment_status"`
		ShippingAddress data.ShippingAddress `json:"shipping_address"`
		Items           []data.OrderItem     `json:"items"`
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

	order := &data.Order{
		UserID:          user.ID,
		TotalAmount:     input.TotalAmount,
		Currency:        input.Currency,
		Status:          input.Status,
		PaymentStatus:   input.PaymentStatus,
		ShippingAddress: input.ShippingAddress,
		Items:           input.Items,
	}
	// app.logger.PrintInfo("Order struct ", map[string]string{
	// 	"user_id": fmt.Sprintf("%d", order.UserID),
	// })

	// Initialize a new Validator instance.
	v := validator.New()
	if data.ValidateOrder(v, order); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Orders.Insert(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/orders/%d", order.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"order": order}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
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

	order, err := app.models.Orders.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
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

	order, err := app.models.Orders.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if order.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	var input struct {
		Status          *string               `json:"status,omitempty"`
		PaymentStatus   *string               `json:"payment_status,omitempty"`
		ShippingAddress *data.ShippingAddress `json:"shipping_address,omitempty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Status != nil {
		order.Status = *input.Status
	}
	if input.PaymentStatus != nil {
		order.PaymentStatus = *input.PaymentStatus
	}
	if input.ShippingAddress != nil {
		order.ShippingAddress = *input.ShippingAddress
	}

	v := validator.New()
	if data.ValidateOrder(v, order); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Orders.Update(order)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "id")
	if err != nil {
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

	order, err := app.models.Orders.Get(id, user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if order.UserID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Orders.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "order deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{
		"id", "-id",
		"total_amount", "-total_amount",
		"created_at", "-created_at",
		"status", "-status",
	}

	data.ValidateFilters(v, input.Filters)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.authenticationRequiredResponse(w, r)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}

	orders, metadata, err := app.models.Orders.GetAll(user.ID, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"orders":   orders,
		"metadata": metadata,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
