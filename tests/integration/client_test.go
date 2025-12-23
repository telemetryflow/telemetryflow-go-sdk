// Package integration provides integration tests for the TelemetryFlow SDK.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

// Skip integration tests in short mode
func skipInShortMode(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}

// Check if required environment variables are set for integration tests
func hasIntegrationEnv() bool {
	return os.Getenv("TELEMETRYFLOW_API_KEY_ID") != "" &&
		os.Getenv("TELEMETRYFLOW_API_KEY_SECRET") != ""
}

// Check if OTLP collector is available for integration tests
func hasOTLPCollector() bool {
	return os.Getenv("OTLP_COLLECTOR_AVAILABLE") == "true"
}

func TestClientLifecycle(t *testing.T) {
	skipInShortMode(t)

	t.Run("should initialize and shutdown gracefully", func(t *testing.T) {
		// Skip if no OTLP collector is available (CI environment)
		if !hasOTLPCollector() {
			t.Skip("Skipping: OTLP_COLLECTOR_AVAILABLE not set (no collector in CI)")
		}

		creds, err := domain.NewCredentials("tfk_test123", "tfs_secret456")
		require.NoError(t, err)

		config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "integration-test")
		require.NoError(t, err)
		config.WithInsecure(true) // For local testing

		client, err := telemetryflow.NewClient(config)
		require.NoError(t, err)

		ctx := context.Background()

		// Client should not be initialized yet
		assert.False(t, client.IsInitialized())

		// Note: Initialize will fail without actual OTLP endpoint
		// This test just verifies the lifecycle methods exist and work
		err = client.Initialize(ctx)
		// We expect this to potentially fail in test environment without collector
		if err != nil {
			t.Skipf("Skipping: OTLP collector not available at localhost:4317 (%v)", err)
			return
		}

		assert.True(t, client.IsInitialized())

		// Shutdown
		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		err = client.Shutdown(shutdownCtx)
		assert.NoError(t, err)
		assert.False(t, client.IsInitialized())
	})

	t.Run("should not initialize twice", func(t *testing.T) {
		// Skip if no OTLP collector is available (CI environment)
		if !hasOTLPCollector() {
			t.Skip("Skipping: OTLP_COLLECTOR_AVAILABLE not set (no collector in CI)")
		}

		creds, _ := domain.NewCredentials("tfk_test", "tfs_secret")
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
		config.WithInsecure(true)

		client, _ := telemetryflow.NewClient(config)
		ctx := context.Background()

		// First initialization attempt
		err1 := client.Initialize(ctx)

		if err1 == nil {
			// Second initialization should fail
			err2 := client.Initialize(ctx)
			assert.Error(t, err2)
			assert.Contains(t, err2.Error(), "already initialized")

			// Cleanup
			_ = client.Shutdown(ctx)
		}
	})
}

func TestClientWithBuilder(t *testing.T) {
	skipInShortMode(t)

	t.Run("should create fully configured client", func(t *testing.T) {
		client, err := telemetryflow.NewBuilder().
			WithAPIKey("tfk_integration", "tfs_integration").
			WithEndpoint("localhost:4317").
			WithService("integration-service", "1.0.0").
			WithEnvironment("test").
			WithGRPC().
			WithInsecure(true).
			WithTimeout(10*time.Second).
			WithSignals(true, true, true).
			WithCustomAttribute("test", "integration").
			Build()

		require.NoError(t, err)
		assert.NotNil(t, client)

		config := client.Config()
		assert.Equal(t, "localhost:4317", config.Endpoint())
		assert.Equal(t, "integration-service", config.ServiceName())
		assert.Equal(t, "1.0.0", config.ServiceVersion())
		assert.Equal(t, "test", config.Environment())
		assert.Equal(t, domain.ProtocolGRPC, config.Protocol())
		assert.True(t, config.IsInsecure())
		assert.Equal(t, 10*time.Second, config.Timeout())
		assert.True(t, config.IsSignalEnabled(domain.SignalMetrics))
		assert.True(t, config.IsSignalEnabled(domain.SignalLogs))
		assert.True(t, config.IsSignalEnabled(domain.SignalTraces))
		assert.Equal(t, "integration", config.CustomAttributes()["test"])
	})
}

