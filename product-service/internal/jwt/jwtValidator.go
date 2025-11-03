// internal/jwt/validator.go
package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slices"
)

type JWTValidator struct {
	pubKey           *ecdsa.PublicKey
	expectedIssuer   string
	expectedAudience string
}

// Config holds JWT validation configuration
type Config struct {
	PublicKeyPath    string
	ExpectedIssuer   string
	ExpectedAudience string
}

func NewJWTValidator(cfg Config) (*JWTValidator, error) {
	// Read and parse public key
	pemData, err := os.ReadFile(cfg.PublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, errors.New("invalid PEM block")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return &JWTValidator{
		pubKey:           pubKey,
		expectedIssuer:   cfg.ExpectedIssuer,
		expectedAudience: cfg.ExpectedAudience,
	}, nil
}

// ValidatedClaims holds the validated JWT claims
type ValidatedClaims struct {
	UserID    int64
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// Validate validates the JWT token and returns the user ID
func (v *JWTValidator) Validate(tokenStr string) (*ValidatedClaims, error) {
	var claims jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := t.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Validate issuer
	if v.expectedIssuer != "" && claims.Issuer != v.expectedIssuer {
		return nil, fmt.Errorf("invalid issuer: expected %s, got %s", v.expectedIssuer, claims.Issuer)
	}

	// Validate audience
	if v.expectedAudience != "" && !slices.Contains(claims.Audience, v.expectedAudience) {
		return nil, fmt.Errorf("invalid audience: token not intended for this service")
	}

	// Validate expiration
	now := time.Now()
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(now) {
		return nil, errors.New("token expired")
	}

	// Validate not before
	if claims.NotBefore != nil && claims.NotBefore.After(now) {
		return nil, errors.New("token not yet valid")
	}

	// Parse user ID from subject
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid subject (user ID): %w", err)
	}

	result := &ValidatedClaims{
		UserID: userID,
	}

	if claims.IssuedAt != nil {
		result.IssuedAt = claims.IssuedAt.Time
	}
	if claims.ExpiresAt != nil {
		result.ExpiresAt = claims.ExpiresAt.Time
	}

	return result, nil
}
