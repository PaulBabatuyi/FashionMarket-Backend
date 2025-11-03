package jwt

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// ecdsaSigner holds the loaded private/public key pair.
type ecdsaSigner struct {
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
	alg  string
}

// NewECDSASigner loads PEM-encoded keys from disk.
// Both private and public keys are required.
// The public key can be extracted from the private key, but we keep them separate for clarity.
func NewECDSASigner(privateKeyPath, publicKeyPath, alg string) (Signer, error) {
	// --- Load private key ---
	privPEM, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key: %w", err)
	}

	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("invalid private key PEM: no PEM block found")
	}

	// Try PKCS#8 first (modern format)
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback to legacy SEC1 format
		priv, err = x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse private key: %w", err)
		}
	}

	ecdsaPriv, ok := priv.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not ECDSA")
	}

	// --- Load public key ---
	pubPEM, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read public key: %w", err)
	}

	block, _ = pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("invalid public key PEM: no PEM block found")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	ecdsaPub, ok := pubInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not ECDSA")
	}

	return &ecdsaSigner{
		priv: ecdsaPriv,
		pub:  ecdsaPub,
		alg:  alg,
	}, nil
}

// Algorithm returns the JWT algorithm identifier.
func (s *ecdsaSigner) Algorithm() string {
	return s.alg
}

// Sign creates a signed JWT using the private key.
func (s *ecdsaSigner) Sign(claims jwt.Claims) (string, error) {
	// Create a new token with the correct signing method
	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.alg), claims)
	// Sign and return compact string: header.payload.signature
	return token.SignedString(s.priv)
}

// PublicKey returns the public key for verification.
func (s *ecdsaSigner) PublicKey() any {
	return s.pub
}
