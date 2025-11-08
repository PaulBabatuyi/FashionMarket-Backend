// internal/cache/user.go
package cache

import (
	"sync"
	"time"

	"github.com/PaulBabatuyi/FashionMarket-Backend/order-service/internal/data"
)

// UserCache provides thread-safe in-memory caching for user data
type UserCache struct {
	mu         sync.RWMutex
	users      map[int64]*cachedUser
	defaultTTL time.Duration
}

type cachedUser struct {
	user      *data.User
	expiresAt time.Time
}

// NewUserCache creates a new user cache with the specified default TTL
func NewUserCache(defaultTTL time.Duration) *UserCache {
	c := &UserCache{
		users:      make(map[int64]*cachedUser),
		defaultTTL: defaultTTL,
	}

	// Start background cleanup goroutine
	go c.cleanupExpired()

	return c
}

// Get retrieves a user from cache. Returns (user, true) if found and not expired,
// (nil, false) otherwise.
func (c *UserCache) Get(userID int64) (*data.User, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cached, exists := c.users[userID]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(cached.expiresAt) {
		return nil, false
	}

	return cached.user, true
}

// Set stores a user in cache with the default TTL
func (c *UserCache) Set(userID int64, user *data.User) {
	c.SetWithTTL(userID, user, c.defaultTTL)
}

// SetWithTTL stores a user in cache with a custom TTL
func (c *UserCache) SetWithTTL(userID int64, user *data.User, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.users[userID] = &cachedUser{
		user:      user,
		expiresAt: time.Now().Add(ttl),
	}
}

// Delete removes a user from cache
func (c *UserCache) Delete(userID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.users, userID)
}

// Clear removes all entries from cache
func (c *UserCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.users = make(map[int64]*cachedUser)
}

// Size returns the current number of cached users
func (c *UserCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.users)
}

// cleanupExpired runs periodically to remove expired entries
func (c *UserCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		now := time.Now()

		for userID, cached := range c.users {
			if now.After(cached.expiresAt) {
				delete(c.users, userID)
			}
		}

		c.mu.Unlock()
	}
}
