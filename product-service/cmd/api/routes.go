package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	// Initialize a new Chi router instance
	router := chi.NewRouter()

	// Set the custom error handler for 404 Not Found responses using http.HandlerFunc
	router.NotFound(http.HandlerFunc(app.notFoundResponse))

	// Set the custom error handler for 405 Method Not Allowed responses using http.HandlerFunc
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	router.MethodFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// Register routes with method, URL patterns, and handler functions
	// app.requireActivatedUser
	router.MethodFunc(http.MethodGet, "/v1/products", app.listProductHandler)
	router.MethodFunc(http.MethodPost, "/v1/products", app.createProductHandler)
	router.MethodFunc(http.MethodGet, "/v1/products/{id}", app.showProductHandler)
	router.MethodFunc(http.MethodPatch, "/v1/products/{id}", app.updateProductHandler)
	router.MethodFunc(http.MethodDelete, "/v1/products/{id}", app.deleteProductHandler)

	// Return the Chi router, which implements http.Handler
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
