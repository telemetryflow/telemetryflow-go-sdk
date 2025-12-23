package application

import (
	"context"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

// Command interface - marker interface for all commands
// Commands are handled by CommandHandler.Handle() which does the actual execution
type Command interface {
	isCommand() // marker method
}

// CommandHandler interface for handling commands
type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}

// ===== METRIC COMMANDS =====

// RecordMetricCommand records a metric data point
type RecordMetricCommand struct {
	Name       string
	Value      float64
	Unit       string
	Attributes map[string]interface{}
	Timestamp  time.Time
}

func (*RecordMetricCommand) isCommand() {}

// RecordCounterCommand increments a counter metric
type RecordCounterCommand struct {
	Name       string
	Value      int64
	Attributes map[string]interface{}
}

func (*RecordCounterCommand) isCommand() {}

// RecordGaugeCommand sets a gauge metric value
type RecordGaugeCommand struct {
	Name       string
	Value      float64
	Attributes map[string]interface{}
}

func (*RecordGaugeCommand) isCommand() {}

// RecordHistogramCommand records a histogram measurement
type RecordHistogramCommand struct {
	Name       string
	Value      float64
	Unit       string
	Attributes map[string]interface{}
}

func (*RecordHistogramCommand) isCommand() {}

// ===== LOG COMMANDS =====

// EmitLogCommand emits a structured log entry
type EmitLogCommand struct {
	Severity   string
	Message    string
	Attributes map[string]interface{}
	Timestamp  time.Time
	TraceID    string // Optional trace correlation
	SpanID     string // Optional span correlation
}

func (*EmitLogCommand) isCommand() {}

// EmitBatchLogsCommand emits multiple logs in a batch
type EmitBatchLogsCommand struct {
	Logs []EmitLogCommand
}

func (*EmitBatchLogsCommand) isCommand() {}

// ===== TRACE COMMANDS =====

// StartSpanCommand starts a new trace span
type StartSpanCommand struct {
	Name       string
	Kind       string // internal, server, client, producer, consumer
	Attributes map[string]interface{}
	ParentID   string // Optional parent span ID
}

func (*StartSpanCommand) isCommand() {}

// EndSpanCommand ends an active span
type EndSpanCommand struct {
	SpanID string
	Error  error // Optional error to record
}

func (*EndSpanCommand) isCommand() {}

// AddSpanEventCommand adds an event to a span
type AddSpanEventCommand struct {
	SpanID     string
	Name       string
	Attributes map[string]interface{}
	Timestamp  time.Time
}

func (*AddSpanEventCommand) isCommand() {}

// ===== INITIALIZATION COMMANDS =====

// InitializeSDKCommand initializes the SDK with configuration
type InitializeSDKCommand struct {
	Config *domain.TelemetryConfig
}

func (*InitializeSDKCommand) isCommand() {}

// ShutdownSDKCommand gracefully shuts down the SDK
type ShutdownSDKCommand struct {
	Timeout time.Duration
}

func (*ShutdownSDKCommand) isCommand() {}

// FlushTelemetryCommand forces a flush of all pending telemetry
type FlushTelemetryCommand struct {
	Timeout time.Duration
}

func (*FlushTelemetryCommand) isCommand() {}

// ===== COMMAND BUS =====

// CommandBus dispatches commands to handlers
type CommandBus struct {
	handlers map[string]CommandHandler
}

// NewCommandBus creates a new command bus
func NewCommandBus() *CommandBus {
	return &CommandBus{
		handlers: make(map[string]CommandHandler),
	}
}

// Register registers a command handler
func (b *CommandBus) Register(commandType string, handler CommandHandler) {
	b.handlers[commandType] = handler
}

// Dispatch dispatches a command to its handler
func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	// Implementation will delegate to specific handlers
	// For now, this is the structure
	return nil
}
