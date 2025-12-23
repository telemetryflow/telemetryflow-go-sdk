package traces

import (
	"context"

	"{{.ModulePath}}/telemetry"
)

// StartSpan starts a new trace span
func StartSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	return telemetry.Client().StartSpan(ctx, name, "internal", attrs)
}

// StartServerSpan starts a server-side span
func StartServerSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	return telemetry.Client().StartSpan(ctx, name, "server", attrs)
}

// StartClientSpan starts a client-side span
func StartClientSpan(ctx context.Context, name string, attrs map[string]interface{}) (string, error) {
	return telemetry.Client().StartSpan(ctx, name, "client", attrs)
}

// EndSpan ends an active span
func EndSpan(ctx context.Context, spanID string, err error) error {
	return telemetry.Client().EndSpan(ctx, spanID, err)
}

// AddEvent adds an event to an active span
func AddEvent(ctx context.Context, spanID, eventName string, attrs map[string]interface{}) error {
	return telemetry.Client().AddSpanEvent(ctx, spanID, eventName, attrs)
}

// Trace is a helper that wraps a function with tracing
func Trace(ctx context.Context, name string, fn func() error) error {
	spanID, err := StartSpan(ctx, name, nil)
	if err != nil {
		return fn()
	}
	defer EndSpan(ctx, spanID, nil)

	if err := fn(); err != nil {
		EndSpan(ctx, spanID, err)
		return err
	}
	return nil
}
