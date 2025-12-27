// Package domain provides core domain types for the TelemetryFlow SDK.
// This package implements Domain-Driven Design (DDD) patterns including
// value objects, entities, and domain services for telemetry configuration.
package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Credentials represents TelemetryFlow API credentials
// This is a Value Object in DDD - immutable and validates itself
type Credentials struct {
	keyID     string
	keySecret string
}

// NewCredentials creates and validates new credentials
func NewCredentials(keyID, keySecret string) (*Credentials, error) {
	if keyID == "" {
		return nil, errors.New("API key ID cannot be empty")
	}
	if keySecret == "" {
		return nil, errors.New("API key secret cannot be empty")
	}

	// Validate key ID format (should start with tfk_)
	if !strings.HasPrefix(keyID, "tfk_") {
		return nil, fmt.Errorf("invalid key ID format: must start with 'tfk_', got: %s", keyID)
	}

	// Validate key secret format (should start with tfs_)
	if !strings.HasPrefix(keySecret, "tfs_") {
		return nil, fmt.Errorf("invalid key secret format: must start with 'tfs_', got prefix: %s", keySecret[:4])
	}

	return &Credentials{
		keyID:     keyID,
		keySecret: keySecret,
	}, nil
}

// KeyID returns the API key ID
func (c *Credentials) KeyID() string {
	return c.keyID
}

// KeySecret returns the API key secret
func (c *Credentials) KeySecret() string {
	return c.keySecret
}

// AuthorizationHeader returns the formatted authorization header value
// Format: "Bearer <keyID>:<keySecret>"
func (c *Credentials) AuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s:%s", c.keyID, c.keySecret)
}

// Equals checks if two credentials are equal
func (c *Credentials) Equals(other *Credentials) bool {
	if other == nil {
		return false
	}
	return c.keyID == other.keyID && c.keySecret == other.keySecret
}

// String returns a safe string representation (hides the secret)
func (c *Credentials) String() string {
	return fmt.Sprintf("Credentials{keyID: %s, keySecret: ***}", c.keyID)
}
