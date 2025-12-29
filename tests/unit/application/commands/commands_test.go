// Package application_test provides unit tests for the application package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package application_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/application"
)

func TestRecordMetricCommand(t *testing.T) {
	t.Run("should create metric command with all fields", func(t *testing.T) {
		now := time.Now()
		cmd := &application.RecordMetricCommand{
			Name:  "http.request.duration",
			Value: 0.125,
			Unit:  "s",
			Attributes: map[string]interface{}{
				"method": "GET",
				"status": 200,
			},
			Timestamp: now,
		}

		assert.Equal(t, "http.request.duration", cmd.Name)
		assert.Equal(t, 0.125, cmd.Value)
		assert.Equal(t, "s", cmd.Unit)
		assert.Equal(t, "GET", cmd.Attributes["method"])
		assert.Equal(t, 200, cmd.Attributes["status"])
		assert.Equal(t, now, cmd.Timestamp)
	})

	t.Run("should allow nil attributes", func(t *testing.T) {
		cmd := &application.RecordMetricCommand{
			Name:       "simple.metric",
			Value:      42.0,
			Unit:       "count",
			Attributes: nil,
			Timestamp:  time.Now(),
		}

		assert.Nil(t, cmd.Attributes)
	})
}

func TestRecordCounterCommand(t *testing.T) {
	t.Run("should create counter command", func(t *testing.T) {
		cmd := &application.RecordCounterCommand{
			Name:  "requests.total",
			Value: 1,
			Attributes: map[string]interface{}{
				"endpoint": "/api/users",
			},
		}

		assert.Equal(t, "requests.total", cmd.Name)
		assert.Equal(t, int64(1), cmd.Value)
		assert.Equal(t, "/api/users", cmd.Attributes["endpoint"])
	})

	t.Run("should handle large counter values", func(t *testing.T) {
		cmd := &application.RecordCounterCommand{
			Name:  "bytes.transferred",
			Value: 9223372036854775807, // max int64
		}

		assert.Equal(t, int64(9223372036854775807), cmd.Value)
	})
}

func TestRecordGaugeCommand(t *testing.T) {
	t.Run("should create gauge command", func(t *testing.T) {
		cmd := &application.RecordGaugeCommand{
			Name:  "memory.usage",
			Value: 512.5,
			Attributes: map[string]interface{}{
				"unit": "MB",
			},
		}

		assert.Equal(t, "memory.usage", cmd.Name)
		assert.Equal(t, 512.5, cmd.Value)
	})

	t.Run("should handle negative values", func(t *testing.T) {
		cmd := &application.RecordGaugeCommand{
			Name:  "temperature",
			Value: -15.5,
		}

		assert.Equal(t, -15.5, cmd.Value)
	})
}

func TestRecordHistogramCommand(t *testing.T) {
	t.Run("should create histogram command", func(t *testing.T) {
		cmd := &application.RecordHistogramCommand{
			Name:  "request.size",
			Value: 1024.0,
			Unit:  "bytes",
			Attributes: map[string]interface{}{
				"content_type": "application/json",
			},
		}

		assert.Equal(t, "request.size", cmd.Name)
		assert.Equal(t, 1024.0, cmd.Value)
		assert.Equal(t, "bytes", cmd.Unit)
	})
}

func TestEmitLogCommand(t *testing.T) {
	t.Run("should create log command with all fields", func(t *testing.T) {
		now := time.Now()
		cmd := &application.EmitLogCommand{
			Severity: "error",
			Message:  "Database connection failed",
			Attributes: map[string]interface{}{
				"host":    "db.example.com",
				"port":    5432,
				"retries": 3,
			},
			Timestamp: now,
			TraceID:   "abc123",
			SpanID:    "def456",
		}

		assert.Equal(t, "error", cmd.Severity)
		assert.Equal(t, "Database connection failed", cmd.Message)
		assert.Equal(t, "db.example.com", cmd.Attributes["host"])
		assert.Equal(t, now, cmd.Timestamp)
		assert.Equal(t, "abc123", cmd.TraceID)
		assert.Equal(t, "def456", cmd.SpanID)
	})

	t.Run("should allow empty trace/span IDs", func(t *testing.T) {
		cmd := &application.EmitLogCommand{
			Severity: "info",
			Message:  "Application started",
		}

		assert.Empty(t, cmd.TraceID)
		assert.Empty(t, cmd.SpanID)
	})
}

