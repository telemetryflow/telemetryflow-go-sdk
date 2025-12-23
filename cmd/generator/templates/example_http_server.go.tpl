package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	// Setup routes with telemetry middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", withTelemetry(handleHome))
	mux.HandleFunc("/api/users", withTelemetry(handleUsers))
	mux.HandleFunc("/health", handleHealth)

	server := &http.Server{
		Addr:         ":{{.Port}}",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	logs.Info("Server starting", map[string]interface{}{
		"port": "{{.Port}}",
	})

	log.Printf("Server starting on :{{.Port}}")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// withTelemetry wraps an HTTP handler with telemetry instrumentation
func withTelemetry(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Start a trace span for this request
		spanID, err := traces.StartSpan(ctx, fmt.Sprintf("HTTP %s %s", r.Method, r.URL.Path), map[string]interface{}{
			"http.method":      r.Method,
			"http.url":         r.URL.String(),
			"http.user_agent":  r.UserAgent(),
			"http.remote_addr": r.RemoteAddr,
		})
		if err == nil {
			defer func() {
				traces.EndSpan(ctx, spanID, nil)
			}()
		}

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the actual handler
		next(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		metrics.RecordHTTPRequest(r.Method, r.URL.Path, wrapped.statusCode, duration)

		// Log the request
		logs.Info("HTTP request completed", map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"status":     wrapped.statusCode,
			"duration_s": duration,
			"user_agent": r.UserAgent(),
		})
	}
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

func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to {{.ProjectName}}!"))
}

func handleUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Create a child span for database operation
	spanID, _ := traces.StartSpan(ctx, "db.query.users", map[string]interface{}{
		"db.operation": "SELECT",
		"db.table":     "users",
	})

	// Simulate database query
	start := time.Now()
	time.Sleep(50 * time.Millisecond)
	queryDuration := time.Since(start).Seconds()

	// Record database metrics
	metrics.RecordHistogram("db.query.duration", queryDuration, "s", map[string]interface{}{
		"operation": "SELECT",
		"table":     "users",
	})

	traces.EndSpan(ctx, spanID, nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"users": []}`))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
}
