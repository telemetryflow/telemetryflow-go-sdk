// Package telemetryflow provides unit tests for the client.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package telemetryflow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func createTestConfig(t *testing.T) *domain.TelemetryConfig {
	creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)

	config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
	require.NoError(t, err)

	return config
}

func TestNewClient(t *testing.T) {
	t.Run("should create client with valid config", func(t *testing.T) {
		config := createTestConfig(t)

		client, err := NewClient(config)

		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.False(t, client.IsInitialized())
	})

	t.Run("should fail with nil config", func(t *testing.T) {
		_, err := NewClient(nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})
}

func TestClient_IsInitialized(t *testing.T) {
	t.Run("should return false for new client", func(t *testing.T) {
		config := createTestConfig(t)
		client, err := NewClient(config)
		require.NoError(t, err)

		assert.False(t, client.IsInitialized())
	})
}

func TestClient_Config(t *testing.T) {
	t.Run("should return config", func(t *testing.T) {
		config := createTestConfig(t)
		client, err := NewClient(config)
		require.NoError(t, err)

		assert.Equal(t, config, client.Config())
		assert.Equal(t, "test-service", client.Config().ServiceName())
	})
}

func TestClient_NotInitializedErrors(t *testing.T) {
	config := createTestConfig(t)
	client, err := NewClient(config)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("Flush should fail when not initialized", func(t *testing.T) {
		err := client.Flush(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("RecordMetric should fail when not initialized", func(t *testing.T) {
		err := client.RecordMetric(ctx, "test.metric", 1.0, "count", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("IncrementCounter should fail when not initialized", func(t *testing.T) {
		err := client.IncrementCounter(ctx, "test.counter", 1, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("RecordGauge should fail when not initialized", func(t *testing.T) {
		err := client.RecordGauge(ctx, "test.gauge", 42.0, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("RecordHistogram should fail when not initialized", func(t *testing.T) {
		err := client.RecordHistogram(ctx, "test.histogram", 100.0, "ms", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("Log should fail when not initialized", func(t *testing.T) {
		err := client.Log(ctx, "info", "test message", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogInfo should fail when not initialized", func(t *testing.T) {
		err := client.LogInfo(ctx, "test message", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogWarn should fail when not initialized", func(t *testing.T) {
		err := client.LogWarn(ctx, "test message", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("LogError should fail when not initialized", func(t *testing.T) {
		err := client.LogError(ctx, "test message", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("StartSpan should fail when not initialized", func(t *testing.T) {
		_, err := client.StartSpan(ctx, "test-span", "internal", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("EndSpan should fail when not initialized", func(t *testing.T) {
		err := client.EndSpan(ctx, "span-id", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("AddSpanEvent should fail when not initialized", func(t *testing.T) {
		err := client.AddSpanEvent(ctx, "span-id", "event-name", nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestClient_Shutdown(t *testing.T) {
	t.Run("should succeed when not initialized", func(t *testing.T) {
		config := createTestConfig(t)
		client, err := NewClient(config)
		require.NoError(t, err)

		err = client.Shutdown(context.Background())
		assert.NoError(t, err)
	})
}
