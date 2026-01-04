// Package instrumentation provides auto-instrumentation for TelemetryFlow SDK.
// This file provides gRPC instrumentation for both client and server.
package instrumentation

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCConfig holds gRPC instrumentation configuration
type GRPCConfig struct {
	*Config
	// RecordMessageEvents records message events for streaming calls
	RecordMessageEvents bool
}

// DefaultGRPCConfig returns default gRPC configuration
func DefaultGRPCConfig() *GRPCConfig {
	return &GRPCConfig{
		Config:              DefaultConfig(),
		RecordMessageEvents: true,
	}
}

// GRPCOption is a function that configures gRPC instrumentation
type GRPCOption func(*GRPCConfig)

// WithRecordMessageEvents enables/disables message event recording
func WithRecordMessageEvents(enabled bool) GRPCOption {
	return func(c *GRPCConfig) {
		c.RecordMessageEvents = enabled
	}
}

// GRPCMetrics holds metrics for gRPC operations
type GRPCMetrics struct {
	ServerRequestCounter  metric.Int64Counter
	ServerRequestDuration metric.Float64Histogram
	ClientRequestCounter  metric.Int64Counter
	ClientRequestDuration metric.Float64Histogram
}

// NewGRPCMetrics creates gRPC metrics
func NewGRPCMetrics(mp metric.MeterProvider) (*GRPCMetrics, error) {
	meter := mp.Meter(InstrumentationName,
		metric.WithInstrumentationVersion(InstrumentationVersion),
	)

	serverRequestCounter, err := meter.Int64Counter("rpc.server.request.count",
		metric.WithDescription("Number of gRPC server requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	serverRequestDuration, err := meter.Float64Histogram("rpc.server.request.duration",
		metric.WithDescription("gRPC server request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	clientRequestCounter, err := meter.Int64Counter("rpc.client.request.count",
		metric.WithDescription("Number of gRPC client requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	clientRequestDuration, err := meter.Float64Histogram("rpc.client.request.duration",
		metric.WithDescription("gRPC client request duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	return &GRPCMetrics{
		ServerRequestCounter:  serverRequestCounter,
		ServerRequestDuration: serverRequestDuration,
		ClientRequestCounter:  clientRequestCounter,
		ClientRequestDuration: clientRequestDuration,
	}, nil
}

// UnaryServerInterceptor returns a gRPC unary server interceptor with tracing and metrics
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
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

	var metrics *GRPCMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewGRPCMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()

		// Extract trace context from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = otel.GetTextMapPropagator().Extract(ctx, metadataCarrier(md))
		}

		// Create span
		spanName := info.FullMethod
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCMethod(info.FullMethod),
		}

		if cfg.ServiceName != "" {
			attrs = append(attrs, semconv.ServiceName(cfg.ServiceName))
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Call handler
		resp, err := handler(ctx, req)

		duration := time.Since(start)
		statusCode := grpcCodes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		// Set gRPC status code
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)))

		// Record metrics
		if metrics != nil {
			metricAttrs := []attribute.KeyValue{
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
				semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)),
			}
			metrics.ServerRequestCounter.Add(ctx, 1, metric.WithAttributes(metricAttrs...))
			metrics.ServerRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttrs...))
		}

		return resp, err
	}
}

// StreamServerInterceptor returns a gRPC stream server interceptor with tracing and metrics
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
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

	var metrics *GRPCMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewGRPCMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		ctx := ss.Context()

		// Extract trace context from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = otel.GetTextMapPropagator().Extract(ctx, metadataCarrier(md))
		}

		// Create span
		spanName := info.FullMethod
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCMethod(info.FullMethod),
		}

		if cfg.ServiceName != "" {
			attrs = append(attrs, semconv.ServiceName(cfg.ServiceName))
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Wrap the stream with context
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		// Call handler
		err := handler(srv, wrappedStream)

		duration := time.Since(start)
		statusCode := grpcCodes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)))

		// Record metrics
		if metrics != nil {
			metricAttrs := []attribute.KeyValue{
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(info.FullMethod),
				semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)),
			}
			metrics.ServerRequestCounter.Add(ctx, 1, metric.WithAttributes(metricAttrs...))
			metrics.ServerRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttrs...))
		}

		return err
	}
}

// UnaryClientInterceptor returns a gRPC unary client interceptor with tracing and metrics
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
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

	var metrics *GRPCMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewGRPCMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
		start := time.Now()

		// Create span
		spanName := method
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCMethod(method),
		}

		if cfg.ServiceName != "" {
			attrs = append(attrs, semconv.ServiceName(cfg.ServiceName))
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attrs...),
		)
		defer span.End()

		// Inject trace context into metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(md))
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Call invoker
		err := invoker(ctx, method, req, reply, cc, callOpts...)

		duration := time.Since(start)
		statusCode := grpcCodes.OK
		if err != nil {
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)))

		// Record metrics
		if metrics != nil {
			metricAttrs := []attribute.KeyValue{
				semconv.RPCSystemGRPC,
				semconv.RPCMethod(method),
				semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)),
			}
			metrics.ClientRequestCounter.Add(ctx, 1, metric.WithAttributes(metricAttrs...))
			metrics.ClientRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttrs...))
		}

		return err
	}
}

