package infrastructure

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/application"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// TelemetryCommandHandler handles all telemetry-related commands
type TelemetryCommandHandler struct {
	config         *domain.TelemetryConfig
	tracerProvider *sdktrace.TracerProvider
	meterProvider  *sdkmetric.MeterProvider
	tracer         trace.Tracer
	meter          otelmetric.Meter
	activeSpans    map[string]trace.Span
	spansMutex     sync.RWMutex
	initialized    bool
	initMutex      sync.Mutex
}

// NewTelemetryCommandHandler creates a new command handler
func NewTelemetryCommandHandler(config *domain.TelemetryConfig) *TelemetryCommandHandler {
	return &TelemetryCommandHandler{
		config:      config,
		activeSpans: make(map[string]trace.Span),
		initialized: false,
	}
}

// Handle dispatches commands to appropriate handlers
func (h *TelemetryCommandHandler) Handle(ctx context.Context, cmd application.Command) error {
	switch c := cmd.(type) {
	case *application.InitializeSDKCommand:
		return h.handleInitializeSDK(ctx, c)
	case *application.ShutdownSDKCommand:
		return h.handleShutdownSDK(ctx, c)
	case *application.FlushTelemetryCommand:
		return h.handleFlushTelemetry(ctx, c)
	case *application.RecordMetricCommand:
		return h.handleRecordMetric(ctx, c)
	case *application.RecordCounterCommand:
		return h.handleRecordCounter(ctx, c)
	case *application.RecordGaugeCommand:
		return h.handleRecordGauge(ctx, c)
	case *application.RecordHistogramCommand:
		return h.handleRecordHistogram(ctx, c)
	case *application.EmitLogCommand:
		return h.handleEmitLog(ctx, c)
	case *application.StartSpanCommand:
		return h.handleStartSpan(ctx, c)
	case *application.EndSpanCommand:
		return h.handleEndSpan(ctx, c)
	case *application.AddSpanEventCommand:
		return h.handleAddSpanEvent(ctx, c)
	default:
		return fmt.Errorf("unknown command type: %T", cmd)
	}
}

// ===== SDK LIFECYCLE HANDLERS =====

func (h *TelemetryCommandHandler) handleInitializeSDK(ctx context.Context, cmd *application.InitializeSDKCommand) error {
	h.initMutex.Lock()
	defer h.initMutex.Unlock()

	if h.initialized {
		return fmt.Errorf("SDK already initialized")
	}

	h.config = cmd.Config

	// Validate configuration
	if err := h.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	factory := NewOTLPExporterFactory(h.config)

	// Create resource
	resource, err := factory.CreateResource(ctx)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// Initialize traces if enabled
	if h.config.IsSignalEnabled(domain.SignalTraces) {
		traceExporter, err := factory.CreateTraceExporter(ctx)
		if err != nil {
			return fmt.Errorf("failed to create trace exporter: %w", err)
		}

		h.tracerProvider = sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(traceExporter,
				sdktrace.WithBatchTimeout(h.config.BatchTimeout()),
				sdktrace.WithMaxExportBatchSize(h.config.BatchMaxSize()),
			),
			sdktrace.WithResource(resource),
		)
		otel.SetTracerProvider(h.tracerProvider)
		h.tracer = h.tracerProvider.Tracer(h.config.ServiceName())
	}

	// Initialize metrics if enabled
	if h.config.IsSignalEnabled(domain.SignalMetrics) {
		metricExporter, err := factory.CreateMetricExporter(ctx)
		if err != nil {
			return fmt.Errorf("failed to create metric exporter: %w", err)
		}

		h.meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter,
				sdkmetric.WithInterval(h.config.BatchTimeout()),
			)),
			sdkmetric.WithResource(resource),
		)
		otel.SetMeterProvider(h.meterProvider)
		h.meter = h.meterProvider.Meter(h.config.ServiceName())
	}

	h.initialized = true
	return nil
}

