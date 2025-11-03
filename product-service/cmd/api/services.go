// cmd/api/services.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/product-service/internal/data"
)

// getUserFromUserService fetches user details from the user-service.
// This is called after JWT validation to get the full user profile.
func (app *application) getUserFromUserService(userID int64) (*data.User, error) {
	// Build the URL
	url := fmt.Sprintf("%s/v1/users/%d", app.config.userService.url, userID)

	// Create the request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set timeout context (optional, but recommended)
	ctx, cancel := app.createRequestContext(3 * time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	// Execute the request
	resp, err := app.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user-service request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	// Parse the response
	var envelope struct {
		User struct {
			ID        int64  `json:"id"`
			Email     string `json:"email"`
			Name      string `json:"name"`
			Activated bool   `json:"activated"`
		} `json:"user"`
	}

	err = json.NewDecoder(resp.Body).Decode(&envelope)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// Convert to internal User type
	user := &data.User{
		ID:        envelope.User.ID,
		Email:     envelope.User.Email,
		Name:      envelope.User.Name,
		Activated: envelope.User.Activated,
	}

	return user, nil
}

// Optional: Helper to create request contexts with timeout
func (app *application) createRequestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}

// getUserWithCache fetches user data with caching support.
// This is the PRIMARY method you should use in middleware.
// It checks cache first, then falls back to user-service.
func (app *application) getUserWithCache(userID int64) (*data.User, error) {
	// Try cache first
	if app.userCache != nil {
		if user, found := app.userCache.Get(userID); found {
			app.logger.PrintInfo("user cache hit", map[string]string{
				"user_id": fmt.Sprintf("%d", userID),
			})
			return user, nil
		}
	}

	// Cache miss - fetch from user-service
	app.logger.PrintInfo("user cache miss", map[string]string{
		"user_id": fmt.Sprintf("%d", userID),
	})

	user, err := app.getUserFromUserService(userID)
	if err != nil {
		return nil, err
	}

	// Store in cache for future requests
	if app.userCache != nil {
		app.userCache.Set(userID, user)
	}

	return user, nil
}

// invalidateUserCache removes a user from cache.
// Call this when you know user data has changed (e.g., after profile update).
func (app *application) invalidateUserCache(userID int64) {
	if app.userCache != nil {
		app.userCache.Delete(userID)
		app.logger.PrintInfo("user cache invalidated", map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		})
	}
}
