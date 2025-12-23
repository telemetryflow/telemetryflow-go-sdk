package telemetryflow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/application"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/infrastructure"
)

// Client is the main SDK client for TelemetryFlow
type Client struct {
	config         *domain.TelemetryConfig
	commandBus     *application.CommandBus
	queryBus       *application.QueryBus
	commandHandler *infrastructure.TelemetryCommandHandler
	initialized    bool
	mu             sync.RWMutex
}

// NewClient creates a new TelemetryFlow client
func NewClient(config *domain.TelemetryConfig) (*Client, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	commandHandler := infrastructure.NewTelemetryCommandHandler(config)

	return &Client{
		config:         config,
		commandBus:     application.NewCommandBus(),
		queryBus:       application.NewQueryBus(),
		commandHandler: commandHandler,
		initialized:    false,
	}, nil
}

// Initialize initializes the SDK and starts exporters
func (c *Client) Initialize(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return fmt.Errorf("client already initialized")
	}

	cmd := &application.InitializeSDKCommand{
		Config: c.config,
	}

	if err := c.commandHandler.Handle(ctx, cmd); err != nil {
		return fmt.Errorf("failed to initialize SDK: %w", err)
	}

	c.initialized = true
	return nil
}

// Shutdown gracefully shuts down the SDK
func (c *Client) Shutdown(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.initialized {
		return nil
	}

	cmd := &application.ShutdownSDKCommand{
		Timeout: 30 * time.Second,
	}

	if err := c.commandHandler.Handle(ctx, cmd); err != nil {
		return fmt.Errorf("failed to shutdown SDK: %w", err)
	}

	c.initialized = false
	return nil
}

// Flush forces a flush of all pending telemetry
func (c *Client) Flush(ctx context.Context) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.FlushTelemetryCommand{
		Timeout: 10 * time.Second,
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// ===== METRICS API =====

// RecordMetric records a generic metric
func (c *Client) RecordMetric(ctx context.Context, name string, value float64, unit string, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.RecordMetricCommand{
		Name:       name,
		Value:      value,
		Unit:       unit,
		Attributes: attributes,
		Timestamp:  time.Now(),
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// IncrementCounter increments a counter metric
func (c *Client) IncrementCounter(ctx context.Context, name string, value int64, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.RecordCounterCommand{
		Name:       name,
		Value:      value,
		Attributes: attributes,
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// RecordGauge records a gauge metric
func (c *Client) RecordGauge(ctx context.Context, name string, value float64, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.RecordGaugeCommand{
		Name:       name,
		Value:      value,
		Attributes: attributes,
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// RecordHistogram records a histogram measurement
func (c *Client) RecordHistogram(ctx context.Context, name string, value float64, unit string, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.RecordHistogramCommand{
		Name:       name,
		Value:      value,
		Unit:       unit,
		Attributes: attributes,
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// ===== LOGS API =====

// Log emits a structured log entry
func (c *Client) Log(ctx context.Context, severity string, message string, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.EmitLogCommand{
		Severity:   severity,
		Message:    message,
		Attributes: attributes,
		Timestamp:  time.Now(),
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// LogInfo emits an info-level log
func (c *Client) LogInfo(ctx context.Context, message string, attributes map[string]interface{}) error {
	return c.Log(ctx, "info", message, attributes)
}

// LogWarn emits a warning-level log
func (c *Client) LogWarn(ctx context.Context, message string, attributes map[string]interface{}) error {
	return c.Log(ctx, "warn", message, attributes)
}

// LogError emits an error-level log
func (c *Client) LogError(ctx context.Context, message string, attributes map[string]interface{}) error {
	return c.Log(ctx, "error", message, attributes)
}

// ===== TRACES API =====

// StartSpan starts a new trace span
func (c *Client) StartSpan(ctx context.Context, name string, kind string, attributes map[string]interface{}) (string, error) {
	if !c.isInitialized() {
		return "", fmt.Errorf("client not initialized")
	}

	cmd := &application.StartSpanCommand{
		Name:       name,
		Kind:       kind,
		Attributes: attributes,
	}

	if err := c.commandHandler.Handle(ctx, cmd); err != nil {
		return "", err
	}

	// In production, return the actual span ID
	return "span-id-placeholder", nil
}

// EndSpan ends an active span
func (c *Client) EndSpan(ctx context.Context, spanID string, err error) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.EndSpanCommand{
		SpanID: spanID,
		Error:  err,
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// AddSpanEvent adds an event to an active span
func (c *Client) AddSpanEvent(ctx context.Context, spanID string, name string, attributes map[string]interface{}) error {
	if !c.isInitialized() {
		return fmt.Errorf("client not initialized")
	}

	cmd := &application.AddSpanEventCommand{
		SpanID:     spanID,
		Name:       name,
		Attributes: attributes,
		Timestamp:  time.Now(),
	}

	return c.commandHandler.Handle(ctx, cmd)
}

// ===== HELPER METHODS =====

func (c *Client) isInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}

// Config returns the client configuration
func (c *Client) Config() *domain.TelemetryConfig {
	return c.config
}

// IsInitialized returns whether the client is initialized
func (c *Client) IsInitialized() bool {
	return c.isInitialized()
}
