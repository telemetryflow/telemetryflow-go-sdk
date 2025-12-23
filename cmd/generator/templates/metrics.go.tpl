package metrics

import (
	"context"

	"{{.ModulePath}}/telemetry"
)

// IncrementCounter increments a counter metric
func IncrementCounter(name string, value int64, labels map[string]interface{}) {
	ctx := context.Background()
	if err := telemetry.Client().IncrementCounter(ctx, name, value, labels); err != nil {
		// Handle error (log, etc.)
	}
}

// RecordGauge records a gauge metric
func RecordGauge(name string, value float64, labels map[string]interface{}) {
	ctx := context.Background()
	if err := telemetry.Client().RecordGauge(ctx, name, value, labels); err != nil {
		// Handle error
	}
}

// RecordHistogram records a histogram measurement
func RecordHistogram(name string, value float64, unit string, labels map[string]interface{}) {
	ctx := context.Background()
	if err := telemetry.Client().RecordHistogram(ctx, name, value, unit, labels); err != nil {
		// Handle error
	}
}

// Common application metrics

// RecordHTTPRequest records metrics for an HTTP request
func RecordHTTPRequest(method, path string, statusCode int, duration float64) {
	RecordHistogram("http.request.duration", duration, "s", map[string]interface{}{
		"method": method,
		"path":   path,
		"status": statusCode,
	})
	IncrementCounter("http.requests.total", 1, map[string]interface{}{
		"method": method,
		"path":   path,
		"status": statusCode,
	})
}

// RecordDatabaseQuery records metrics for a database query
func RecordDatabaseQuery(operation, table string, duration float64, success bool) {
	RecordHistogram("db.query.duration", duration, "s", map[string]interface{}{
		"operation": operation,
		"table":     table,
		"success":   success,
	})
}
