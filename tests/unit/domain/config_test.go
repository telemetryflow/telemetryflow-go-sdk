// Package domain_test provides unit tests for the domain package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func createValidCredentials(t *testing.T) *domain.Credentials {
	creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)
	return creds
}

func TestNewTelemetryConfig(t *testing.T) {
	t.Run("should create config with valid inputs", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := domain.NewTelemetryConfig(creds, "api.telemetryflow.io:4317", "test-service")

		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "api.telemetryflow.io:4317", config.Endpoint())
		assert.Equal(t, "test-service", config.ServiceName())
	})

	t.Run("should set correct defaults", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

		require.NoError(t, err)
		assert.Equal(t, domain.ProtocolGRPC, config.Protocol())
		assert.False(t, config.IsInsecure())
		assert.Equal(t, 30*time.Second, config.Timeout())
		assert.True(t, config.IsRetryEnabled())
		assert.Equal(t, 3, config.MaxRetries())
		assert.Equal(t, 5*time.Second, config.RetryBackoff())
		assert.True(t, config.IsCompressionEnabled())
		assert.Equal(t, "1.0.0", config.ServiceVersion())
		assert.Equal(t, "production", config.Environment())
		assert.Equal(t, 10*time.Second, config.BatchTimeout())
		assert.Equal(t, 512, config.BatchMaxSize())
		assert.Equal(t, 1000, config.RateLimit())
	})

	t.Run("should enable all signals by default", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

		assert.True(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.True(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.True(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should reject nil credentials", func(t *testing.T) {
		config, err := domain.NewTelemetryConfig(nil, "localhost:4317", "my-service")

		require.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "credentials")
	})

	t.Run("should reject empty endpoint", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := domain.NewTelemetryConfig(creds, "", "my-service")

		require.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "endpoint")
	})

	t.Run("should reject empty service name", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "")

		require.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "service name")
	})
}

func TestTelemetryConfig_WithProtocol(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should set gRPC protocol", func(t *testing.T) {
		config.WithProtocol(domain.ProtocolGRPC)
		assert.Equal(t, domain.ProtocolGRPC, config.Protocol())
	})

	t.Run("should set HTTP protocol", func(t *testing.T) {
		config.WithProtocol(domain.ProtocolHTTP)
		assert.Equal(t, domain.ProtocolHTTP, config.Protocol())
	})

	t.Run("should return config for chaining", func(t *testing.T) {
		result := config.WithProtocol(domain.ProtocolGRPC)
		assert.Same(t, config, result)
	})
}

func TestTelemetryConfig_WithInsecure(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should enable insecure mode", func(t *testing.T) {
		config.WithInsecure(true)
		assert.True(t, config.IsInsecure())
	})

	t.Run("should disable insecure mode", func(t *testing.T) {
		config.WithInsecure(false)
		assert.False(t, config.IsInsecure())
	})
}

func TestTelemetryConfig_WithTimeout(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should set custom timeout", func(t *testing.T) {
		config.WithTimeout(60 * time.Second)
		assert.Equal(t, 60*time.Second, config.Timeout())
	})

	t.Run("should accept short timeout", func(t *testing.T) {
		config.WithTimeout(1 * time.Second)
		assert.Equal(t, 1*time.Second, config.Timeout())
	})
}

func TestTelemetryConfig_WithRetry(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should configure retry settings", func(t *testing.T) {
		config.WithRetry(true, 5, 10*time.Second)

		assert.True(t, config.IsRetryEnabled())
		assert.Equal(t, 5, config.MaxRetries())
		assert.Equal(t, 10*time.Second, config.RetryBackoff())
	})

	t.Run("should disable retry", func(t *testing.T) {
		config.WithRetry(false, 0, 0)

		assert.False(t, config.IsRetryEnabled())
		assert.Equal(t, 0, config.MaxRetries())
	})
}

