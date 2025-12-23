// Package e2e provides end-to-end tests for the TelemetryFlow SDK.
//
// These tests require a running TelemetryFlow backend or OTLP collector.
// Set the following environment variables to run:
//   - TELEMETRYFLOW_API_KEY_ID
//   - TELEMETRYFLOW_API_KEY_SECRET
//   - TELEMETRYFLOW_ENDPOINT (optional, defaults to api.telemetryflow.io:4317)
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package e2e

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

// skipE2E skips the test if e2e environment is not configured
func skipE2E(t *testing.T) {
	if os.Getenv("TELEMETRYFLOW_E2E") != "true" {
		t.Skip("Skipping e2e test: set TELEMETRYFLOW_E2E=true to run")
	}
	if os.Getenv("TELEMETRYFLOW_API_KEY_ID") == "" {
		t.Skip("Skipping e2e test: TELEMETRYFLOW_API_KEY_ID not set")
	}
	if os.Getenv("TELEMETRYFLOW_API_KEY_SECRET") == "" {
		t.Skip("Skipping e2e test: TELEMETRYFLOW_API_KEY_SECRET not set")
	}
}

func TestE2E_FullTelemetryPipeline(t *testing.T) {
	skipE2E(t)

	t.Run("should send all telemetry types successfully", func(t *testing.T) {
		// Create client from environment
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()

		// Initialize
		err = client.Initialize(ctx)
		require.NoError(t, err)

		defer func() {
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = client.Flush(shutdownCtx)
			_ = client.Shutdown(shutdownCtx)
		}()

		assert.True(t, client.IsInitialized())

		// Send metrics
		t.Log("Sending metrics...")
		err = client.IncrementCounter(ctx, "e2e.requests.total", 1, map[string]interface{}{
			"method": "GET",
			"status": 200,
		})
		assert.NoError(t, err)

		err = client.RecordGauge(ctx, "e2e.memory.usage", 512.5, map[string]interface{}{
			"unit": "MB",
		})
		assert.NoError(t, err)

		err = client.RecordHistogram(ctx, "e2e.request.duration", 0.125, "s", map[string]interface{}{
			"endpoint": "/api/test",
		})
		assert.NoError(t, err)

		// Send logs
		t.Log("Sending logs...")
		err = client.LogInfo(ctx, "E2E test started", map[string]interface{}{
			"test_id": "e2e-001",
		})
		assert.NoError(t, err)

		err = client.LogWarn(ctx, "E2E test warning example", map[string]interface{}{
			"warning_code": "W001",
		})
		assert.NoError(t, err)

		// Create traces
		t.Log("Creating traces...")
		spanID, err := client.StartSpan(ctx, "e2e.test.operation", "internal", map[string]interface{}{
			"test": "e2e",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, spanID)

		time.Sleep(100 * time.Millisecond)

		err = client.AddSpanEvent(ctx, spanID, "processing.step", map[string]interface{}{
			"step": 1,
		})
		assert.NoError(t, err)

		err = client.EndSpan(ctx, spanID, nil)
		assert.NoError(t, err)

		// Flush and verify no errors
		t.Log("Flushing telemetry...")
		err = client.Flush(ctx)
		assert.NoError(t, err)

		t.Log("E2E test completed successfully")
	})
}

func TestE2E_MetricsBatch(t *testing.T) {
	skipE2E(t)

	t.Run("should handle batch of metrics", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() {
			_ = client.Flush(ctx)
			_ = client.Shutdown(ctx)
		}()

		// Send batch of metrics
		for i := 0; i < 100; i++ {
			err = client.IncrementCounter(ctx, "e2e.batch.counter", 1, map[string]interface{}{
				"batch": "test",
				"index": i,
			})
			assert.NoError(t, err)
		}

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestE2E_ConcurrentTelemetry(t *testing.T) {
	skipE2E(t)

	t.Run("should handle concurrent telemetry from multiple goroutines", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() {
			_ = client.Flush(ctx)
			_ = client.Shutdown(ctx)
		}()

		var wg sync.WaitGroup
		errors := make(chan error, 300)

		// Launch goroutines for metrics
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				err := client.IncrementCounter(ctx, "e2e.concurrent.counter", 1, map[string]interface{}{
					"goroutine": id,
				})
				if err != nil {
					errors <- err
				}
			}(i)
		}

		// Launch goroutines for logs
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				err := client.LogInfo(ctx, fmt.Sprintf("Concurrent log %d", id), map[string]interface{}{
					"goroutine": id,
				})
				if err != nil {
					errors <- err
				}
			}(i)
		}

		// Launch goroutines for traces
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				spanID, err := client.StartSpan(ctx, "e2e.concurrent.span", "internal", map[string]interface{}{
					"goroutine": id,
				})
				if err != nil {
					errors <- err
					return
				}
				time.Sleep(10 * time.Millisecond)
				err = client.EndSpan(ctx, spanID, nil)
				if err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		var errCount int
		for err := range errors {
			t.Logf("Error: %v", err)
			errCount++
		}

		assert.Zero(t, errCount, "Expected no errors during concurrent operations")

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestE2E_LongRunningSession(t *testing.T) {
	skipE2E(t)
	if os.Getenv("TELEMETRYFLOW_E2E_LONG") != "true" {
		t.Skip("Skipping long-running test: set TELEMETRYFLOW_E2E_LONG=true to run")
	}

	t.Run("should maintain connection over extended period", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() {
			client.Flush(ctx)
			client.Shutdown(ctx)
		}()

		// Send telemetry every second for 30 seconds
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		timeout := time.After(30 * time.Second)
		count := 0

		for {
			select {
			case <-timeout:
				t.Logf("Sent %d telemetry items over 30 seconds", count)
				return
			case <-ticker.C:
				count++
				err := client.IncrementCounter(ctx, "e2e.long_running.counter", 1, map[string]interface{}{
					"tick": count,
				})
				assert.NoError(t, err)

				err = client.LogInfo(ctx, fmt.Sprintf("Long running test tick %d", count), nil)
				assert.NoError(t, err)
			}
		}
	})
}

