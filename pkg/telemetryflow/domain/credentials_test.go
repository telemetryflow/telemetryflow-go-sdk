// Package domain provides unit tests for the domain layer.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCredentials(t *testing.T) {
	t.Run("should create credentials with valid inputs", func(t *testing.T) {
		creds, err := NewCredentials("tfk_test", "tfs_secret")

		require.NoError(t, err)
		assert.NotNil(t, creds)
		assert.Equal(t, "tfk_test", creds.KeyID())
		assert.Equal(t, "tfs_secret", creds.KeySecret())
	})

	t.Run("should fail with empty key ID", func(t *testing.T) {
		_, err := NewCredentials("", "tfs_secret")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key ID")
	})

	t.Run("should fail with empty key secret", func(t *testing.T) {
		_, err := NewCredentials("tfk_test", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "key secret")
	})

	t.Run("should fail with invalid key ID prefix", func(t *testing.T) {
		_, err := NewCredentials("invalid_key", "tfs_secret")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tfk_")
	})

	t.Run("should fail with invalid key secret prefix", func(t *testing.T) {
		_, err := NewCredentials("tfk_test", "invalid_secret")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tfs_")
	})
}

func TestCredentials_AuthorizationHeader(t *testing.T) {
	t.Run("should generate correct authorization header", func(t *testing.T) {
		creds, err := NewCredentials("tfk_test", "tfs_secret")
		require.NoError(t, err)

		header := creds.AuthorizationHeader()

		assert.Contains(t, header, "Bearer")
		assert.NotEmpty(t, header)
	})
}

func TestCredentials_Equals(t *testing.T) {
	t.Run("should return true for equal credentials", func(t *testing.T) {
		creds1, _ := NewCredentials("tfk_test", "tfs_secret")
		creds2, _ := NewCredentials("tfk_test", "tfs_secret")

		assert.True(t, creds1.Equals(creds2))
	})

	t.Run("should return false for different key IDs", func(t *testing.T) {
		creds1, _ := NewCredentials("tfk_test1", "tfs_secret")
		creds2, _ := NewCredentials("tfk_test2", "tfs_secret")

		assert.False(t, creds1.Equals(creds2))
	})

	t.Run("should return false for different key secrets", func(t *testing.T) {
		creds1, _ := NewCredentials("tfk_test", "tfs_secret1")
		creds2, _ := NewCredentials("tfk_test", "tfs_secret2")

		assert.False(t, creds1.Equals(creds2))
	})

	t.Run("should return false when comparing with nil", func(t *testing.T) {
		creds, _ := NewCredentials("tfk_test", "tfs_secret")

		assert.False(t, creds.Equals(nil))
	})
}

func TestCredentials_String(t *testing.T) {
	t.Run("should hide secret in string representation", func(t *testing.T) {
		creds, _ := NewCredentials("tfk_test", "tfs_secret")

		str := creds.String()

		assert.Contains(t, str, "tfk_test")
		assert.NotContains(t, str, "tfs_secret")
		assert.Contains(t, str, "***")
	})
}
