// Package main demonstrates TelemetryFlow Go SDK v1.1.2 basic usage.
//
// This example shows:
// - Simple initialization from environment variables
// - Builder pattern with TFO v2 API support
// - Collector identity configuration (aligned with tfoidentityextension)
// - TFO v2-only mode for TFO-Collector v1.1.2 compatibility
//
// Run with:
//
//	export TELEMETRYFLOW_API_KEY_ID=tfk_your_key
//	export TELEMETRYFLOW_API_KEY_SECRET=tfs_your_secret
//	export TELEMETRYFLOW_SERVICE_NAME=my-service
//	go run main.go
package main

import (
	"context"
	"log"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

func main() {
	// Example 1: Simple initialization from environment variables
	simpleExample()

	// Example 2: Builder pattern with TFO v2 API configuration
	builderExample()

	// Example 3: TFO v2-only mode (aligned with TFO-Collector v1.1.2)
	v2OnlyExample()

	// Example 4: Complete example with all signal types
	completeExample()
}

func simpleExample() {
	log.Println("=== Simple Example (Auto-configuration) ===")

	// Create client from environment variables
	// Requires: TELEMETRYFLOW_API_KEY_ID, TELEMETRYFLOW_API_KEY_SECRET,
	//           TELEMETRYFLOW_SERVICE_NAME in environment
	// Automatically reads TFO v2 settings and collector identity from env
	client, err := telemetryflow.NewFromEnv()
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Initialize the SDK
	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown client: %v", err)
		}
	}()

	// Send a simple metric
	if err := client.IncrementCounter(ctx, "requests.total", 1, map[string]interface{}{
		"method": "GET",
		"status": 200,
	}); err != nil {
		log.Printf("Failed to send metric: %v", err)
	}

	// Flush and wait
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
	log.Println("Simple example completed")
}

func builderExample() {
	log.Println("=== Builder Example (TFO v2 API + Collector Identity) ===")

	// Build client with fluent API including TFO v2 features
	// aligned with tfoexporter and tfoidentityextension
	client := telemetryflow.NewBuilder().
		WithAPIKey("tfk_your_key_id", "tfs_your_key_secret").
		WithEndpoint("api.telemetryflow.id:4317").
		WithService("my-go-service", "1.0.0").
		WithServiceNamespace("telemetryflow").
		WithEnvironment("production").
		WithGRPC().
		// TFO v2 API settings (aligned with tfoexporter)
		WithV2API(true).
		// Collector Identity (aligned with tfoidentityextension)
		WithCollectorName("My Go Service SDK").
		WithCollectorDescription("TelemetryFlow Go SDK embedded collector").
		WithDatacenter("aws-us-east-1").
		WithCollectorTag("team", "backend").
		WithCollectorTag("environment", "production").
		WithEnrichResources(true).
		// Enable exemplars for metrics-to-traces correlation
		WithExemplars(true).
		// Custom resource attributes
		WithCustomAttribute("cloud.provider", "aws").
		WithCustomAttribute("cloud.region", "us-east-1").
		MustBuild()

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown client: %v", err)
		}
	}()

	// Send metrics
	if err := client.RecordGauge(ctx, "cpu.usage", 75.5, map[string]interface{}{
		"host": "server-01",
	}); err != nil {
		log.Printf("Failed to record gauge: %v", err)
	}

	if err := client.RecordHistogram(ctx, "request.duration", 0.25, "s", map[string]interface{}{
		"endpoint": "/api/users",
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
	log.Println("Builder example completed")
}

func v2OnlyExample() {
	log.Println("=== TFO v2-Only Mode Example ===")

	// Build client in v2-only mode for TFO-Collector v1.1.2 (OCB-native)
	// This mode uses only v2 endpoints: /v2/traces, /v2/metrics, /v2/logs
	client := telemetryflow.NewBuilder().
		WithAPIKey("tfk_your_key_id", "tfs_your_key_secret").
		WithEndpoint("api.telemetryflow.id:4317").
		WithService("v2-only-service", "1.1.2").
		WithServiceNamespace("telemetryflow").
		WithEnvironment("production").
		WithGRPC().
		// Enable v2-only mode - only uses v2 endpoints
		WithV2Only().
		// Collector Identity is required for v2-only mode
		WithCollectorName("TFO v2 Service").
		WithCollectorDescription("TelemetryFlow Go SDK - TFO v2-only mode").
		WithDatacenter("default").
		WithCollectorTag("mode", "v2-only").
		WithCollectorTag("sdk_version", "1.1.2").
		WithEnrichResources(true).
		WithExemplars(true).
		MustBuild()

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown client: %v", err)
		}
	}()

	// Log that we're in v2-only mode
	if err := client.LogInfo(ctx, "SDK initialized in v2-only mode", map[string]interface{}{
		"sdk_version":           "1.1.2",
		"tfo_collector_version": "1.1.2",
		"mode":                  "v2-only",
	}); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	// Send a metric through v2 endpoint
	if err := client.IncrementCounter(ctx, "v2.requests.total", 1, map[string]interface{}{
		"api_version": "v2",
	}); err != nil {
		log.Printf("Failed to send metric: %v", err)
	}

	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
	log.Println("TFO v2-only example completed")
}

