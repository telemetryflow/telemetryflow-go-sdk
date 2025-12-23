// Package client_test provides unit tests for the client package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package client_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func TestNewBuilder(t *testing.T) {
	t.Run("should create builder with defaults", func(t *testing.T) {
		builder := telemetryflow.NewBuilder()

		assert.NotNil(t, builder)
	})
}

func TestBuilder_WithAPIKey(t *testing.T) {
	t.Run("should set API key credentials", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test123", "tfs_secret456").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.NoError(t, err)
		assert.NotNil(t, client)

		config := client.Config()
		assert.Equal(t, "tfk_test123", config.Credentials().KeyID())
	})
}

func TestBuilder_WithAPIKeyFromEnv(t *testing.T) {
	t.Run("should read credentials from environment", func(t *testing.T) {
		// Set environment variables
		os.Setenv("TELEMETRYFLOW_API_KEY_ID", "tfk_env_test")
		os.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "tfs_env_secret")
		defer func() {
			os.Unsetenv("TELEMETRYFLOW_API_KEY_ID")
			os.Unsetenv("TELEMETRYFLOW_API_KEY_SECRET")
		}()

		client, err := telemetryflow.NewBuilder().
			WithAPIKeyFromEnv().
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.NoError(t, err)
		assert.NotNil(t, client)

		config := client.Config()
		assert.Equal(t, "tfk_env_test", config.Credentials().KeyID())
	})

	t.Run("should error when env vars not set", func(t *testing.T) {
		// Ensure env vars are not set
		os.Unsetenv("TELEMETRYFLOW_API_KEY_ID")
		os.Unsetenv("TELEMETRYFLOW_API_KEY_SECRET")

		_, err := telemetryflow.NewBuilder().
			WithAPIKeyFromEnv().
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.Error(t, err)
	})
}

func TestBuilder_WithEndpoint(t *testing.T) {
	t.Run("should set custom endpoint", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("custom.endpoint.id:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "custom.endpoint.id:4317", client.Config().Endpoint())
	})
}

func TestBuilder_WithEndpointFromEnv(t *testing.T) {
	t.Run("should read endpoint from environment", func(t *testing.T) {
		os.Setenv("TELEMETRYFLOW_ENDPOINT", "env.endpoint.id:4317")
		defer os.Unsetenv("TELEMETRYFLOW_ENDPOINT")

		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpointFromEnv().
			WithService("test-service", "1.0.0").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "env.endpoint.id:4317", client.Config().Endpoint())
	})

	t.Run("should use default when env not set", func(t *testing.T) {
		os.Unsetenv("TELEMETRYFLOW_ENDPOINT")

		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpointFromEnv().
			WithService("test-service", "1.0.0").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "api.telemetryflow.id:4317", client.Config().Endpoint())
	})
}

func TestBuilder_WithService(t *testing.T) {
	t.Run("should set service name and version", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("my-service", "2.0.0").
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.Equal(t, "my-service", config.ServiceName())
		assert.Equal(t, "2.0.0", config.ServiceVersion())
	})
}

func TestBuilder_WithEnvironment(t *testing.T) {
	t.Run("should set environment", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithEnvironment("staging").
			Build()

		require.NoError(t, err)
		assert.Equal(t, "staging", client.Config().Environment())
	})
}

func TestBuilder_WithProtocol(t *testing.T) {
	t.Run("should set gRPC protocol", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithGRPC().
			Build()

		require.NoError(t, err)
		assert.Equal(t, domain.ProtocolGRPC, client.Config().Protocol())
	})

	t.Run("should set HTTP protocol", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4318").
			WithService("test-service", "1.0.0").
			WithHTTP().
			Build()

		require.NoError(t, err)
		assert.Equal(t, domain.ProtocolHTTP, client.Config().Protocol())
	})
}

func TestBuilder_WithInsecure(t *testing.T) {
	t.Run("should enable insecure mode", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithInsecure(true).
			Build()

		require.NoError(t, err)
		assert.True(t, client.Config().IsInsecure())
	})
}

func TestBuilder_WithTimeout(t *testing.T) {
	t.Run("should set custom timeout", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithTimeout(60 * time.Second).
			Build()

		require.NoError(t, err)
		assert.Equal(t, 60*time.Second, client.Config().Timeout())
	})
}

