package main

import (
	"context"
	"log"
	"time"

	"{{.ModulePath}}/telemetry"
	"{{.ModulePath}}/telemetry/logs"
	"{{.ModulePath}}/telemetry/metrics"
	"{{.ModulePath}}/telemetry/traces"
)

func main() {
	// Initialize TelemetryFlow SDK
	if err := telemetry.Init(); err != nil {
		log.Fatal(err)
	}
	defer telemetry.Shutdown()

	ctx := context.Background()

	// Example: Record metrics
	metrics.IncrementCounter("app.started", 1, map[string]interface{}{
		"version": "1.0.0",
	})

	// Example: Log a message
	logs.Info("Application started", map[string]interface{}{
		"environment": "{{.Environment}}",
	})

	// Example: Create a trace span
	spanID, err := traces.StartSpan(ctx, "main.process", map[string]interface{}{
		"operation": "startup",
	})
	if err != nil {
		log.Printf("Failed to start span: %v", err)
	}

	// Simulate work
	time.Sleep(100 * time.Millisecond)

	// End the span
	if err := traces.EndSpan(ctx, spanID, nil); err != nil {
		log.Printf("Failed to end span: %v", err)
	}

	// Record a gauge metric
	metrics.RecordGauge("app.ready", 1, nil)

	log.Println("Basic TelemetryFlow example completed")
}
