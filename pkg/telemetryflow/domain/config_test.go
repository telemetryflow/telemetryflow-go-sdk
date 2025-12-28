// Package domain provides unit tests for the domain layer config.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidCredentials(t *testing.T) *Credentials {
	creds, err := NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)
	return creds
}

func TestNewTelemetryConfig(t *testing.T) {
	t.Run("should create config with valid inputs", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := NewTelemetryConfig(creds, "api.telemetryflow.io:4317", "test-service")

		require.NoError(t, err)
		assert.NotNil(t, config)
		assert.Equal(t, "api.telemetryflow.io:4317", config.Endpoint())
		assert.Equal(t, "test-service", config.ServiceName())
	})

	t.Run("should set correct defaults", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, err := NewTelemetryConfig(creds, "localhost:4317", "my-service")

		require.NoError(t, err)
		assert.Equal(t, ProtocolGRPC, config.Protocol())
		assert.False(t, config.IsInsecure())
		assert.Equal(t, 30*time.Second, config.Timeout())
		assert.True(t, config.IsRetryEnabled())
		assert.Equal(t, 3, config.MaxRetries())
		assert.Equal(t, 5*time.Second, config.RetryBackoff())
		assert.True(t, config.IsCompressionEnabled())
		assert.Equal(t, "1.0.0", config.ServiceVersion())
		assert.Equal(t, "production", config.Environment())
	})

	t.Run("should fail with nil credentials", func(t *testing.T) {
		_, err := NewTelemetryConfig(nil, "localhost:4317", "my-service")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "credentials")
	})

	t.Run("should fail with empty endpoint", func(t *testing.T) {
		creds := createValidCredentials(t)
		_, err := NewTelemetryConfig(creds, "", "my-service")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint")
	})

	t.Run("should fail with empty service name", func(t *testing.T) {
		creds := createValidCredentials(t)
		_, err := NewTelemetryConfig(creds, "localhost:4317", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service name")
	})
}

func TestTelemetryConfig_FluentSetters(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

	t.Run("should set protocol", func(t *testing.T) {
		config.WithProtocol(ProtocolHTTP)
		assert.Equal(t, ProtocolHTTP, config.Protocol())
	})

	t.Run("should set insecure mode", func(t *testing.T) {
		config.WithInsecure(true)
		assert.True(t, config.IsInsecure())
	})

	t.Run("should set timeout", func(t *testing.T) {
		config.WithTimeout(60 * time.Second)
		assert.Equal(t, 60*time.Second, config.Timeout())
	})

	t.Run("should set service version", func(t *testing.T) {
		config.WithServiceVersion("2.0.0")
		assert.Equal(t, "2.0.0", config.ServiceVersion())
	})

	t.Run("should set service namespace", func(t *testing.T) {
		config.WithServiceNamespace("custom-ns")
		assert.Equal(t, "custom-ns", config.ServiceNamespace())
	})

	t.Run("should set environment", func(t *testing.T) {
		config.WithEnvironment("staging")
		assert.Equal(t, "staging", config.Environment())
	})

	t.Run("should set collector ID", func(t *testing.T) {
		config.WithCollectorID("collector-123")
		assert.Equal(t, "collector-123", config.CollectorID())
	})

	t.Run("should set batch settings", func(t *testing.T) {
		config.WithBatchSettings(20*time.Second, 256)
		assert.Equal(t, 20*time.Second, config.BatchTimeout())
		assert.Equal(t, 256, config.BatchMaxSize())
	})

	t.Run("should set exemplars enabled", func(t *testing.T) {
		config.WithExemplars(true)
		assert.True(t, config.IsExemplarsEnabled())
	})
}

func TestTelemetryConfig_Signals(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

	t.Run("should enable all signals by default", func(t *testing.T) {
		assert.True(t, config.IsSignalEnabled(SignalTraces))
		assert.True(t, config.IsSignalEnabled(SignalMetrics))
		assert.True(t, config.IsSignalEnabled(SignalLogs))
	})

	t.Run("should disable specific signals", func(t *testing.T) {
		config.WithSignals(true, false, true)

		assert.True(t, config.IsSignalEnabled(SignalMetrics))
		assert.False(t, config.IsSignalEnabled(SignalLogs))
		assert.True(t, config.IsSignalEnabled(SignalTraces))
	})
}

func TestTelemetryConfig_CustomAttributes(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

	t.Run("should add custom attributes", func(t *testing.T) {
		config.WithCustomAttribute("key1", "value1")
		config.WithCustomAttribute("key2", "value2")

		attrs := config.CustomAttributes()
		assert.Equal(t, "value1", attrs["key1"])
		assert.Equal(t, "value2", attrs["key2"])
	})
}

func TestTelemetryConfig_Retry(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

	t.Run("should configure retry settings", func(t *testing.T) {
		config.WithRetry(true, 5, 10*time.Second)

		assert.True(t, config.IsRetryEnabled())
		assert.Equal(t, 5, config.MaxRetries())
		assert.Equal(t, 10*time.Second, config.RetryBackoff())
	})

	t.Run("should disable retry", func(t *testing.T) {
		config.WithRetry(false, 0, 0)

		assert.False(t, config.IsRetryEnabled())
	})
}

func TestTelemetryConfig_Compression(t *testing.T) {
	creds := createValidCredentials(t)
	config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

	t.Run("should enable compression by default", func(t *testing.T) {
		assert.True(t, config.IsCompressionEnabled())
	})

	t.Run("should disable compression", func(t *testing.T) {
		config.WithCompression(false)
		assert.False(t, config.IsCompressionEnabled())
	})
}

func TestTelemetryConfig_Validate(t *testing.T) {
	t.Run("should validate successfully with valid config", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

		err := config.Validate()
		assert.NoError(t, err)
	})
}

func TestTelemetryConfig_Credentials(t *testing.T) {
	t.Run("should return credentials", func(t *testing.T) {
		creds := createValidCredentials(t)
		config, _ := NewTelemetryConfig(creds, "localhost:4317", "test-service")

		assert.Equal(t, creds, config.Credentials())
	})
}

func TestProtocol(t *testing.T) {
	t.Run("should define gRPC protocol", func(t *testing.T) {
		assert.Equal(t, Protocol("grpc"), ProtocolGRPC)
	})

	t.Run("should define HTTP protocol", func(t *testing.T) {
		assert.Equal(t, Protocol("http"), ProtocolHTTP)
	})
}

func TestSignalType(t *testing.T) {
	t.Run("should define trace signal", func(t *testing.T) {
		assert.Equal(t, SignalType("traces"), SignalTraces)
	})

	t.Run("should define metrics signal", func(t *testing.T) {
		assert.Equal(t, SignalType("metrics"), SignalMetrics)
	})

	t.Run("should define logs signal", func(t *testing.T) {
		assert.Equal(t, SignalType("logs"), SignalLogs)
	})
}
