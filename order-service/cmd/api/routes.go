package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	// 	// Initialize a new Chi router instance
	router := chi.NewRouter()

	// Set the custom error handler for 404 Not Found responses using http.HandlerFunc
	router.NotFound(http.HandlerFunc(app.notFoundResponse))
	// Set the custom error handler for 405 Method Not Allowed responses using http.HandlerFunc
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	router.MethodFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Public routes
	// router.MethodFunc(http.MethodGet, "/v1/orders", app.listOrderHandler)
	router.MethodFunc(http.MethodGet, "/v1/orders/{id}", app.getOrderHandler)

	// Protected routes - require activated user
	// app.requireActivatedUser()
	router.MethodFunc(http.MethodPost, "/v1/orders", app.requireActivatedUser(app.createOrderHandler))
	router.MethodFunc(http.MethodPatch, "/v1/orders/{id}", app.requireActivatedUser(app.updateOrderHandler))
	router.MethodFunc(http.MethodDelete, "/v1/orders/{id}", app.requireActivatedUser(app.deleteOrderHandler))

	// Public routes
	// router.MethodFunc(http.MethodGet, "/v1/orders", app.listOrderHandler)
	// Protected routes - require activated user
	// app.requireActivatedUser()
	router.MethodFunc(http.MethodPost, "/v1/orders/{orderID}/items", app.createOrderItemHandler)
	router.MethodFunc(http.MethodGet, "/v1/orders/{orderID}/items/{id}", app.getOrderItemHandler)
	router.MethodFunc(http.MethodPatch, "/v1/orders/{orderID}/items/{id}", app.updateOrderItemHandler)
	router.MethodFunc(http.MethodDelete, "/v1/orders/{orderID}/items/{id}", app.deleteOrderItemHandler)

	//   router.Method(http.MethodGet, "/debug/vars", expvar.Handler())

	// Return the Chi router, which implements http.Handler
	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
	// return app.recoverPanic(app.rateLimit(app.authenticate(router)))

}
