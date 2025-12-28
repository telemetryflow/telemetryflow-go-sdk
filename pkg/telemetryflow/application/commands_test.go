// Package application provides unit tests for the application layer commands.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommandBus(t *testing.T) {
	t.Run("should create new command bus", func(t *testing.T) {
		bus := NewCommandBus()

		assert.NotNil(t, bus)
		assert.NotNil(t, bus.handlers)
	})

	t.Run("should register handler", func(t *testing.T) {
		bus := NewCommandBus()
		handler := &mockCommandHandler{}

		bus.Register("test", handler)

		assert.Equal(t, handler, bus.handlers["test"])
	})
}

type mockCommandHandler struct{}

func (h *mockCommandHandler) Handle(ctx context.Context, cmd Command) error {
	return nil
}

func TestInitializeSDKCommand(t *testing.T) {
	t.Run("should create initialize command", func(t *testing.T) {
		cmd := &InitializeSDKCommand{
			Config: nil, // would be a real config in production
		}

		assert.NotNil(t, cmd)
	})
}

func TestShutdownSDKCommand(t *testing.T) {
	t.Run("should create shutdown command with timeout", func(t *testing.T) {
		cmd := &ShutdownSDKCommand{
			Timeout: 30 * time.Second,
		}

		assert.Equal(t, 30*time.Second, cmd.Timeout)
	})
}

func TestFlushTelemetryCommand(t *testing.T) {
	t.Run("should create flush command with timeout", func(t *testing.T) {
		cmd := &FlushTelemetryCommand{
			Timeout: 10 * time.Second,
		}

		assert.Equal(t, 10*time.Second, cmd.Timeout)
	})
}

func TestRecordMetricCommand(t *testing.T) {
	t.Run("should create record metric command", func(t *testing.T) {
		now := time.Now()
		attrs := map[string]interface{}{"key": "value"}

		cmd := &RecordMetricCommand{
			Name:       "test.metric",
			Value:      42.5,
			Unit:       "count",
			Attributes: attrs,
			Timestamp:  now,
		}

		assert.Equal(t, "test.metric", cmd.Name)
		assert.Equal(t, 42.5, cmd.Value)
		assert.Equal(t, "count", cmd.Unit)
		assert.Equal(t, attrs, cmd.Attributes)
		assert.Equal(t, now, cmd.Timestamp)
	})
}

func TestRecordCounterCommand(t *testing.T) {
	t.Run("should create record counter command", func(t *testing.T) {
		attrs := map[string]interface{}{"key": "value"}

		cmd := &RecordCounterCommand{
			Name:       "test.counter",
			Value:      10,
			Attributes: attrs,
		}

		assert.Equal(t, "test.counter", cmd.Name)
		assert.Equal(t, int64(10), cmd.Value)
		assert.Equal(t, attrs, cmd.Attributes)
	})
}

func TestRecordGaugeCommand(t *testing.T) {
	t.Run("should create record gauge command", func(t *testing.T) {
		attrs := map[string]interface{}{"key": "value"}

		cmd := &RecordGaugeCommand{
			Name:       "test.gauge",
			Value:      100.5,
			Attributes: attrs,
		}

		assert.Equal(t, "test.gauge", cmd.Name)
		assert.Equal(t, 100.5, cmd.Value)
		assert.Equal(t, attrs, cmd.Attributes)
	})
}

func TestRecordHistogramCommand(t *testing.T) {
	t.Run("should create record histogram command", func(t *testing.T) {
		attrs := map[string]interface{}{"key": "value"}

		cmd := &RecordHistogramCommand{
			Name:       "test.histogram",
			Value:      250.0,
			Unit:       "ms",
			Attributes: attrs,
		}

		assert.Equal(t, "test.histogram", cmd.Name)
		assert.Equal(t, 250.0, cmd.Value)
		assert.Equal(t, "ms", cmd.Unit)
		assert.Equal(t, attrs, cmd.Attributes)
	})
}

func TestEmitLogCommand(t *testing.T) {
	t.Run("should create emit log command", func(t *testing.T) {
		now := time.Now()
		attrs := map[string]interface{}{"key": "value"}

		cmd := &EmitLogCommand{
			Severity:   "info",
			Message:    "test message",
			Attributes: attrs,
			Timestamp:  now,
		}

		assert.Equal(t, "info", cmd.Severity)
		assert.Equal(t, "test message", cmd.Message)
		assert.Equal(t, attrs, cmd.Attributes)
		assert.Equal(t, now, cmd.Timestamp)
	})
}

func TestStartSpanCommand(t *testing.T) {
	t.Run("should create start span command", func(t *testing.T) {
		attrs := map[string]interface{}{"key": "value"}

		cmd := &StartSpanCommand{
			Name:       "test-span",
			Kind:       "server",
			Attributes: attrs,
		}

		assert.Equal(t, "test-span", cmd.Name)
		assert.Equal(t, "server", cmd.Kind)
		assert.Equal(t, attrs, cmd.Attributes)
	})
}

func TestEndSpanCommand(t *testing.T) {
	t.Run("should create end span command", func(t *testing.T) {
		cmd := &EndSpanCommand{
			SpanID: "span-123",
			Error:  nil,
		}

		assert.Equal(t, "span-123", cmd.SpanID)
		assert.Nil(t, cmd.Error)
	})

	t.Run("should create end span command with error", func(t *testing.T) {
		err := assert.AnError
		cmd := &EndSpanCommand{
			SpanID: "span-123",
			Error:  err,
		}

		assert.Equal(t, "span-123", cmd.SpanID)
		assert.Equal(t, err, cmd.Error)
	})
}

func TestAddSpanEventCommand(t *testing.T) {
	t.Run("should create add span event command", func(t *testing.T) {
		now := time.Now()
		attrs := map[string]interface{}{"key": "value"}

		cmd := &AddSpanEventCommand{
			SpanID:     "span-123",
			Name:       "event-name",
			Attributes: attrs,
			Timestamp:  now,
		}

		assert.Equal(t, "span-123", cmd.SpanID)
		assert.Equal(t, "event-name", cmd.Name)
		assert.Equal(t, attrs, cmd.Attributes)
		assert.Equal(t, now, cmd.Timestamp)
	})
}