func completeExample() {
	log.Println("=== Complete Example (All Signals with TFO v2) ===")

	// Create client with full TFO v1.1.2 configuration
	client := telemetryflow.NewBuilder().
		WithAPIKey("tfk_your_key_id", "tfs_your_key_secret").
		WithEndpoint("api.telemetryflow.id:4317").
		WithService("complete-service", "2.0.0").
		WithServiceNamespace("telemetryflow").
		WithEnvironment("staging").
		WithSignals(true, true, true). // Enable all signals
		// TFO v2 API enabled by default
		WithV2API(true).
		// Collector Identity
		WithCollectorName("Complete Service SDK").
		WithDatacenter("aws-us-east-1").
		WithEnrichResources(true).
		// Enable exemplars for metrics-to-traces correlation
		WithExemplars(true).
		MustBuild()

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer func() {
		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown client: %v", err)
		}
	}()

	// === METRICS ===
	log.Println("Sending metrics...")

	// Counter
	if err := client.IncrementCounter(ctx, "api.requests", 1, map[string]interface{}{
		"method": "POST",
		"path":   "/api/orders",
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	// Gauge
	if err := client.RecordGauge(ctx, "memory.usage", 512.0, map[string]interface{}{
		"unit": "MB",
	}); err != nil {
		log.Printf("Failed to record gauge: %v", err)
	}

	// Histogram
	if err := client.RecordHistogram(ctx, "db.query.duration", 0.125, "s", map[string]interface{}{
		"query_type": "SELECT",
		"table":      "users",
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	// === LOGS ===
	log.Println("Sending logs...")

	if err := client.LogInfo(ctx, "Application started successfully", map[string]interface{}{
		"version":     "2.0.0",
		"port":        8080,
		"sdk_version": "1.1.2",
	}); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	if err := client.LogWarn(ctx, "High memory usage detected", map[string]interface{}{
		"memory_mb": 512,
		"threshold": 400,
	}); err != nil {
		log.Printf("Failed to log warning: %v", err)
	}

	if err := client.LogError(ctx, "Failed to connect to database", map[string]interface{}{
		"error":    "connection timeout",
		"host":     "db.example.com",
		"attempts": 3,
	}); err != nil {
		log.Printf("Failed to log error: %v", err)
	}

	// === TRACES ===
	log.Println("Creating traces...")

	// Start a span
	spanID, err := client.StartSpan(ctx, "process-order", "internal", map[string]interface{}{
		"order_id":    "12345",
		"customer_id": "67890",
	})
	if err != nil {
		log.Printf("Failed to start span: %v", err)
	}

	// Simulate some work
	time.Sleep(100 * time.Millisecond)

	// Add event to span
	if err := client.AddSpanEvent(ctx, spanID, "validation.complete", map[string]interface{}{
		"valid": true,
	}); err != nil {
		log.Printf("Failed to add span event: %v", err)
	}

	time.Sleep(50 * time.Millisecond)

	// End span
	if err := client.EndSpan(ctx, spanID, nil); err != nil {
		log.Printf("Failed to end span: %v", err)
	}

	// === FLUSH & SHUTDOWN ===
	log.Println("Flushing telemetry...")
	if err := client.Flush(ctx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}

	log.Println("Complete example finished")
}
