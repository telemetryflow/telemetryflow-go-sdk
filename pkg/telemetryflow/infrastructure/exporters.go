package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// OTLPExporterFactory creates OTLP exporters based on configuration
type OTLPExporterFactory struct {
	config *domain.TelemetryConfig
}

// NewOTLPExporterFactory creates a new exporter factory
func NewOTLPExporterFactory(config *domain.TelemetryConfig) *OTLPExporterFactory {
	return &OTLPExporterFactory{
		config: config,
	}
}

// CreateResource creates an OTLP resource with service information
func (f *OTLPExporterFactory) CreateResource(ctx context.Context) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{
		semconv.ServiceName(f.config.ServiceName()),
		semconv.ServiceVersion(f.config.ServiceVersion()),
		semconv.DeploymentEnvironment(f.config.Environment()),
	}

	// Add custom attributes
	for key, value := range f.config.CustomAttributes() {
		attrs = append(attrs, attribute.String(key, value))
	}

	return resource.New(
		ctx,
		resource.WithAttributes(attrs...),
		resource.WithProcessRuntimeDescription(),
		resource.WithTelemetrySDK(),
	)
}

// CreateTraceExporter creates a trace exporter based on protocol
func (f *OTLPExporterFactory) CreateTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	if !f.config.IsSignalEnabled(domain.SignalTraces) {
		return nil, fmt.Errorf("traces signal is not enabled")
	}

	switch f.config.Protocol() {
	case domain.ProtocolGRPC:
		return f.createGRPCTraceExporter(ctx)
	case domain.ProtocolHTTP:
		return f.createHTTPTraceExporter(ctx)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", f.config.Protocol())
	}
}

// CreateMetricExporter creates a metric exporter based on protocol
func (f *OTLPExporterFactory) CreateMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	if !f.config.IsSignalEnabled(domain.SignalMetrics) {
		return nil, fmt.Errorf("metrics signal is not enabled")
	}

	switch f.config.Protocol() {
	case domain.ProtocolGRPC:
		return f.createGRPCMetricExporter(ctx)
	case domain.ProtocolHTTP:
		return f.createHTTPMetricExporter(ctx)
	default:
		return nil, fmt.Errorf("unsupported protocol: %s", f.config.Protocol())
	}
}

// ===== GRPC EXPORTERS =====

func (f *OTLPExporterFactory) createGRPCTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(f.config.Endpoint()),
		otlptracegrpc.WithTimeout(f.config.Timeout()),
		otlptracegrpc.WithHeaders(f.getAuthHeaders()),
		otlptracegrpc.WithDialOption(grpc.WithUnaryInterceptor(f.authInterceptor())),
	}

	if f.config.IsInsecure() {
		opts = append(opts, otlptracegrpc.WithInsecure())
	} else {
		opts = append(opts, otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	if f.config.IsCompressionEnabled() {
		opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	}

	if f.config.IsRetryEnabled() {
		opts = append(opts, otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: f.config.RetryBackoff(),
			MaxInterval:     f.config.RetryBackoff() * 2,
			MaxElapsedTime:  time.Duration(f.config.MaxRetries()) * f.config.RetryBackoff(),
		}))
	}

	return otlptracegrpc.New(ctx, opts...)
}

func (f *OTLPExporterFactory) createGRPCMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(f.config.Endpoint()),
		otlpmetricgrpc.WithTimeout(f.config.Timeout()),
		otlpmetricgrpc.WithHeaders(f.getAuthHeaders()),
		otlpmetricgrpc.WithDialOption(grpc.WithUnaryInterceptor(f.authInterceptor())),
	}

	if f.config.IsInsecure() {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	} else {
		opts = append(opts, otlpmetricgrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, "")))
	}

	if f.config.IsCompressionEnabled() {
		opts = append(opts, otlpmetricgrpc.WithCompressor("gzip"))
	}

	if f.config.IsRetryEnabled() {
		opts = append(opts, otlpmetricgrpc.WithRetry(otlpmetricgrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: f.config.RetryBackoff(),
			MaxInterval:     f.config.RetryBackoff() * 2,
			MaxElapsedTime:  time.Duration(f.config.MaxRetries()) * f.config.RetryBackoff(),
		}))
	}

	return otlpmetricgrpc.New(ctx, opts...)
}

// ===== HTTP EXPORTERS =====

func (f *OTLPExporterFactory) createHTTPTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(f.config.Endpoint()),
		otlptracehttp.WithTimeout(f.config.Timeout()),
		otlptracehttp.WithHeaders(f.getAuthHeaders()),
	}

	if f.config.IsInsecure() {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	if f.config.IsCompressionEnabled() {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	if f.config.IsRetryEnabled() {
		opts = append(opts, otlptracehttp.WithRetry(otlptracehttp.RetryConfig{
			Enabled:         true,
			InitialInterval: f.config.RetryBackoff(),
			MaxInterval:     f.config.RetryBackoff() * 2,
			MaxElapsedTime:  time.Duration(f.config.MaxRetries()) * f.config.RetryBackoff(),
		}))
	}

	return otlptracehttp.New(ctx, opts...)
}

func (f *OTLPExporterFactory) createHTTPMetricExporter(ctx context.Context) (sdkmetric.Exporter, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(f.config.Endpoint()),
		otlpmetrichttp.WithTimeout(f.config.Timeout()),
		otlpmetrichttp.WithHeaders(f.getAuthHeaders()),
	}

	if f.config.IsInsecure() {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}

	if f.config.IsCompressionEnabled() {
		opts = append(opts, otlpmetrichttp.WithCompression(otlpmetrichttp.GzipCompression))
	}

	if f.config.IsRetryEnabled() {
		opts = append(opts, otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{
			Enabled:         true,
			InitialInterval: f.config.RetryBackoff(),
			MaxInterval:     f.config.RetryBackoff() * 2,
			MaxElapsedTime:  time.Duration(f.config.MaxRetries()) * f.config.RetryBackoff(),
		}))
	}

	return otlpmetrichttp.New(ctx, opts...)
}

// ===== HELPER METHODS =====

// getAuthHeaders returns headers with TelemetryFlow authentication
func (f *OTLPExporterFactory) getAuthHeaders() map[string]string {
	return map[string]string{
		"authorization": f.config.Credentials().AuthorizationHeader(),
		"content-type":  "application/x-protobuf",
	}
}

// authInterceptor creates a gRPC interceptor that adds authentication
func (f *OTLPExporterFactory) authInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Add authentication to context
		ctx = metadata.AppendToOutgoingContext(ctx,
			"authorization", f.config.Credentials().AuthorizationHeader(),
		)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
