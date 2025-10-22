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

// fashion=# SELECT * FROM users where id = 13;
// -[ RECORD 1 ]-+---------------------------------------------------------------------------------------------------------------------------
// id            | 13
// email         | pau4@example.com
// password_hash | \x243261243132246a7a6d71524c4d4f42424271754233366430686a584f544d305a45764a46324a6235366f6434786642614e6b31387a6e756a316e71
// name          | paul four
// avatar_url    |
// activated     | f
// created_at    | 2025-10-20 14:25:46+01
// updated_at    | 2025-10-20 14:25:46+01
// version       | 1
// {"token": "ICTLTRWXESSNFX3KQSIBSYGSS4"}

// BODY='{"name":"Ice money","description":"richies cloths","price":43333.00,
// "image_url":"http://google.com","stock":50,"category":["men","boy"]}'
// curl -i -X POST -H "Content-Type: application/json" -H "Authorization: Bearer 2QTDZ7RI3SMRJUKCDAEUSHBWHU" -d "$BODY" localhost:4000/v1/products

//  curl -X POST -H "Content-Type: application/json" -d '{"email": "pau5@example.com", "password": "Pa55word"}' localhost:4000/v1/tokens/authentication{
//         "authentication token": {
//                 "token": "IPWOF6PLODNYOOCD3D23DXWP6U",
//                 "expiry": "2025-10-23T12:59:21.4757788+01:00",
//                 "created_at": "2025-10-22T12:59:21.4757788+01:00"
//         }
// }
