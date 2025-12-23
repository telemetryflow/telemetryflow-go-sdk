package logs

import (
	"context"

	"{{.ModulePath}}/telemetry"
)

// Info logs an info-level message
func Info(message string, attrs map[string]interface{}) {
	ctx := context.Background()
	telemetry.Client().LogInfo(ctx, message, attrs)
}

// Warn logs a warning-level message
func Warn(message string, attrs map[string]interface{}) {
	ctx := context.Background()
	telemetry.Client().LogWarn(ctx, message, attrs)
}

// Error logs an error-level message
func Error(message string, attrs map[string]interface{}) {
	ctx := context.Background()
	telemetry.Client().LogError(ctx, message, attrs)
}

// Debug logs a debug-level message
func Debug(message string, attrs map[string]interface{}) {
	ctx := context.Background()
	telemetry.Client().Log(ctx, "debug", message, attrs)
}

// WithContext logs with a specific context (for trace correlation)
func WithContext(ctx context.Context, severity, message string, attrs map[string]interface{}) {
	telemetry.Client().Log(ctx, severity, message, attrs)
}
