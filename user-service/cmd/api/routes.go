package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	// Initialize a new Chi router instance
	router := chi.NewRouter()

	// Set the custom error handler for 404Not Found responses using http.HandlerFunc
	router.NotFound(http.HandlerFunc(app.notFoundResponse))

	// Set the custom error handler for 405 Method Not Allowed responses using http.HandlerFunc
	router.MethodNotAllowed(http.HandlerFunc(app.methodNotAllowedResponse))

	router.MethodFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//api endpoint needed by other services to interact with user-service
	router.MethodFunc(http.MethodGet, "/v1/users/{id}", app.requirePermission("products:read", app.getUserHandler))

	// Register routes with method, URL patterns, and handler functions
	// app.requireActivatedUser
	router.MethodFunc(http.MethodPost, "/v1/auth/register", app.registerUserHandler)
	router.MethodFunc(http.MethodPatch, "/v1/auth/activate", app.activateUserHandler)
	router.MethodFunc(http.MethodPost, "/v1/auth/token", app.createAuthenticationTokenHandler)

	router.MethodFunc(http.MethodPut, "/v1/users/password", app.updateUserPasswordHandler)
	router.MethodFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.MethodFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	// Register a new GET /debug/vars endpoint pointing to the expvar handler for metric
	router.Method(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}

// DELETE /v1/auth/logout              # Logout (delete tokens)
// POST   /v1/auth/login               # Login with email/password
// POST /v1/auth/password-reset          reset password