func TestE2E_ErrorRecovery(t *testing.T) {
	skipE2E(t)

	t.Run("should continue after errors", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)
		defer func() {
			client.Flush(ctx)
			client.Shutdown(ctx)
		}()

		// Send some valid telemetry
		err = client.IncrementCounter(ctx, "e2e.recovery.before", 1, nil)
		assert.NoError(t, err)

		// Create a span and end it with an error
		spanID, err := client.StartSpan(ctx, "e2e.error.span", "internal", nil)
		assert.NoError(t, err)

		err = client.EndSpan(ctx, spanID, fmt.Errorf("simulated error"))
		assert.NoError(t, err)

		// Log an error
		err = client.LogError(ctx, "Simulated error for recovery test", map[string]interface{}{
			"error_code": "E001",
		})
		assert.NoError(t, err)

		// Continue sending valid telemetry after the error
		err = client.IncrementCounter(ctx, "e2e.recovery.after", 1, nil)
		assert.NoError(t, err)

		err = client.LogInfo(ctx, "Recovery successful", nil)
		assert.NoError(t, err)

		err = client.Flush(ctx)
		assert.NoError(t, err)
	})
}

func TestE2E_GracefulShutdown(t *testing.T) {
	skipE2E(t)

	t.Run("should flush pending data on shutdown", func(t *testing.T) {
		client, err := telemetryflow.NewFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		err = client.Initialize(ctx)
		require.NoError(t, err)

		// Send a batch of telemetry
		for i := 0; i < 50; i++ {
			_ = client.IncrementCounter(ctx, "e2e.shutdown.counter", 1, map[string]interface{}{
				"batch": i,
			})
		}

		// Shutdown with timeout - should flush pending data
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		start := time.Now()
		err = client.Flush(shutdownCtx)
		assert.NoError(t, err)

		err = client.Shutdown(shutdownCtx)
		assert.NoError(t, err)
		elapsed := time.Since(start)

		t.Logf("Shutdown completed in %v", elapsed)
		assert.False(t, client.IsInitialized())
	})
}

func TestE2E_MultipleServices(t *testing.T) {
	skipE2E(t)

	t.Run("should handle multiple service instances", func(t *testing.T) {
		// Create multiple clients with different service names
		services := []string{"service-a", "service-b", "service-c"}
		var clients []*telemetryflow.Client

		for _, svc := range services {
			client, err := telemetryflow.NewBuilder().
				WithAPIKeyFromEnv().
				WithEndpointFromEnv().
				WithService(svc, "1.0.0").
				WithEnvironmentFromEnv().
				Build()
			require.NoError(t, err)

			ctx := context.Background()
			err = client.Initialize(ctx)
			require.NoError(t, err)

			clients = append(clients, client)
		}

		// Send telemetry from each service
		ctx := context.Background()
		for i, client := range clients {
			serviceName := services[i]
			err := client.IncrementCounter(ctx, "e2e.multi_service.request", 1, map[string]interface{}{
				"service": serviceName,
			})
			assert.NoError(t, err)

			err = client.LogInfo(ctx, fmt.Sprintf("Request from %s", serviceName), nil)
			assert.NoError(t, err)
		}

		// Shutdown all clients
		for _, client := range clients {
			client.Flush(ctx)
			client.Shutdown(ctx)
		}
	})
}

// Benchmark tests
func BenchmarkE2E_MetricThroughput(b *testing.B) {
	if os.Getenv("TELEMETRYFLOW_E2E") != "true" {
		b.Skip("Skipping e2e benchmark: set TELEMETRYFLOW_E2E=true to run")
	}

	client, err := telemetryflow.NewFromEnv()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		b.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		client.Flush(ctx)
		client.Shutdown(ctx)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.IncrementCounter(ctx, "e2e.benchmark.counter", 1, nil)
	}
}

func BenchmarkE2E_LogThroughput(b *testing.B) {
	if os.Getenv("TELEMETRYFLOW_E2E") != "true" {
		b.Skip("Skipping e2e benchmark: set TELEMETRYFLOW_E2E=true to run")
	}

	client, err := telemetryflow.NewFromEnv()
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		b.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		client.Flush(ctx)
		client.Shutdown(ctx)
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.LogInfo(ctx, "Benchmark log message", nil)
	}
}
