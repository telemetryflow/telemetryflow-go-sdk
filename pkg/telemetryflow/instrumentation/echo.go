// Package instrumentation provides auto-instrumentation for TelemetryFlow SDK.
// This file provides Echo framework instrumentation.
package instrumentation

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// EchoContext is an interface that matches echo.Context
// This allows the instrumentation to work without importing Echo directly
type EchoContext interface {
	Request() interface {
		Header() interface{ Get(string) string }
	}
	Response() interface {
		Status() int
		Size() int64
	}
	Path() string
	RealIP() string
	SetRequest(interface{})
	Next() error
}

// EchoMiddlewareConfig holds Echo middleware configuration
type EchoMiddlewareConfig struct {
	*Config
	// Skipper defines a function to skip middleware
	Skipper func(c interface{}) bool
	// SpanNameFormatter formats the span name
	SpanNameFormatter func(c interface{}) string
}

// DefaultEchoMiddlewareConfig returns default Echo middleware configuration
func DefaultEchoMiddlewareConfig() *EchoMiddlewareConfig {
	return &EchoMiddlewareConfig{
		Config: DefaultConfig(),
		Skipper: func(c interface{}) bool {
			return false
		},
	}
}

// EchoMiddlewareFunc creates an Echo middleware function for tracing and metrics.
// This returns a generic middleware function that wraps the Echo handler.
//
// Usage with Echo:
//
//	import (
//	    "github.com/labstack/echo/v4"
//	    "github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/instrumentation"
//	)
//
//	e := echo.New()
//	e.Use(instrumentation.EchoMiddlewareFunc(
//	    instrumentation.WithServiceInfo("my-service", "1.0.0"),
//	))
func EchoMiddlewareFunc(opts ...Option) func(next interface{}) interface{} {
	cfg := ApplyOptions(opts...)

	tracer := cfg.TracerProvider.Tracer(InstrumentationName,
		trace.WithInstrumentationVersion(InstrumentationVersion),
	)

	var metrics *HTTPRequestMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewHTTPRequestMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	// Return a generic middleware wrapper
	// The actual Echo middleware should be created by the user
	return func(next interface{}) interface{} {
		return &echoMiddlewareHandler{
			next:    next,
			tracer:  tracer,
			config:  cfg,
			metrics: metrics,
		}
	}
}

type echoMiddlewareHandler struct {
	next    interface{}
	tracer  trace.Tracer
	config  *Config
	metrics *HTTPRequestMetrics
}

// EchoInstrumentationHelper provides helper methods for Echo instrumentation
type EchoInstrumentationHelper struct {
	tracer  trace.Tracer
	config  *Config
	metrics *HTTPRequestMetrics
}

// NewEchoInstrumentationHelper creates a new Echo instrumentation helper
func NewEchoInstrumentationHelper(opts ...Option) *EchoInstrumentationHelper {
	cfg := ApplyOptions(opts...)

	var tracer trace.Tracer
	if cfg.TracerProvider != nil {
		tracer = cfg.TracerProvider.Tracer(InstrumentationName,
			trace.WithInstrumentationVersion(InstrumentationVersion),
		)
	} else {
		tracer = otel.Tracer(InstrumentationName,
			trace.WithInstrumentationVersion(InstrumentationVersion),
		)
	}

	var metrics *HTTPRequestMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewHTTPRequestMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return &EchoInstrumentationHelper{
		tracer:  tracer,
		config:  cfg,
		metrics: metrics,
	}
}

// GetTracer returns the tracer
func (h *EchoInstrumentationHelper) GetTracer() trace.Tracer {
	return h.tracer
}

// GetConfig returns the configuration
func (h *EchoInstrumentationHelper) GetConfig() *Config {
	return h.config
}

// GetMetrics returns the HTTP metrics recorder
func (h *EchoInstrumentationHelper) GetMetrics() *HTTPRequestMetrics {
	return h.metrics
}

// CreateSpanName creates a span name from method and path
func (h *EchoInstrumentationHelper) CreateSpanName(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}

// CreateSpanAttributes creates span attributes for an HTTP request
func (h *EchoInstrumentationHelper) CreateSpanAttributes(method, path, host, userAgent, clientIP string) []attribute.KeyValue {
	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(method),
		semconv.URLPath(path),
		semconv.ServerAddress(host),
	}

	if userAgent != "" {
		attrs = append(attrs, semconv.UserAgentOriginal(userAgent))
	}
	if clientIP != "" {
		attrs = append(attrs, semconv.ClientAddress(clientIP))
	}

	if h.config.ServiceName != "" {
		attrs = append(attrs, semconv.ServiceName(h.config.ServiceName))
	}
	if h.config.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(h.config.ServiceVersion))
	}

	return attrs
}

// SetSpanStatus sets the span status based on HTTP status code
func (h *EchoInstrumentationHelper) SetSpanStatus(span trace.Span, statusCode int) {
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(statusCode))
	if statusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", statusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}
}

// RecordMetrics records HTTP request metrics
func (h *EchoInstrumentationHelper) RecordMetrics(method, path string, statusCode int, duration time.Duration, requestSize, responseSize int64) {
	if h.metrics != nil {
		h.metrics.RecordRequest(context.Background(), method, path, statusCode, duration, requestSize, responseSize)
	}
}

// TraceEchoRequest is a helper function to trace an Echo request
// This function should be called at the beginning of an Echo handler
//
// Usage:
//
//	func handler(c echo.Context) error {
//	    ctx, span := helper.TraceEchoRequest(c.Request().Context(), c.Request().Method, c.Path(), c.Request().Host)
//	    defer span.End()
//	    c.SetRequest(c.Request().WithContext(ctx))
//	    // ... handler logic
//	}
func (h *EchoInstrumentationHelper) TraceEchoRequest(parentCtx context.Context, method, path, host string) (context.Context, trace.Span) {
	spanName := h.CreateSpanName(method, path)
	attrs := h.CreateSpanAttributes(method, path, host, "", "")

	return h.tracer.Start(parentCtx, spanName,
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(attrs...),
	)
}
