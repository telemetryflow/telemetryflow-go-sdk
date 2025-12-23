// Package main demonstrates TelemetryFlow SDK integration with an HTTP server.
//
// This example shows:
// - HTTP middleware for automatic request tracing
// - Request duration histograms
// - Error counting and logging
// - Graceful shutdown with telemetry flush
//
// Run with:
//
//	export TELEMETRYFLOW_API_KEY_ID=tfk_your_key
//	export TELEMETRYFLOW_API_KEY_SECRET=tfs_your_secret
//	go run main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

func main() {
	// Initialize TelemetryFlow client
	var err error
	client, err = telemetryflow.NewBuilder().
		WithAPIKeyFromEnv().
		WithEndpointFromEnv().
		WithService("http-server-example", "1.0.0").
		WithEnvironmentFromEnv().
		WithGRPC().
		WithCustomAttribute("example", "http-server").
		Build()

	if err != nil {
		log.Fatalf("Failed to create TelemetryFlow client: %v", err)
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize TelemetryFlow: %v", err)
	}

	// Create HTTP server with telemetry middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/api/users", handleUsers)
	mux.HandleFunc("/api/orders", handleOrders)
	mux.HandleFunc("/health", handleHealth)

	// Wrap with telemetry middleware
	handler := TelemetryMiddleware(mux)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Println("Starting HTTP server on :8080")
		client.LogInfo(ctx, "HTTP server started", map[string]interface{}{
			"port":    8080,
			"version": "1.0.0",
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	client.LogInfo(ctx, "Server shutdown initiated", nil)

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Flush and shutdown telemetry
	client.Flush(shutdownCtx)
	client.Shutdown(shutdownCtx)

	log.Println("Server stopped")
}

// TelemetryMiddleware wraps an http.Handler with telemetry instrumentation
func TelemetryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Start span for this request
		spanID, err := client.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path), "server", map[string]interface{}{
			"http.method":      r.Method,
			"http.url":         r.URL.String(),
			"http.path":        r.URL.Path,
			"http.user_agent":  r.UserAgent(),
			"http.remote_addr": r.RemoteAddr,
		})
		if err != nil {
			log.Printf("Failed to start span: %v", err)
		}

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Record request metrics
		client.IncrementCounter(ctx, "http.requests.total", 1, map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": wrapped.statusCode,
		})

		client.RecordHistogram(ctx, "http.request.duration", duration.Seconds(), "s", map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.Path,
			"status": wrapped.statusCode,
		})

		// Log errors
		if wrapped.statusCode >= 400 {
			client.IncrementCounter(ctx, "http.errors.total", 1, map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"status": wrapped.statusCode,
			})

			if wrapped.statusCode >= 500 {
				client.LogError(ctx, "HTTP request failed", map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"status":      wrapped.statusCode,
					"duration_ms": duration.Milliseconds(),
				})
			}
		}

		// End span
		if spanID != "" {
			var spanErr error
			if wrapped.statusCode >= 500 {
				spanErr = fmt.Errorf("HTTP %d", wrapped.statusCode)
			}
			client.EndSpan(ctx, spanID, spanErr)
		}
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Handler functions

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	response := map[string]string{
		"message": "Welcome to TelemetryFlow HTTP Server Example",
		"version": "1.0.0",
	}
	writeJSON(w, http.StatusOK, response)
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		// Simulate database query
		spanID, _ := client.StartSpan(ctx, "db.query.users", "client", map[string]interface{}{
			"db.operation": "SELECT",
			"db.table":     "users",
		})

		time.Sleep(50 * time.Millisecond) // Simulate query time

		client.AddSpanEvent(ctx, spanID, "query.executed", map[string]interface{}{
			"rows_returned": 10,
		})
		client.EndSpan(ctx, spanID, nil)

		users := []map[string]interface{}{
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
			{"id": 3, "name": "Charlie", "email": "charlie@example.com"},
		}
		writeJSON(w, http.StatusOK, users)

	case http.MethodPost:
		client.LogInfo(ctx, "Creating new user", map[string]interface{}{
			"source": "api",
		})
		response := map[string]interface{}{
			"id":      4,
			"message": "User created successfully",
		}
		writeJSON(w, http.StatusCreated, response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	switch r.Method {
	case http.MethodGet:
		orders := []map[string]interface{}{
			{"id": "ord_001", "total": 99.99, "status": "completed"},
			{"id": "ord_002", "total": 149.50, "status": "pending"},
		}
		writeJSON(w, http.StatusOK, orders)

	case http.MethodPost:
		// Start order processing span
		spanID, _ := client.StartSpan(ctx, "order.process", "internal", map[string]interface{}{
			"order.source": "api",
		})

		// Simulate order validation
		client.AddSpanEvent(ctx, spanID, "order.validated", nil)
		time.Sleep(30 * time.Millisecond)

		// Simulate payment processing
		paymentSpanID, _ := client.StartSpan(ctx, "payment.process", "client", map[string]interface{}{
			"payment.provider": "stripe",
		})
		time.Sleep(100 * time.Millisecond)
		client.EndSpan(ctx, paymentSpanID, nil)

		client.AddSpanEvent(ctx, spanID, "payment.completed", map[string]interface{}{
			"payment.method": "credit_card",
		})

		// Record business metric
		client.IncrementCounter(ctx, "orders.created", 1, map[string]interface{}{
			"source":         "api",
			"payment_method": "credit_card",
		})

		client.RecordHistogram(ctx, "order.value", 199.99, "usd", map[string]interface{}{
			"source": "api",
		})

		client.EndSpan(ctx, spanID, nil)

		response := map[string]interface{}{
			"id":      "ord_003",
			"total":   199.99,
			"status":  "pending",
			"message": "Order created successfully",
		}
		writeJSON(w, http.StatusCreated, response)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	// Record health check
	client.RecordGauge(r.Context(), "health.check", 1, map[string]interface{}{
		"endpoint": "/health",
	})

	response := map[string]string{
		"status":  "healthy",
		"version": "1.0.0",
	}
	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
