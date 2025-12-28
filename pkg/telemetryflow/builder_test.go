// Package telemetryflow provides unit tests for the telemetryflow package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package telemetryflow

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func TestNewBuilder(t *testing.T) {
	t.Run("should create builder with defaults", func(t *testing.T) {
		builder := NewBuilder()

		assert.NotNil(t, builder)
		assert.Equal(t, domain.ProtocolGRPC, builder.protocol)
		assert.False(t, builder.insecure)
		assert.Equal(t, 30*time.Second, builder.timeout)
		assert.True(t, builder.enableMetrics)
		assert.True(t, builder.enableLogs)
		assert.True(t, builder.enableTraces)
		assert.True(t, builder.enableExemplars)
		assert.Equal(t, "telemetryflow", builder.serviceNamespace)
		assert.NotNil(t, builder.customAttrs)
		assert.Empty(t, builder.errors)
	})
}

func TestBuilder_WithAPIKey(t *testing.T) {
	t.Run("should set API key credentials", func(t *testing.T) {
		builder := NewBuilder().WithAPIKey("tfk_test", "tfs_secret")

		assert.Equal(t, "tfk_test", builder.apiKeyID)
		assert.Equal(t, "tfs_secret", builder.apiKeySecret)
	})
}

func TestBuilder_WithEndpoint(t *testing.T) {
	t.Run("should set endpoint", func(t *testing.T) {
		builder := NewBuilder().WithEndpoint("localhost:4317")

		assert.Equal(t, "localhost:4317", builder.endpoint)
	})
}

func TestBuilder_WithService(t *testing.T) {
	t.Run("should set service name and version", func(t *testing.T) {
		builder := NewBuilder().WithService("my-service", "1.0.0")

		assert.Equal(t, "my-service", builder.serviceName)
		assert.Equal(t, "1.0.0", builder.serviceVersion)
	})
}

func TestBuilder_WithEnvironment(t *testing.T) {
	t.Run("should set environment", func(t *testing.T) {
		builder := NewBuilder().WithEnvironment("staging")

		assert.Equal(t, "staging", builder.environment)
	})
}

func TestBuilder_WithProtocol(t *testing.T) {
	t.Run("should set gRPC protocol", func(t *testing.T) {
		builder := NewBuilder().WithGRPC()

		assert.Equal(t, domain.ProtocolGRPC, builder.protocol)
	})

	t.Run("should set HTTP protocol", func(t *testing.T) {
		builder := NewBuilder().WithHTTP()

		assert.Equal(t, domain.ProtocolHTTP, builder.protocol)
	})
}

func TestBuilder_WithInsecure(t *testing.T) {
	t.Run("should enable insecure mode", func(t *testing.T) {
		builder := NewBuilder().WithInsecure(true)

		assert.True(t, builder.insecure)
	})

	t.Run("should disable insecure mode", func(t *testing.T) {
		builder := NewBuilder().WithInsecure(false)

		assert.False(t, builder.insecure)
	})
}

func TestBuilder_WithTimeout(t *testing.T) {
	t.Run("should set custom timeout", func(t *testing.T) {
		builder := NewBuilder().WithTimeout(60 * time.Second)

		assert.Equal(t, 60*time.Second, builder.timeout)
	})
}

func TestBuilder_WithSignals(t *testing.T) {
	t.Run("should enable/disable signals", func(t *testing.T) {
		builder := NewBuilder().WithSignals(true, false, true)

		assert.True(t, builder.enableMetrics)
		assert.False(t, builder.enableLogs)
		assert.True(t, builder.enableTraces)
	})

	t.Run("should enable only metrics", func(t *testing.T) {
		builder := NewBuilder().WithMetricsOnly()

		assert.True(t, builder.enableMetrics)
		assert.False(t, builder.enableLogs)
		assert.False(t, builder.enableTraces)
	})

	t.Run("should enable only logs", func(t *testing.T) {
		builder := NewBuilder().WithLogsOnly()

		assert.False(t, builder.enableMetrics)
		assert.True(t, builder.enableLogs)
		assert.False(t, builder.enableTraces)
	})

	t.Run("should enable only traces", func(t *testing.T) {
		builder := NewBuilder().WithTracesOnly()

		assert.False(t, builder.enableMetrics)
		assert.False(t, builder.enableLogs)
		assert.True(t, builder.enableTraces)
	})
}

func TestBuilder_WithCustomAttribute(t *testing.T) {
	t.Run("should add custom attribute", func(t *testing.T) {
		builder := NewBuilder().
			WithCustomAttribute("key1", "value1").
			WithCustomAttribute("key2", "value2")

		assert.Equal(t, "value1", builder.customAttrs["key1"])
		assert.Equal(t, "value2", builder.customAttrs["key2"])
	})
}

func TestBuilder_WithCollectorID(t *testing.T) {
	t.Run("should set collector ID", func(t *testing.T) {
		builder := NewBuilder().WithCollectorID("collector-123")

		assert.Equal(t, "collector-123", builder.collectorID)
	})
}

func TestBuilder_WithServiceNamespace(t *testing.T) {
	t.Run("should set service namespace", func(t *testing.T) {
		builder := NewBuilder().WithServiceNamespace("custom-namespace")

		assert.Equal(t, "custom-namespace", builder.serviceNamespace)
	})
}

func TestBuilder_WithExemplars(t *testing.T) {
	t.Run("should enable exemplars", func(t *testing.T) {
		builder := NewBuilder().WithExemplars(true)

		assert.True(t, builder.enableExemplars)
	})

	t.Run("should disable exemplars", func(t *testing.T) {
		builder := NewBuilder().WithExemplars(false)

		assert.False(t, builder.enableExemplars)
	})
}

func TestBuilder_WithAPIKeyFromEnv(t *testing.T) {
	t.Run("should read API key from environment", func(t *testing.T) {
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "tfk_env")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "tfs_env_secret")

		builder := NewBuilder().WithAPIKeyFromEnv()

		assert.Equal(t, "tfk_env", builder.apiKeyID)
		assert.Equal(t, "tfs_env_secret", builder.apiKeySecret)
		assert.Empty(t, builder.errors)
	})

	t.Run("should add error when env vars not set", func(t *testing.T) {
		// Ensure env vars are not set by using empty values
		t.Setenv("TELEMETRYFLOW_API_KEY_ID", "")
		t.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "")

		builder := NewBuilder().WithAPIKeyFromEnv()

		assert.Len(t, builder.errors, 2)
	})
}

func TestBuilder_Build(t *testing.T) {
	t.Run("should build client with valid config", func(t *testing.T) {
		client, err := NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithInsecure(true).
			Build()

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "test-service", client.Config().ServiceName())
	})

	t.Run("should fail without API credentials", func(t *testing.T) {
		_, err := NewBuilder().
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "API credentials are required")
	})

	t.Run("should fail without endpoint", func(t *testing.T) {
		_, err := NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithService("test-service", "1.0.0").
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint is required")
	})

	t.Run("should fail without service name", func(t *testing.T) {
		_, err := NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			Build()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service name is required")
	})
}

func TestBuilder_MustBuild(t *testing.T) {
	t.Run("should panic on invalid config", func(t *testing.T) {
		assert.Panics(t, func() {
			NewBuilder().MustBuild()
		})
	})
}

func TestNewSimple(t *testing.T) {
	t.Run("should create client with simple config", func(t *testing.T) {
		client, err := NewSimple("tfk_test", "tfs_secret", "localhost:4317", "simple-service")

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "simple-service", client.Config().ServiceName())
		assert.Equal(t, "1.0.0", client.Config().ServiceVersion())
	})
}
