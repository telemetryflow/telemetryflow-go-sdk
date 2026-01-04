// Package instrumentation provides auto-instrumentation helpers for TelemetryFlow SDK.
// This package offers convenient wrappers for popular Go frameworks and libraries
// to automatically capture traces, metrics, and logs.
package instrumentation

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// InstrumentationName is the name used for OpenTelemetry instrumentation
	InstrumentationName = "github.com/telemetryflow/telemetryflow-go-sdk/instrumentation"
	// InstrumentationVersion is the version of the instrumentation
	InstrumentationVersion = "1.1.2"
)

// Config holds common instrumentation configuration
type Config struct {
	// TracerProvider is the tracer provider to use
	TracerProvider trace.TracerProvider
	// MeterProvider is the meter provider to use
	MeterProvider metric.MeterProvider
	// Propagators are the propagators to use for context propagation
	Propagators propagation.TextMapPropagator
	// ServiceName is the name of the instrumented service
	ServiceName string
	// ServiceVersion is the version of the instrumented service
	ServiceVersion string
	// EnableMetrics enables automatic metrics collection
	EnableMetrics bool
	// EnableTracing enables automatic tracing
	EnableTracing bool
	// RecordRequestBody records request body in spans (use with caution for PII)
	RecordRequestBody bool
	// RecordResponseBody records response body in spans (use with caution for PII)
	RecordResponseBody bool
	// FilterFunc is a function to filter requests from instrumentation
	FilterFunc func(*http.Request) bool
}

// DefaultConfig returns the default instrumentation configuration
func DefaultConfig() *Config {
	return &Config{
		TracerProvider: otel.GetTracerProvider(),
		MeterProvider:  otel.GetMeterProvider(),
		Propagators:    otel.GetTextMapPropagator(),
		EnableMetrics:  true,
		EnableTracing:  true,
	}
}

// Option is a function that configures the instrumentation
type Option func(*Config)

// WithTracerProvider sets the tracer provider
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(c *Config) {
		c.TracerProvider = tp
	}
}

// WithMeterProvider sets the meter provider
func WithMeterProvider(mp metric.MeterProvider) Option {
	return func(c *Config) {
		c.MeterProvider = mp
	}
}

// WithPropagators sets the propagators
func WithPropagators(p propagation.TextMapPropagator) Option {
	return func(c *Config) {
		c.Propagators = p
	}
}

// WithServiceInfo sets the service name and version
func WithServiceInfo(name, version string) Option {
	return func(c *Config) {
		c.ServiceName = name
		c.ServiceVersion = version
	}
}

// WithMetrics enables or disables metrics collection
func WithMetrics(enabled bool) Option {
	return func(c *Config) {
		c.EnableMetrics = enabled
	}
}

// WithTracing enables or disables tracing
func WithTracing(enabled bool) Option {
	return func(c *Config) {
		c.EnableTracing = enabled
	}
}

// WithFilter sets a filter function for requests
func WithFilter(fn func(*http.Request) bool) Option {
	return func(c *Config) {
		c.FilterFunc = fn
	}
}

