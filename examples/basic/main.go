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

	// Example 2: Builder pattern with custom configuration
	builderExample()

	// Example 3: Complete example with all signal types
	completeExample()
}

func simpleExample() {
	log.Println("=== Simple Example ===")

	// Create client from environment variables
	// Requires: TELEMETRYFLOW_API_KEY_ID, TELEMETRYFLOW_API_KEY_SECRET,
	//           TELEMETRYFLOW_SERVICE_NAME in environment
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
	log.Println("=== Builder Example ===")

	// Build client with fluent API
	client := telemetryflow.NewBuilder().
		WithAPIKey("tfk_your_key_id", "tfs_your_key_secret").
		WithEndpoint("api.telemetryflow.id:4317").
		WithService("my-go-service", "1.0.0").
		WithEnvironment("production").
		WithGRPC().
		WithCustomAttribute("team", "backend").
		WithCustomAttribute("region", "us-east-1").
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

func completeExample() {
	log.Println("=== Complete Example ===")

	// Create client
	client := telemetryflow.NewBuilder().
		WithAPIKey("tfk_your_key_id", "tfs_your_key_secret").
		WithEndpoint("api.telemetryflow.id:4317").
		WithService("complete-service", "2.0.0").
		WithEnvironment("staging").
		WithSignals(true, true, true). // Enable all signals
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
		"version": "2.0.0",
		"port":    8080,
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