func TestClientMetricsIntegration(t *testing.T) {
	skipInShortMode(t)
	if !hasIntegrationEnv() {
		t.Skip("Skipping: TELEMETRYFLOW_API_KEY_ID and TELEMETRYFLOW_API_KEY_SECRET not set")
	}

	t.Run("should send metrics when initialized", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() { _ = client.Shutdown(ctx) }()

		// Send various metrics
		err = client.IncrementCounter(ctx, "integration.counter", 1, map[string]interface{}{
			"test": "true",
		})
		assert.NoError(t, err)

		err = client.RecordGauge(ctx, "integration.gauge", 42.5, map[string]interface{}{
			"unit": "percent",
		})
		assert.NoError(t, err)

		err = client.RecordHistogram(ctx, "integration.histogram", 0.125, "s", map[string]interface{}{
			"operation": "test",
		})
		assert.NoError(t, err)

		err = client.RecordMetric(ctx, "integration.custom", 100.0, "items", map[string]interface{}{
			"category": "test",
		})
		assert.NoError(t, err)

		// Flush to ensure data is sent
		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestClientLogsIntegration(t *testing.T) {
	skipInShortMode(t)
	if !hasIntegrationEnv() {
		t.Skip("Skipping: TELEMETRYFLOW_API_KEY_ID and TELEMETRYFLOW_API_KEY_SECRET not set")
	}

	t.Run("should send logs when initialized", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() { _ = client.Shutdown(ctx) }()

		// Send logs at different levels
		err = client.LogInfo(ctx, "Integration test info log", map[string]interface{}{
			"test":  "true",
			"level": "info",
		})
		assert.NoError(t, err)

		err = client.LogWarn(ctx, "Integration test warning log", map[string]interface{}{
			"test":  "true",
			"level": "warn",
		})
		assert.NoError(t, err)

		err = client.LogError(ctx, "Integration test error log", map[string]interface{}{
			"test":  "true",
			"level": "error",
		})
		assert.NoError(t, err)

		err = client.Log(ctx, "debug", "Integration test debug log", map[string]interface{}{
			"test": "true",
		})
		assert.NoError(t, err)

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestClientTracesIntegration(t *testing.T) {
	skipInShortMode(t)
	if !hasIntegrationEnv() {
		t.Skip("Skipping: TELEMETRYFLOW_API_KEY_ID and TELEMETRYFLOW_API_KEY_SECRET not set")
	}

	t.Run("should create and end spans when initialized", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() { _ = client.Shutdown(ctx) }()

		// Create a span
		spanID, err := client.StartSpan(ctx, "integration.test.span", "internal", map[string]interface{}{
			"test": "true",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, spanID)

		// Add event to span
		err = client.AddSpanEvent(ctx, spanID, "test.event", map[string]interface{}{
			"step": "processing",
		})
		assert.NoError(t, err)

		// Simulate work
		time.Sleep(50 * time.Millisecond)

		// End span successfully
		err = client.EndSpan(ctx, spanID, nil)
		assert.NoError(t, err)

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("should record span with error", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() { _ = client.Shutdown(ctx) }()

		spanID, err := client.StartSpan(ctx, "integration.error.span", "internal", nil)
		assert.NoError(t, err)

		// End span with error
		testErr := assert.AnError
		err = client.EndSpan(ctx, spanID, testErr)
		assert.NoError(t, err)

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestClientConcurrency(t *testing.T) {
	skipInShortMode(t)

	t.Run("should handle concurrent operations", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_concurrent", "tfs_concurrent")
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "concurrent-test")
		config.WithInsecure(true)

		client, err := telemetryflow.NewClient(config)
		require.NoError(t, err)

		// Even without initialization, concurrent calls should not panic
		ctx := context.Background()
		done := make(chan bool, 100)

		// Launch concurrent operations
		for i := 0; i < 100; i++ {
			go func(id int) {
				_ = client.IncrementCounter(ctx, "concurrent.counter", 1, map[string]interface{}{
					"goroutine": id,
				})
				done <- true
			}(i)
		}

		// Wait for all goroutines
		for i := 0; i < 100; i++ {
			<-done
		}
	})

	t.Run("should safely check initialization concurrently", func(t *testing.T) {
		creds, _ := domain.NewCredentials("tfk_test", "tfs_test")
		config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
		client, _ := telemetryflow.NewClient(config)

		done := make(chan bool, 100)

		for i := 0; i < 100; i++ {
			go func() {
				_ = client.IsInitialized()
				done <- true
			}()
		}

		for i := 0; i < 100; i++ {
			<-done
		}
	})
}

func TestSignalConfiguration(t *testing.T) {
	skipInShortMode(t)

	testCases := []struct {
		name    string
		metrics bool
		logs    bool
		traces  bool
	}{
		{"metrics only", true, false, false},
		{"logs only", false, true, false},
		{"traces only", false, false, true},
		{"metrics and logs", true, true, false},
		{"metrics and traces", true, false, true},
		{"logs and traces", false, true, true},
		{"all signals", true, true, true},
		{"no signals", false, false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := telemetryflow.NewBuilder().
				WithAPIKey("tfk_signal_test", "tfs_signal_test").
				WithEndpoint("localhost:4317").
				WithService("signal-test", "1.0.0").
				WithSignals(tc.metrics, tc.logs, tc.traces).
				Build()

			require.NoError(t, err)

			config := client.Config()
			assert.Equal(t, tc.metrics, config.IsSignalEnabled(domain.SignalMetrics))
			assert.Equal(t, tc.logs, config.IsSignalEnabled(domain.SignalLogs))
			assert.Equal(t, tc.traces, config.IsSignalEnabled(domain.SignalTraces))
		})
	}
}

// Benchmark tests for integration scenarios
func BenchmarkClientCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = telemetryflow.NewBuilder().
			WithAPIKey("tfk_bench", "tfs_bench").
			WithEndpoint("localhost:4317").
			WithService("bench-service", "1.0.0").
			Build()
	}
}
