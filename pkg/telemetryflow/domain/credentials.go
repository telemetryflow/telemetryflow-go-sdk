// Package domain provides core domain types for the TelemetryFlow SDK. This package implements Domain-Driven Design (DDD) patterns including value objects, entities, and domain services for telemetry configuration.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
