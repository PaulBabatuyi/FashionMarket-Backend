package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// Signer abstracts the algorithm-specific details.
// The rest of the app only talks to this interface.
type Signer interface {
	// Algorithm returns the JWT "alg" header value (e.g. "ES256").
	Algorithm() string

	// Sign creates a compact JWT string from the claims.
	Sign(claims jwt.Claims) (string, error)

	// PublicKey returns the key used for verification.
	// This is used by jwt.ParseWithClaims().
	PublicKey() any
}