// ApplyOptions applies configuration options
func ApplyOptions(opts ...Option) *Config {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// HTTPRequestMetrics holds metrics for HTTP requests
type HTTPRequestMetrics struct {
	RequestCounter    metric.Int64Counter
	RequestDuration   metric.Float64Histogram
	RequestBodySize   metric.Int64Histogram
	ResponseBodySize  metric.Int64Histogram
	ActiveRequests    metric.Int64UpDownCounter
}

// NewHTTPRequestMetrics creates HTTP request metrics
func NewHTTPRequestMetrics(mp metric.MeterProvider) (*HTTPRequestMetrics, error) {
	meter := mp.Meter(InstrumentationName,
		metric.WithInstrumentationVersion(InstrumentationVersion),
	)

	requestCounter, err := meter.Int64Counter("http.server.request.count",
		metric.WithDescription("Number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram("http.server.request.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	requestBodySize, err := meter.Int64Histogram("http.server.request.body.size",
		metric.WithDescription("HTTP request body size in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	responseBodySize, err := meter.Int64Histogram("http.server.response.body.size",
		metric.WithDescription("HTTP response body size in bytes"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return nil, err
	}

	activeRequests, err := meter.Int64UpDownCounter("http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	return &HTTPRequestMetrics{
		RequestCounter:   requestCounter,
		RequestDuration:  requestDuration,
		RequestBodySize:  requestBodySize,
		ResponseBodySize: responseBodySize,
		ActiveRequests:   activeRequests,
	}, nil
}

// RecordRequest records HTTP request metrics
func (m *HTTPRequestMetrics) RecordRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(method),
		semconv.HTTPRouteKey.String(path),
		semconv.HTTPResponseStatusCodeKey.Int(statusCode),
	}

	m.RequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.RequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))

	if requestSize > 0 {
		m.RequestBodySize.Record(ctx, requestSize, metric.WithAttributes(attrs...))
	}
	if responseSize > 0 {
		m.ResponseBodySize.Record(ctx, responseSize, metric.WithAttributes(attrs...))
	}
}

// DatabaseMetrics holds metrics for database operations
type DatabaseMetrics struct {
	QueryCounter     metric.Int64Counter
	QueryDuration    metric.Float64Histogram
	ConnectionsOpen  metric.Int64UpDownCounter
	ConnectionsIdle  metric.Int64UpDownCounter
	ConnectionErrors metric.Int64Counter
}

// NewDatabaseMetrics creates database metrics
func NewDatabaseMetrics(mp metric.MeterProvider) (*DatabaseMetrics, error) {
	meter := mp.Meter(InstrumentationName,
		metric.WithInstrumentationVersion(InstrumentationVersion),
	)

	queryCounter, err := meter.Int64Counter("db.client.operation.count",
		metric.WithDescription("Number of database operations"),
		metric.WithUnit("{operation}"),
	)
	if err != nil {
		return nil, err
	}

	queryDuration, err := meter.Float64Histogram("db.client.operation.duration",
		metric.WithDescription("Database operation duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	connectionsOpen, err := meter.Int64UpDownCounter("db.client.connection.count",
		metric.WithDescription("Number of open database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	connectionsIdle, err := meter.Int64UpDownCounter("db.client.connection.idle",
		metric.WithDescription("Number of idle database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	connectionErrors, err := meter.Int64Counter("db.client.connection.error",
		metric.WithDescription("Number of database connection errors"),
		metric.WithUnit("{error}"),
	)
	if err != nil {
		return nil, err
	}

	return &DatabaseMetrics{
		QueryCounter:     queryCounter,
		QueryDuration:    queryDuration,
		ConnectionsOpen:  connectionsOpen,
		ConnectionsIdle:  connectionsIdle,
		ConnectionErrors: connectionErrors,
	}, nil
}

// RecordQuery records database query metrics
func (m *DatabaseMetrics) RecordQuery(ctx context.Context, operation, table string, duration time.Duration, hasError bool) {
	attrs := []attribute.KeyValue{
		attribute.String("db.operation", operation),
		attribute.String("db.sql.table", table),
		attribute.Bool("error", hasError),
	}

	m.QueryCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.QueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// SpanNameFormatter is a function type for formatting span names
type SpanNameFormatter func(operation string, r *http.Request) string

// DefaultSpanNameFormatter returns the default span name formatter
func DefaultSpanNameFormatter(operation string, r *http.Request) string {
	if r != nil {
		return r.Method + " " + r.URL.Path
	}
	return operation
}

// ExtractTraceParent extracts the trace parent from HTTP headers
func ExtractTraceParent(ctx context.Context, headers http.Header) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(headers))
}

// InjectTraceParent injects the trace parent into HTTP headers
func InjectTraceParent(ctx context.Context, headers http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(headers))
}
