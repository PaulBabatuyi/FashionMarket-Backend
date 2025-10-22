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

	router.MethodFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.MethodFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.MethodFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Return the Chi router, which implements http.Handler
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}

// POST   /v1/auth/register            # Register with email/password
// POST   /v1/auth/activate            # Activate account (if using token-based activation)
// GET    /v1/auth/me                  # Get current user

// DELETE /v1/auth/logout              # Logout (delete tokens)
// POST   /v1/auth/login               # Login with email/password

// GET    /v1/products                 # List all products (with filters/search via Elasticsearch)
// POST   /v1/products                 # Create product (as current user)
// GET    /v1/products/:id             # Get single product
// PATCH  /v1/products/:id             # Update product (own only)
// DELETE /v1/products/:id             # Delete product (own only)

// GET    /v1/orders                   # List users orders
// POST   /v1/orders                   # Create order
// GET    /v1/orders/:id               # Get order details
// PATCH  /v1/orders/:id               # Update order status (e.g., cancel)

// POST   /v1/payments/initialize      # Initialize Stripe payment
// POST   /v1/payments/verify          # Verify payment (webhook)
// GET    /v1/payments/:id             # Get payment details

// GET    /v1/healthcheck              # API health status

