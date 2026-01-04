// Package instrumentation provides auto-instrumentation for TelemetryFlow SDK.
package instrumentation

import (
	"context"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware provides auto-instrumentation for net/http handlers
type HTTPMiddleware struct {
	tracer       trace.Tracer
	config       *Config
	metrics      *HTTPRequestMetrics
	spanFormatter SpanNameFormatter
}

// NewHTTPMiddleware creates a new HTTP middleware with the given options
func NewHTTPMiddleware(opts ...Option) *HTTPMiddleware {
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
			metrics = nil // Fail silently, continue without metrics
		}
	}

	return &HTTPMiddleware{
		tracer:        tracer,
		config:        cfg,
		metrics:       metrics,
		spanFormatter: DefaultSpanNameFormatter,
	}
}

// WithSpanNameFormatter sets a custom span name formatter
func (m *HTTPMiddleware) WithSpanNameFormatter(fn SpanNameFormatter) *HTTPMiddleware {
	m.spanFormatter = fn
	return m
}

// Handler wraps an http.Handler with tracing and metrics
func (m *HTTPMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply filter if configured
		if m.config.FilterFunc != nil && m.config.FilterFunc(r) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Extract trace context from incoming request
		ctx := ExtractTraceParent(r.Context(), r.Header)

		// Determine span name
		spanName := m.spanFormatter("", r)

		// Create span attributes
		attrs := []attribute.KeyValue{
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.URLPath(r.URL.Path),
			semconv.URLScheme(r.URL.Scheme),
			semconv.ServerAddress(r.Host),
			semconv.UserAgentOriginal(r.UserAgent()),
		}

		if r.URL.RawQuery != "" {
			attrs = append(attrs, semconv.URLQuery(r.URL.RawQuery))
		}

		// Add service attributes if configured
		if m.config.ServiceName != "" {
			attrs = append(attrs, semconv.ServiceName(m.config.ServiceName))
		}
		if m.config.ServiceVersion != "" {
			attrs = append(attrs, semconv.ServiceVersion(m.config.ServiceVersion))
		}

		// Start span
		ctx, span := m.tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Track active requests
		if m.metrics != nil {
			m.metrics.ActiveRequests.Add(ctx, 1)
			defer m.metrics.ActiveRequests.Add(ctx, -1)
		}

		// Create wrapped response writer to capture status code and size
		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Update request context
		r = r.WithContext(ctx)

		// Call next handler
		next.ServeHTTP(rw, r)

		// Record span status
		statusCode := rw.statusCode
		if statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Add response attributes
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(statusCode))

		// Record metrics
		duration := time.Since(start)
		if m.metrics != nil {
			m.metrics.RecordRequest(ctx, r.Method, r.URL.Path, statusCode, duration, r.ContentLength, int64(rw.size))
		}
	})
}

// HandlerFunc wraps an http.HandlerFunc with tracing and metrics
func (m *HTTPMiddleware) HandlerFunc(next http.HandlerFunc) http.HandlerFunc {
	return m.Handler(next).ServeHTTP
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Unwrap returns the original ResponseWriter (for http.Hijacker, etc.)
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

// HTTPClient provides auto-instrumentation for http.Client
type HTTPClient struct {
	client  *http.Client
	tracer  trace.Tracer
	config  *Config
	metrics *HTTPRequestMetrics
}

// NewHTTPClient creates a new instrumented HTTP client
func NewHTTPClient(client *http.Client, opts ...Option) *HTTPClient {
	if client == nil {
		client = http.DefaultClient
	}

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

	return &HTTPClient{
		client:  client,
		tracer:  tracer,
		config:  cfg,
		metrics: metrics,
	}
}

// Do executes an HTTP request with tracing and metrics
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.do(req.Context(), req)
}

func (c *HTTPClient) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	start := time.Now()

	// Create span name
	spanName := req.Method + " " + req.URL.Host + req.URL.Path

	// Create span attributes
	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.URLFull(req.URL.String()),
		semconv.ServerAddress(req.URL.Host),
	}

	// Start span
	ctx, span := c.tracer.Start(ctx, spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	// Inject trace context into outgoing request
	InjectTraceParent(ctx, req.Header)

	// Update request context
	req = req.WithContext(ctx)

	// Execute request
	resp, err := c.client.Do(req)

	duration := time.Since(start)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	// Add response attributes
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(resp.StatusCode))

	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	// Record metrics
	if c.metrics != nil {
		c.metrics.RecordRequest(ctx, req.Method, req.URL.Path, resp.StatusCode, duration, req.ContentLength, resp.ContentLength)
	}

	return resp, nil
}

// Get performs an HTTP GET request with tracing
func (c *HTTPClient) Get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(ctx, req)
}

// Post performs an HTTP POST request with tracing
func (c *HTTPClient) Post(ctx context.Context, url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.do(ctx, req)
}

// WrapRoundTripper creates an instrumented RoundTripper
func WrapRoundTripper(rt http.RoundTripper, opts ...Option) http.RoundTripper {
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

	if rt == nil {
		rt = http.DefaultTransport
	}

	return &instrumentedRoundTripper{
		base:   rt,
		tracer: tracer,
		config: cfg,
	}
}

type instrumentedRoundTripper struct {
	base   http.RoundTripper
	tracer trace.Tracer
	config *Config
}

func (t *instrumentedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	spanName := req.Method + " " + req.URL.Host + req.URL.Path

	attrs := []attribute.KeyValue{
		semconv.HTTPRequestMethodKey.String(req.Method),
		semconv.URLFull(req.URL.String()),
		semconv.ServerAddress(req.URL.Host),
	}

	ctx, span := t.tracer.Start(req.Context(), spanName,
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attrs...),
	)
	defer span.End()

	// Inject trace context
	InjectTraceParent(ctx, req.Header)
	req = req.WithContext(ctx)

	resp, err := t.base.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(resp.StatusCode))
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return resp, nil
}