func TestTelemetryConfig_WithSignals(t *testing.T) {
	creds := createValidCredentials(t)

	t.Run("should enable only metrics", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithSignals(true, false, false)

		assert.True(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.False(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should enable only logs", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithSignals(false, true, false)

		assert.False(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.True(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.False(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should enable only traces", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithSignals(false, false, true)

		assert.False(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.True(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should disable all signals", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithSignals(false, false, false)

		assert.False(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.False(t, config.IsSignalEnabled(domain.SignalTraces))
	})
}

func TestTelemetryConfig_WithCustomAttribute(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should add custom attributes", func(t *testing.T) {
		config.
			WithCustomAttribute("team", "backend").
			WithCustomAttribute("region", "us-east-1").
			WithCustomAttribute("datacenter", "dc1")

		attrs := config.CustomAttributes()
		assert.Equal(t, "backend", attrs["team"])
		assert.Equal(t, "us-east-1", attrs["region"])
		assert.Equal(t, "dc1", attrs["datacenter"])
	})

	t.Run("should overwrite existing attribute", func(t *testing.T) {
		config.WithCustomAttribute("team", "backend")
		config.WithCustomAttribute("team", "frontend")

		attrs := config.CustomAttributes()
		assert.Equal(t, "frontend", attrs["team"])
	})
}

func TestTelemetryConfig_WithBatchSettings(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should configure batch settings", func(t *testing.T) {
		config.WithBatchSettings(30*time.Second, 1024)

		assert.Equal(t, 30*time.Second, config.BatchTimeout())
		assert.Equal(t, 1024, config.BatchMaxSize())
	})
}

func TestTelemetryConfig_WithRateLimit(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

	t.Run("should configure rate limit", func(t *testing.T) {
		config.WithRateLimit(5000)

		assert.Equal(t, 5000, config.RateLimit())
	})
}

func TestTelemetryConfig_Validate(t *testing.T) {
	creds := createValidCredentials(t)

	t.Run("should pass validation for valid config", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("should fail validation for negative timeout", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithTimeout(-1 * time.Second)

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "timeout")
	})

	t.Run("should fail validation for negative max retries", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithRetry(true, -1, 5*time.Second)

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "retries")
	})

	t.Run("should fail validation for zero batch size", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithBatchSettings(10*time.Second, 0)

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "batch")
	})

	t.Run("should fail validation for negative rate limit", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")
		config.WithRateLimit(-1)

		err := config.Validate()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit")
	})
}

func TestTelemetryConfig_String(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := domain.NewTelemetryConfig(creds, "api.telemetryflow.io:4317", "my-service")
	config.WithEnvironment("staging")

	t.Run("should return string representation", func(t *testing.T) {
		str := config.String()

		assert.Contains(t, str, "api.telemetryflow.io:4317")
		assert.Contains(t, str, "my-service")
		assert.Contains(t, str, "staging")
	})
}

func TestTelemetryConfig_BuilderChaining(t *testing.T) {
	creds := createValidCredentials(t)

	t.Run("should support method chaining", func(t *testing.T) {
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

		result := config.
			WithProtocol(domain.ProtocolHTTP).
			WithInsecure(true).
			WithTimeout(60 * time.Second).
			WithServiceVersion("2.0.0").
			WithEnvironment("staging").
			WithRetry(true, 5, 10*time.Second).
			WithCompression(false).
			WithBatchSettings(30*time.Second, 1024).
			WithRateLimit(5000).
			WithCustomAttribute("team", "backend")

		assert.Same(t, config, result)
		assert.Equal(t, domain.ProtocolHTTP, config.Protocol())
		assert.True(t, config.IsInsecure())
		assert.Equal(t, 60*time.Second, config.Timeout())
		assert.Equal(t, "2.0.0", config.ServiceVersion())
		assert.Equal(t, "staging", config.Environment())
		assert.Equal(t, 5, config.MaxRetries())
		assert.False(t, config.IsCompressionEnabled())
		assert.Equal(t, 1024, config.BatchMaxSize())
		assert.Equal(t, 5000, config.RateLimit())
	})
}

// Benchmark tests
func BenchmarkNewTelemetryConfig(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_bench", "tfs_bench")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = domain.NewTelemetryConfig(creds, "localhost:4317", "bench-service")
	}
}

func BenchmarkTelemetryConfig_Validate(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_bench", "tfs_bench")
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "bench-service")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = config.Validate()
	}
}