func (h *TelemetryCommandHandler) handleShutdownSDK(ctx context.Context, cmd *application.ShutdownSDKCommand) error {
	if !h.initialized {
		return fmt.Errorf("SDK not initialized")
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, cmd.Timeout)
	defer cancel()

	var shutdownErrors []error

	// Shutdown tracer provider
	if h.tracerProvider != nil {
		if err := h.tracerProvider.Shutdown(shutdownCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("tracer provider shutdown: %w", err))
		}
	}

	// Shutdown meter provider
	if h.meterProvider != nil {
		if err := h.meterProvider.Shutdown(shutdownCtx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("meter provider shutdown: %w", err))
		}
	}

	h.initialized = false

	if len(shutdownErrors) > 0 {
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	return nil
}

func (h *TelemetryCommandHandler) handleFlushTelemetry(ctx context.Context, cmd *application.FlushTelemetryCommand) error {
	if !h.initialized {
		return fmt.Errorf("SDK not initialized")
	}

	flushCtx, cancel := context.WithTimeout(ctx, cmd.Timeout)
	defer cancel()

	var flushErrors []error

	// Force flush tracer provider
	if h.tracerProvider != nil {
		if err := h.tracerProvider.ForceFlush(flushCtx); err != nil {
			flushErrors = append(flushErrors, fmt.Errorf("tracer provider flush: %w", err))
		}
	}

	// Force flush meter provider
	if h.meterProvider != nil {
		if err := h.meterProvider.ForceFlush(flushCtx); err != nil {
			flushErrors = append(flushErrors, fmt.Errorf("meter provider flush: %w", err))
		}
	}

	if len(flushErrors) > 0 {
		return fmt.Errorf("flush errors: %v", flushErrors)
	}

	return nil
}

// ===== METRIC HANDLERS =====

func (h *TelemetryCommandHandler) handleRecordMetric(ctx context.Context, cmd *application.RecordMetricCommand) error {
	if !h.initialized || h.meter == nil {
		return fmt.Errorf("metrics not initialized")
	}

	// Create a float64 gauge
	gauge, err := h.meter.Float64ObservableGauge(
		cmd.Name,
		otelmetric.WithUnit(cmd.Unit),
		otelmetric.WithDescription(fmt.Sprintf("Metric: %s", cmd.Name)),
	)
	if err != nil {
		return fmt.Errorf("failed to create gauge: %w", err)
	}

	// Record the value
	_, err = h.meter.RegisterCallback(
		func(ctx context.Context, o otelmetric.Observer) error {
			attrs := convertAttributes(cmd.Attributes)
			o.ObserveFloat64(gauge, cmd.Value, otelmetric.WithAttributes(attrs...))
			return nil
		},
		gauge,
	)

	return err
}

func (h *TelemetryCommandHandler) handleRecordCounter(ctx context.Context, cmd *application.RecordCounterCommand) error {
	if !h.initialized || h.meter == nil {
		return fmt.Errorf("metrics not initialized")
	}

	counter, err := h.meter.Int64Counter(
		cmd.Name,
		otelmetric.WithDescription(fmt.Sprintf("Counter: %s", cmd.Name)),
	)
	if err != nil {
		return fmt.Errorf("failed to create counter: %w", err)
	}

	attrs := convertAttributes(cmd.Attributes)
	counter.Add(ctx, cmd.Value, otelmetric.WithAttributes(attrs...))
	return nil
}

func (h *TelemetryCommandHandler) handleRecordGauge(ctx context.Context, cmd *application.RecordGaugeCommand) error {
	// Similar to handleRecordMetric
	return h.handleRecordMetric(ctx, &application.RecordMetricCommand{
		Name:       cmd.Name,
		Value:      cmd.Value,
		Attributes: cmd.Attributes,
		Timestamp:  time.Now(),
	})
}

func (h *TelemetryCommandHandler) handleRecordHistogram(ctx context.Context, cmd *application.RecordHistogramCommand) error {
	if !h.initialized || h.meter == nil {
		return fmt.Errorf("metrics not initialized")
	}

	histogram, err := h.meter.Float64Histogram(
		cmd.Name,
		otelmetric.WithUnit(cmd.Unit),
		otelmetric.WithDescription(fmt.Sprintf("Histogram: %s", cmd.Name)),
	)
	if err != nil {
		return fmt.Errorf("failed to create histogram: %w", err)
	}

	attrs := convertAttributes(cmd.Attributes)
	histogram.Record(ctx, cmd.Value, otelmetric.WithAttributes(attrs...))
	return nil
}