func TestEmitBatchLogsCommand(t *testing.T) {
	t.Run("should create batch logs command", func(t *testing.T) {
		logs := []application.EmitLogCommand{
			{Severity: "info", Message: "Log 1"},
			{Severity: "warn", Message: "Log 2"},
			{Severity: "error", Message: "Log 3"},
		}

		cmd := &application.EmitBatchLogsCommand{
			Logs: logs,
		}

		assert.Len(t, cmd.Logs, 3)
		assert.Equal(t, "info", cmd.Logs[0].Severity)
		assert.Equal(t, "error", cmd.Logs[2].Severity)
	})

	t.Run("should handle empty batch", func(t *testing.T) {
		cmd := &application.EmitBatchLogsCommand{
			Logs: []application.EmitLogCommand{},
		}

		assert.Empty(t, cmd.Logs)
	})
}

func TestStartSpanCommand(t *testing.T) {
	t.Run("should create span command with all fields", func(t *testing.T) {
		cmd := &application.StartSpanCommand{
			Name: "http.request",
			Kind: "server",
			Attributes: map[string]interface{}{
				"http.method": "POST",
				"http.url":    "/api/orders",
			},
			ParentID: "parent-span-123",
		}

		assert.Equal(t, "http.request", cmd.Name)
		assert.Equal(t, "server", cmd.Kind)
		assert.Equal(t, "POST", cmd.Attributes["http.method"])
		assert.Equal(t, "parent-span-123", cmd.ParentID)
	})

	t.Run("should allow different span kinds", func(t *testing.T) {
		kinds := []string{"internal", "server", "client", "producer", "consumer"}

		for _, kind := range kinds {
			cmd := &application.StartSpanCommand{
				Name: "test.span",
				Kind: kind,
			}
			assert.Equal(t, kind, cmd.Kind)
		}
	})
}

func TestEndSpanCommand(t *testing.T) {
	t.Run("should create end span command without error", func(t *testing.T) {
		cmd := &application.EndSpanCommand{
			SpanID: "span-123",
			Error:  nil,
		}

		assert.Equal(t, "span-123", cmd.SpanID)
		assert.Nil(t, cmd.Error)
	})

	t.Run("should create end span command with error", func(t *testing.T) {
		err := assert.AnError
		cmd := &application.EndSpanCommand{
			SpanID: "span-456",
			Error:  err,
		}

		assert.Equal(t, "span-456", cmd.SpanID)
		assert.Equal(t, err, cmd.Error)
	})
}

func TestAddSpanEventCommand(t *testing.T) {
	t.Run("should create span event command", func(t *testing.T) {
		now := time.Now()
		cmd := &application.AddSpanEventCommand{
			SpanID: "span-789",
			Name:   "validation.complete",
			Attributes: map[string]interface{}{
				"valid": true,
			},
			Timestamp: now,
		}

		assert.Equal(t, "span-789", cmd.SpanID)
		assert.Equal(t, "validation.complete", cmd.Name)
		assert.True(t, cmd.Attributes["valid"].(bool))
		assert.Equal(t, now, cmd.Timestamp)
	})
}

func TestInitializeSDKCommand(t *testing.T) {
	t.Run("should create initialize command with config", func(t *testing.T) {
		// Note: We can't test with actual config here without importing domain
		// This test just verifies the struct works
		cmd := &application.InitializeSDKCommand{
			Config: nil, // Would be *domain.TelemetryConfig in real usage
		}

		assert.Nil(t, cmd.Config)
	})
}

func TestShutdownSDKCommand(t *testing.T) {
	t.Run("should create shutdown command with timeout", func(t *testing.T) {
		cmd := &application.ShutdownSDKCommand{
			Timeout: 30 * time.Second,
		}

		assert.Equal(t, 30*time.Second, cmd.Timeout)
	})
}

func TestFlushTelemetryCommand(t *testing.T) {
	t.Run("should create flush command with timeout", func(t *testing.T) {
		cmd := &application.FlushTelemetryCommand{
			Timeout: 10 * time.Second,
		}

		assert.Equal(t, 10*time.Second, cmd.Timeout)
	})
}

func TestCommandBus(t *testing.T) {
	t.Run("should create new command bus", func(t *testing.T) {
		bus := application.NewCommandBus()

		assert.NotNil(t, bus)
	})
}

// Benchmark tests
func BenchmarkRecordMetricCommand_Create(b *testing.B) {
	attrs := map[string]interface{}{
		"method": "GET",
		"status": 200,
	}
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = &application.RecordMetricCommand{
			Name:       "http.request.duration",
			Value:      0.125,
			Unit:       "s",
			Attributes: attrs,
			Timestamp:  now,
		}
	}
}
