// Package domain_test provides unit tests for the domain package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func TestNewCredentials(t *testing.T) {
	t.Run("should create credentials with valid inputs", func(t *testing.T) {
		creds, err := domain.NewCredentials("tfk_test123", "tfs_secret456")
		require.NoError(t, err)
		assert.NotNil(t, creds)
		assert.Equal(t, "tfk_test123", creds.KeyID())
		assert.Equal(t, "tfs_secret456", creds.KeySecret())
	})

	t.Run("should reject empty key ID", func(t *testing.T) {
		creds, err := domain.NewCredentials("", "tfs_secret456")
		require.Error(t, err)
		assert.Nil(t, creds)
		assert.Contains(t, err.Error(), "key ID")
	})

	t.Run("should reject empty key secret", func(t *testing.T) {
		creds, err := domain.NewCredentials("tfk_test123", "")
		require.Error(t, err)
		assert.Nil(t, creds)
		assert.Contains(t, err.Error(), "key secret")
	})

	t.Run("should reject key ID without tfk_ prefix", func(t *testing.T) {
		creds, err := domain.NewCredentials("invalid_key", "tfs_secret456")
		require.Error(t, err)
		assert.Nil(t, creds)
		assert.Contains(t, err.Error(), "tfk_")
	})

	t.Run("should reject key secret without tfs_ prefix", func(t *testing.T) {
		creds, err := domain.NewCredentials("tfk_test123", "invalid_secret")
		require.Error(t, err)
		assert.Nil(t, creds)
		assert.Contains(t, err.Error(), "tfs_")
	})

	t.Run("should accept various valid key formats", func(t *testing.T) {
		testCases := []struct {
			keyID     string
			keySecret string
		}{
			{"tfk_short", "tfs_short"},
			{"tfk_with_underscore_123", "tfs_with_underscore_456"},
			{"tfk_UPPERCASE", "tfs_UPPERCASE"},
			{"tfk_MixedCase123", "tfs_MixedCase456"},
			{"tfk_" + string(make([]byte, 100)), "tfs_" + string(make([]byte, 100))}, // long keys
		}

		for _, tc := range testCases {
			creds, err := domain.NewCredentials(tc.keyID, tc.keySecret)
			assert.NoError(t, err, "should accept keyID=%s", tc.keyID)
			assert.NotNil(t, creds)
		}
	})
}

func TestCredentials_AuthorizationHeader(t *testing.T) {
	t.Run("should format authorization header correctly", func(t *testing.T) {
		creds, err := domain.NewCredentials("tfk_mykey", "tfs_mysecret")
		require.NoError(t, err)

		header := creds.AuthorizationHeader()
		assert.Equal(t, "Bearer tfk_mykey:tfs_mysecret", header)
	})

	t.Run("should handle special characters in keys", func(t *testing.T) {
		creds, err := domain.NewCredentials("tfk_key123", "tfs_secret456")
		require.NoError(t, err)

		header := creds.AuthorizationHeader()
		assert.Contains(t, header, "Bearer")
		assert.Contains(t, header, "tfk_key123")
		assert.Contains(t, header, "tfs_secret456")
	})
}

func TestCredentials_Equals(t *testing.T) {
	t.Run("should return true for equal credentials", func(t *testing.T) {
		creds1, _ := domain.NewCredentials("tfk_same", "tfs_same")
		creds2, _ := domain.NewCredentials("tfk_same", "tfs_same")

		assert.True(t, creds1.Equals(creds2))
		assert.True(t, creds2.Equals(creds1))
	})

	t.Run("should return false for different key IDs", func(t *testing.T) {
		creds1, _ := domain.NewCredentials("tfk_key1", "tfs_same")
		creds2, _ := domain.NewCredentials("tfk_key2", "tfs_same")

		assert.False(t, creds1.Equals(creds2))
	})

	t.Run("should return false for different key secrets", func(t *testing.T) {
		creds1, _ := domain.NewCredentials("tfk_same", "tfs_secret1")
		creds2, _ := domain.NewCredentials("tfk_same", "tfs_secret2")

		assert.False(t, creds1.Equals(creds2))
	})

	t.Run("should return false when comparing to nil", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_key", "tfs_secret")

		assert.False(t, creds.Equals(nil))
	})
}

func TestCredentials_String(t *testing.T) {
	t.Run("should hide secret in string representation", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_visible", "tfs_hidden_secret")

		str := creds.String()
		assert.Contains(t, str, "tfk_visible")
		assert.NotContains(t, str, "tfs_hidden_secret")
		assert.Contains(t, str, "***")
	})
}

// Benchmark tests
func BenchmarkNewCredentials(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = domain.NewCredentials("tfk_benchmark_key", "tfs_benchmark_secret")
	}
}

func BenchmarkCredentials_AuthorizationHeader(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_benchmark_key", "tfs_benchmark_secret")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creds.AuthorizationHeader()
	}
}

func BenchmarkCredentials_Equals(b *testing.B) {
	creds1, _ := domain.NewCredentials("tfk_benchmark_key", "tfs_benchmark_secret")
	creds2, _ := domain.NewCredentials("tfk_benchmark_key", "tfs_benchmark_secret")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = creds1.Equals(creds2)
	}
}