// StreamClientInterceptor returns a gRPC stream client interceptor with tracing and metrics
func StreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
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

	var metrics *GRPCMetrics
	if cfg.EnableMetrics && cfg.MeterProvider != nil {
		var err error
		metrics, err = NewGRPCMetrics(cfg.MeterProvider)
		if err != nil {
			metrics = nil
		}
	}

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, callOpts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()

		// Create span
		spanName := method
		attrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCMethod(method),
		}

		if cfg.ServiceName != "" {
			attrs = append(attrs, semconv.ServiceName(cfg.ServiceName))
		}

		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindClient),
			trace.WithAttributes(attrs...),
		)

		// Inject trace context into metadata
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}
		otel.GetTextMapPropagator().Inject(ctx, metadataCarrier(md))
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Call streamer
		stream, err := streamer(ctx, desc, cc, method, callOpts...)

		if err != nil {
			duration := time.Since(start)
			statusCode := grpcCodes.Unknown
			if s, ok := status.FromError(err); ok {
				statusCode = s.Code()
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)))
			span.End()

			// Record metrics
			if metrics != nil {
				metricAttrs := []attribute.KeyValue{
					semconv.RPCSystemGRPC,
					semconv.RPCMethod(method),
					semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)),
				}
				metrics.ClientRequestCounter.Add(ctx, 1, metric.WithAttributes(metricAttrs...))
				metrics.ClientRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttrs...))
			}

			return nil, err
		}

		return &wrappedClientStream{
			ClientStream: stream,
			span:         span,
			method:       method,
			metrics:      metrics,
			startTime:    start,
		}, nil
	}
}

// metadataCarrier implements propagation.TextMapCarrier for gRPC metadata
type metadataCarrier metadata.MD

func (mc metadataCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (mc metadataCarrier) Set(key, value string) {
	metadata.MD(mc).Set(key, value)
}

func (mc metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range mc {
		keys = append(keys, k)
	}
	return keys
}

// Ensure interface compliance
var _ propagation.TextMapCarrier = metadataCarrier{}

// wrappedServerStream wraps grpc.ServerStream with context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// wrappedClientStream wraps grpc.ClientStream with tracing
type wrappedClientStream struct {
	grpc.ClientStream
	span      trace.Span
	method    string
	metrics   *GRPCMetrics
	startTime time.Time
	finished  bool
}

func (w *wrappedClientStream) CloseSend() error {
	err := w.ClientStream.CloseSend()
	if err != nil {
		w.span.RecordError(err)
	}
	return err
}

func (w *wrappedClientStream) RecvMsg(m interface{}) error {
	err := w.ClientStream.RecvMsg(m)
	if err != nil {
		w.finish(err)
	}
	return err
}

func (w *wrappedClientStream) finish(err error) {
	if w.finished {
		return
	}
	w.finished = true

	duration := time.Since(w.startTime)
	statusCode := grpcCodes.OK
	if err != nil {
		if s, ok := status.FromError(err); ok {
			statusCode = s.Code()
		}
		w.span.RecordError(err)
		w.span.SetStatus(codes.Error, err.Error())
	} else {
		w.span.SetStatus(codes.Ok, "")
	}

	w.span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)))
	w.span.End()

	// Record metrics
	if w.metrics != nil {
		ctx := context.Background()
		metricAttrs := []attribute.KeyValue{
			semconv.RPCSystemGRPC,
			semconv.RPCMethod(w.method),
			semconv.RPCGRPCStatusCodeKey.Int(int(statusCode)),
		}
		w.metrics.ClientRequestCounter.Add(ctx, 1, metric.WithAttributes(metricAttrs...))
		w.metrics.ClientRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(metricAttrs...))
	}
}

// GRPCServerDialOptions returns gRPC dial options for server instrumentation
func GRPCServerDialOptions(opts ...Option) []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.UnaryInterceptor(UnaryServerInterceptor(opts...)),
		grpc.StreamInterceptor(StreamServerInterceptor(opts...)),
	}
}

// GRPCClientDialOptions returns gRPC dial options for client instrumentation
func GRPCClientDialOptions(opts ...Option) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(opts...)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(opts...)),
	}
}

// GRPCServerChainInterceptor returns chained server interceptors
func GRPCServerChainInterceptor(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		chain := handler
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			next := chain
			chain = func(ctx context.Context, req interface{}) (interface{}, error) {
				return interceptor(ctx, req, info, func(ctx context.Context, req interface{}) (interface{}, error) {
					return next(ctx, req)
				})
			}
		}
		return chain(ctx, req)
	}
}
