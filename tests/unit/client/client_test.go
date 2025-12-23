// Package client_test provides unit tests for the client package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package client_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func createTestClient(t *testing.T) *telemetryflow.Client {
	creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)

	config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
	require.NoError(t, err)

	client, err := telemetryflow.NewClient(config)
	require.NoError(t, err)

	return client
}

func TestNewClient(t *testing.T) {
	t.Run("should create client with valid config", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_test", "tfs_secret")
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")

		client, err := telemetryflow.NewClient(config)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.False(t, client.IsInitialized())
	})

	t.Run("should reject nil config", func(t *testing.T) {
		client, err := telemetryflow.NewClient(nil)

		require.Error(t, err)
		assert.Nil(t, client)
		assert.Contains(t, err.Error(), "nil")
	})

	t.Run("should reject invalid config", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_test", "tfs_secret")
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
		// Invalidate config
		config.WithTimeout(-1 * time.Second)

		client, err := telemetryflow.NewClient(config)

		require.Error(t, err)
		assert.Nil(t, client)
	})
}

func TestClient_Config(t *testing.T) {
	t.Run("should return configuration", func(t *testing.T) {
		client := createTestClient(t)

		config := client.Config()

		assert.NotNil(t, config)
		assert.Equal(t, "localhost:4317", config.Endpoint())
		assert.Equal(t, "test-service", config.ServiceName())
	})
}

func TestClient_IsInitialized(t *testing.T) {
	t.Run("should return false before initialization", func(t *testing.T) {
		client := createTestClient(t)

		assert.False(t, client.IsInitialized())
	})
}

func TestClient_Metrics(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// Note: These tests verify that the methods exist and return errors
	// when client is not initialized. Full functionality is tested in integration tests.

	t.Run("RecordMetric should fail when not initialized", func(t *testing.T) {
		err := client.RecordMetric(ctx, "test.metric", 1.0, "count", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("IncrementCounter should fail when not initialized", func(t *testing.T) {
		err := client.IncrementCounter(ctx, "test.counter", 1, nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("RecordGauge should fail when not initialized", func(t *testing.T) {
		err := client.RecordGauge(ctx, "test.gauge", 42.0, nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("RecordHistogram should fail when not initialized", func(t *testing.T) {
		err := client.RecordHistogram(ctx, "test.histogram", 0.5, "s", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestClient_Logs(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	t.Run("Log should fail when not initialized", func(t *testing.T) {
		err := client.Log(ctx, "info", "test message", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogInfo should fail when not initialized", func(t *testing.T) {
		err := client.LogInfo(ctx, "test message", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogWarn should fail when not initialized", func(t *testing.T) {
		err := client.LogWarn(ctx, "test warning", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogError should fail when not initialized", func(t *testing.T) {
		err := client.LogError(ctx, "test error", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestClient_Traces(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	t.Run("StartSpan should fail when not initialized", func(t *testing.T) {
		spanID, err := client.StartSpan(ctx, "test.span", "internal", nil)

		require.Error(t, err)
		assert.Empty(t, spanID)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("EndSpan should fail when not initialized", func(t *testing.T) {
		err := client.EndSpan(ctx, "test-span-id", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("AddSpanEvent should fail when not initialized", func(t *testing.T) {
		err := client.AddSpanEvent(ctx, "test-span-id", "event", nil)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestClient_Flush(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	t.Run("should fail when not initialized", func(t *testing.T) {
		err := client.Flush(ctx)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestClient_Shutdown(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	t.Run("should succeed even when not initialized", func(t *testing.T) {
		// Shutdown on uninitialized client should be a no-op
		err := client.Shutdown(ctx)

		assert.NoError(t, err)
	})
}

func TestClient_MethodSignatures(t *testing.T) {
	// These tests ensure the method signatures match expected patterns
	client := createTestClient(t)
	ctx := context.Background()

	t.Run("RecordMetric accepts all parameters", func(t *testing.T) {
		_ = client.RecordMetric(ctx, "metric.name", 1.5, "unit", map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		})
	})

	t.Run("IncrementCounter accepts int64 value", func(t *testing.T) {
		_ = client.IncrementCounter(ctx, "counter.name", int64(100), map[string]interface{}{
			"label": "test",
		})
	})

	t.Run("RecordGauge accepts float64 value", func(t *testing.T) {
		_ = client.RecordGauge(ctx, "gauge.name", float64(99.9), map[string]interface{}{
			"host": "server-1",
		})
	})

	t.Run("RecordHistogram accepts value and unit", func(t *testing.T) {
		_ = client.RecordHistogram(ctx, "histogram.name", 0.125, "s", map[string]interface{}{
			"endpoint": "/api/test",
		})
	})

	t.Run("Log accepts severity and message", func(t *testing.T) {
		_ = client.Log(ctx, "info", "Test message", map[string]interface{}{
			"user_id": "123",
		})
	})

	t.Run("StartSpan returns span ID", func(t *testing.T) {
		spanID, _ := client.StartSpan(ctx, "test.span", "server", map[string]interface{}{
			"operation": "test",
		})
		// When not initialized, spanID will be empty
		_ = spanID
	})

	t.Run("EndSpan accepts optional error", func(t *testing.T) {
		_ = client.EndSpan(ctx, "span-id", nil)
		_ = client.EndSpan(ctx, "span-id", assert.AnError)
	})
}

// Benchmark tests
func BenchmarkClient_RecordMetric(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_bench", "tfs_bench")
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "bench-service")
	client, _ := telemetryflow.NewClient(config)
	ctx := context.Background()

	// Note: This benchmarks the uninitialized path
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.RecordMetric(ctx, "bench.metric", 1.0, "count", nil)
	}
}

func BenchmarkClient_IncrementCounter(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_bench", "tfs_bench")
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "bench-service")
	client, _ := telemetryflow.NewClient(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.IncrementCounter(ctx, "bench.counter", 1, nil)
	}
}

func BenchmarkClient_LogInfo(b *testing.B) {
	creds, _ := domain.NewCredentials("tfk_bench", "tfs_bench")
	config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "bench-service")
	client, _ := telemetryflow.NewClient(config)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.LogInfo(ctx, "benchmark log message", nil)
	}
}