func TestBuilder_WithSignals(t *testing.T) {
	t.Run("should configure enabled signals", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithSignals(true, false, true). // metrics and traces only
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.True(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.True(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should enable only metrics", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithMetricsOnly().
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.True(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.False(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should enable only logs", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithLogsOnly().
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.False(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.True(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.False(t, config.IsSignalEnabled(domain.SignalTraces))
	})

	t.Run("should enable only traces", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithTracesOnly().
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.False(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.False(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.True(t, config.IsSignalEnabled(domain.SignalTraces))
	})
}

func TestBuilder_WithCustomAttribute(t *testing.T) {
	t.Run("should add custom attributes", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithCustomAttribute("team", "backend").
			WithCustomAttribute("region", "us-east-1").
			Build()

		require.NoError(t, err)
		attrs := client.Config().CustomAttributes()
		assert.Equal(t, "backend", attrs["team"])
		assert.Equal(t, "us-east-1", attrs["region"])
	})
}

func TestBuilder_WithAutoConfiguration(t *testing.T) {
	t.Run("should configure from all environment variables", func(t *testing.T) {
		// Set all environment variables
		os.Setenv("TELEMETRYFLOW_API_KEY_ID", "tfk_auto_test")
		os.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "tfs_auto_secret")
		os.Setenv("TELEMETRYFLOW_ENDPOINT", "auto.endpoint.id:4317")
		os.Setenv("TELEMETRYFLOW_SERVICE_NAME", "auto-service")
		os.Setenv("TELEMETRYFLOW_SERVICE_VERSION", "3.0.0")
		os.Setenv("ENV", "production")
		defer func() {
			os.Unsetenv("TELEMETRYFLOW_API_KEY_ID")
			os.Unsetenv("TELEMETRYFLOW_API_KEY_SECRET")
			os.Unsetenv("TELEMETRYFLOW_ENDPOINT")
			os.Unsetenv("TELEMETRYFLOW_SERVICE_NAME")
			os.Unsetenv("TELEMETRYFLOW_SERVICE_VERSION")
			os.Unsetenv("ENV")
		}()

		client, err := telemetryflow.NewBuilder().
			WithAutoConfiguration().
			Build()

		require.NoError(t, err)
		config := client.Config()
		assert.Equal(t, "tfk_auto_test", config.Credentials().KeyID())
		assert.Equal(t, "auto.endpoint.id:4317", config.Endpoint())
		assert.Equal(t, "auto-service", config.ServiceName())
		assert.Equal(t, "production", config.Environment())
	})
}

func TestBuilder_Build(t *testing.T) {
	t.Run("should fail without API key", func(t *testing.T) {
		_, err := telemetryflow.NewBuilder().
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "credentials")
	})

	t.Run("should fail without endpoint", func(t *testing.T) {
		_, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithService("test-service", "1.0.0").
			Build()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "endpoint")
	})

	t.Run("should fail without service name", func(t *testing.T) {
		_, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			Build()

		require.Error(t, err)
		assert.Contains(t, err.Error(), "service name")
	})

	t.Run("should fail with invalid API key format", func(t *testing.T) {
		_, err := telemetryflow.NewBuilder().
			WithAPIKey("invalid_key", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			Build()

		require.Error(t, err)
	})
}

func TestBuilder_MustBuild(t *testing.T) {
	t.Run("should panic on error", func(t *testing.T) {
		assert.Panics(t, func() {
			telemetryflow.NewBuilder().
				WithEndpoint("localhost:4317").
				MustBuild()
		})
	})

	t.Run("should return client on success", func(t *testing.T) {
		assert.NotPanics(t, func() {
			client := telemetryflow.NewBuilder().
				WithAPIKey("tfk_test", "tfs_secret").
				WithEndpoint("localhost:4317").
				WithService("test-service", "1.0.0").
				MustBuild()

			assert.NotNil(t, client)
		})
	})
}

func TestBuilder_MethodChaining(t *testing.T) {
	t.Run("should support full method chaining", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_test", "tfs_secret").
			WithEndpoint("localhost:4317").
			WithService("test-service", "1.0.0").
			WithEnvironment("production").
			WithGRPC().
			WithInsecure(false).
			WithTimeout(30*time.Second).
			WithSignals(true, true, true).
			WithCustomAttribute("team", "backend").
			WithCustomAttribute("region", "us-east-1").
			Build()

		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

// Convenience constructors tests

func TestNewFromEnv(t *testing.T) {
	t.Run("should create client from environment", func(t *testing.T) {
		os.Setenv("TELEMETRYFLOW_API_KEY_ID", "tfk_env_client")
		os.Setenv("TELEMETRYFLOW_API_KEY_SECRET", "tfs_env_secret")
		os.Setenv("TELEMETRYFLOW_SERVICE_NAME", "env-service")
		defer func() {
			os.Unsetenv("TELEMETRYFLOW_API_KEY_ID")
			os.Unsetenv("TELEMETRYFLOW_API_KEY_SECRET")
			os.Unsetenv("TELEMETRYFLOW_SERVICE_NAME")
		}()

		client, err := telemetryflow.NewFromEnv()

		require.NoError(t, err)
		assert.NotNil(t, client)
	})
}

func TestNewSimple(t *testing.T) {
	t.Run("should create client with minimal config", func(t *testing.T) {
		client, err := telemetryflow.NewSimple(
			"tfk_simple",
			"tfs_simple",
			"localhost:4317",
			"simple-service",
		)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "simple-service", client.Config().ServiceName())
	})
}

// Benchmark tests
func BenchmarkBuilder_Build(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = telemetryflow.NewBuilder().
			WithAPIKey("tfk_bench", "tfs_bench").
			WithEndpoint("localhost:4317").
			WithService("bench-service", "1.0.0").
			Build()
	}
}
