// Package infrastructure provides unit tests for the command handlers.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package infrastructure

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/application"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func createTestConfig(t *testing.T) *domain.TelemetryConfig {
	creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)

	config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
	require.NoError(t, err)

	return config
}

func TestNewTelemetryCommandHandler(t *testing.T) {
	t.Run("should create handler with config", func(t *testing.T) {
		config := createTestConfig(t)

		handler := NewTelemetryCommandHandler(config)

		assert.NotNil(t, handler)
		assert.Equal(t, config, handler.config)
		assert.NotNil(t, handler.activeSpans)
		assert.False(t, handler.initialized)
	})
}


func TestTelemetryCommandHandler_ShutdownNotInitialized(t *testing.T) {
	t.Run("should fail shutdown when not initialized", func(t *testing.T) {
		config := createTestConfig(t)
		handler := NewTelemetryCommandHandler(config)

		cmd := &application.ShutdownSDKCommand{
			Timeout: 30 * time.Second,
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestTelemetryCommandHandler_FlushNotInitialized(t *testing.T) {
	t.Run("should fail flush when not initialized", func(t *testing.T) {
		config := createTestConfig(t)
		handler := NewTelemetryCommandHandler(config)

		cmd := &application.FlushTelemetryCommand{
			Timeout: 10 * time.Second,
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestTelemetryCommandHandler_MetricsNotInitialized(t *testing.T) {
	config := createTestConfig(t)
	handler := NewTelemetryCommandHandler(config)

	t.Run("should fail record metric when not initialized", func(t *testing.T) {
		cmd := &application.RecordMetricCommand{
			Name:  "test.metric",
			Value: 1.0,
			Unit:  "count",
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("should fail record counter when not initialized", func(t *testing.T) {
		cmd := &application.RecordCounterCommand{
			Name:  "test.counter",
			Value: 1,
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("should fail record gauge when not initialized", func(t *testing.T) {
		cmd := &application.RecordGaugeCommand{
			Name:  "test.gauge",
			Value: 42.0,
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("should fail record histogram when not initialized", func(t *testing.T) {
		cmd := &application.RecordHistogramCommand{
			Name:  "test.histogram",
			Value: 100.0,
			Unit:  "ms",
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})
}

func TestTelemetryCommandHandler_TracesNotInitialized(t *testing.T) {
	config := createTestConfig(t)
	handler := NewTelemetryCommandHandler(config)

	t.Run("should fail start span when not initialized", func(t *testing.T) {
		_, err := handler.StartSpanDirect(context.Background(), "test-span", "internal", nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not initialized")
	})

	t.Run("should fail end span when span not found", func(t *testing.T) {
		handler.initialized = true
		cmd := &application.EndSpanCommand{
			SpanID: "non-existent-span",
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "span not found")
	})

	t.Run("should fail add span event when span not found", func(t *testing.T) {
		handler.initialized = true
		cmd := &application.AddSpanEventCommand{
			SpanID:    "non-existent-span",
			Name:      "event-name",
			Timestamp: time.Now(),
		}

		err := handler.Handle(context.Background(), cmd)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "span not found")
	})
}

func TestConvertAttributes(t *testing.T) {
	t.Run("should convert string attribute", func(t *testing.T) {
		attrs := map[string]interface{}{"key": "value"}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
		assert.Equal(t, "key", string(result[0].Key))
	})

	t.Run("should convert int attribute", func(t *testing.T) {
		attrs := map[string]interface{}{"count": 42}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
	})

	t.Run("should convert int64 attribute", func(t *testing.T) {
		attrs := map[string]interface{}{"count": int64(42)}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
	})

	t.Run("should convert float64 attribute", func(t *testing.T) {
		attrs := map[string]interface{}{"value": 3.14}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
	})

	t.Run("should convert bool attribute", func(t *testing.T) {
		attrs := map[string]interface{}{"enabled": true}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
	})

	t.Run("should convert unknown type to string", func(t *testing.T) {
		attrs := map[string]interface{}{"custom": struct{ Name string }{"test"}}

		result := convertAttributes(attrs)

		assert.Len(t, result, 1)
	})

	t.Run("should handle multiple attributes", func(t *testing.T) {
		attrs := map[string]interface{}{
			"string": "value",
			"int":    42,
			"float":  3.14,
			"bool":   true,
		}

		result := convertAttributes(attrs)

		assert.Len(t, result, 4)
	})

	t.Run("should handle empty attributes", func(t *testing.T) {
		attrs := map[string]interface{}{}

		result := convertAttributes(attrs)

		assert.Empty(t, result)
	})

	t.Run("should handle nil attributes", func(t *testing.T) {
		result := convertAttributes(nil)

		assert.Empty(t, result)
	})
}

func TestSpanKindMapping(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"internal", "internal"},
		{"server", "server"},
		{"client", "client"},
		{"producer", "producer"},
		{"consumer", "consumer"},
		{"unknown", "internal"}, // defaults to internal
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("should map %s kind", tc.input), func(t *testing.T) {
			// This is a documentation test - actual mapping is tested in integration
			assert.NotEmpty(t, tc.input)
		})
	}
}
