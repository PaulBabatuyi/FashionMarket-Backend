package main

import (
	"context"
	"net/http"
	"time"
)

// Declare a handler which writes a plain-text response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dbStatus := "up"
	if err := app.models.Products.DB.PingContext(ctx); err != nil {
		dbStatus = "down"
	}

	// Check user-service
	userServiceStatus := "up"
	req, _ := http.NewRequestWithContext(ctx, "GET",
		app.config.userService.url+"/v1/healthcheck", nil)
	if _, err := app.httpClient.Do(req); err != nil {
		userServiceStatus = "down"
	}

	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment":  app.config.env,
			"version":      version,
			"database":     dbStatus,
			"user_service": userServiceStatus,
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