// ===== LOG HANDLERS =====

func (h *TelemetryCommandHandler) handleEmitLog(ctx context.Context, cmd *application.EmitLogCommand) error {
	// For logs, we'll use the trace span to emit log events
	// In a production SDK, you'd integrate with go.opentelemetry.io/otel/log when stable

	if h.tracer != nil {
		span := trace.SpanFromContext(ctx)
		if !span.IsRecording() {
			// If no active span, this is a standalone log
			// In production, use proper log exporter
			return fmt.Errorf("log export not fully implemented - use trace spans for now")
		}

		attrs := convertAttributes(cmd.Attributes)
		span.AddEvent(cmd.Message, trace.WithAttributes(attrs...))
	}

	return nil
}

// ===== TRACE HANDLERS =====

// StartSpanDirect starts a new span and returns its ID directly
// This is the preferred method for starting spans as it returns the span ID
func (h *TelemetryCommandHandler) StartSpanDirect(ctx context.Context, name string, kind string, attributes map[string]interface{}) (string, error) {
	if !h.initialized || h.tracer == nil {
		return "", fmt.Errorf("traces not initialized")
	}

	attrs := convertAttributes(attributes)

	var spanKind trace.SpanKind
	switch kind {
	case "internal":
		spanKind = trace.SpanKindInternal
	case "server":
		spanKind = trace.SpanKindServer
	case "client":
		spanKind = trace.SpanKindClient
	case "producer":
		spanKind = trace.SpanKindProducer
	case "consumer":
		spanKind = trace.SpanKindConsumer
	default:
		spanKind = trace.SpanKindInternal
	}

	_, span := h.tracer.Start(ctx, name,
		trace.WithSpanKind(spanKind),
		trace.WithAttributes(attrs...),
	)

	// Store active span
	spanID := span.SpanContext().SpanID().String()
	h.spansMutex.Lock()
	h.activeSpans[spanID] = span
	h.spansMutex.Unlock()

	return spanID, nil
}

func (h *TelemetryCommandHandler) handleStartSpan(ctx context.Context, cmd *application.StartSpanCommand) error {
	_, err := h.StartSpanDirect(ctx, cmd.Name, cmd.Kind, cmd.Attributes)
	return err
}

func (h *TelemetryCommandHandler) handleEndSpan(ctx context.Context, cmd *application.EndSpanCommand) error {
	h.spansMutex.Lock()
	span, exists := h.activeSpans[cmd.SpanID]
	if exists {
		delete(h.activeSpans, cmd.SpanID)
	}
	h.spansMutex.Unlock()

	if !exists {
		return fmt.Errorf("span not found: %s", cmd.SpanID)
	}

	if cmd.Error != nil {
		span.RecordError(cmd.Error)
	}

	span.End()
	return nil
}

func (h *TelemetryCommandHandler) handleAddSpanEvent(ctx context.Context, cmd *application.AddSpanEventCommand) error {
	h.spansMutex.RLock()
	span, exists := h.activeSpans[cmd.SpanID]
	h.spansMutex.RUnlock()

	if !exists {
		return fmt.Errorf("span not found: %s", cmd.SpanID)
	}

	attrs := convertAttributes(cmd.Attributes)
	span.AddEvent(cmd.Name,
		trace.WithTimestamp(cmd.Timestamp),
		trace.WithAttributes(attrs...),
	)

	return nil
}

// ===== HELPER FUNCTIONS =====

func convertAttributes(attrs map[string]interface{}) []attribute.KeyValue {
	result := make([]attribute.KeyValue, 0, len(attrs))
	for key, value := range attrs {
		switch v := value.(type) {
		case string:
			result = append(result, attribute.String(key, v))
		case int:
			result = append(result, attribute.Int(key, v))
		case int64:
			result = append(result, attribute.Int64(key, v))
		case float64:
			result = append(result, attribute.Float64(key, v))
		case bool:
			result = append(result, attribute.Bool(key, v))
		default:
			result = append(result, attribute.String(key, fmt.Sprintf("%v", v)))
		}
	}
	return result
}
