package application

import (
	"context"
	"time"
)

// Query interface - all queries implement this
type Query interface {
	Execute(ctx context.Context) (interface{}, error)
}

// QueryHandler interface for handling queries
type QueryHandler interface {
	Handle(ctx context.Context, query Query) (interface{}, error)
}

// ===== METRIC QUERIES =====

// GetMetricQuery retrieves a specific metric
type GetMetricQuery struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	GroupBy   []string
	Filters   map[string]string
}

// MetricQueryResult represents the result of a metric query
type MetricQueryResult struct {
	Name       string
	DataPoints []MetricDataPoint
}

// MetricDataPoint represents a single metric data point
type MetricDataPoint struct {
	Timestamp  time.Time
	Value      float64
	Attributes map[string]interface{}
}

// AggregateMetricsQuery aggregates metrics over a time range
type AggregateMetricsQuery struct {
	Name        string
	StartTime   time.Time
	EndTime     time.Time
	Aggregation string // sum, avg, min, max, count
	GroupBy     []string
	Filters     map[string]string
}

// AggregateMetricsResult represents aggregated metric results
type AggregateMetricsResult struct {
	Name   string
	Value  float64
	Groups []AggregatedGroup
}

// AggregatedGroup represents a grouped aggregation result
type AggregatedGroup struct {
	GroupBy map[string]string
	Value   float64
}

// ===== LOG QUERIES =====

// GetLogsQuery retrieves logs based on criteria
type GetLogsQuery struct {
	StartTime  time.Time
	EndTime    time.Time
	Severity   []string // filter by severity levels
	SearchText string   // full-text search
	Filters    map[string]string
	Limit      int
	Offset     int
}

// LogsQueryResult represents the result of a logs query
type LogsQueryResult struct {
	Logs       []LogEntry
	TotalCount int
	HasMore    bool
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp  time.Time
	Severity   string
	Message    string
	Attributes map[string]interface{}
	TraceID    string
	SpanID     string
}

// ===== TRACE QUERIES =====

// GetTraceQuery retrieves a complete trace by ID
type GetTraceQuery struct {
	TraceID string
}

// TraceQueryResult represents a complete trace
type TraceQueryResult struct {
	TraceID   string
	Spans     []SpanEntry
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
}

// SpanEntry represents a single span in a trace
type SpanEntry struct {
	SpanID     string
	ParentID   string
	Name       string
	Kind       string
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Attributes map[string]interface{}
	Events     []SpanEvent
	Status     SpanStatus
}

// SpanEvent represents an event within a span
type SpanEvent struct {
	Name       string
	Timestamp  time.Time
	Attributes map[string]interface{}
}

// SpanStatus represents the status of a span
type SpanStatus struct {
	Code    string
	Message string
}

// SearchTracesQuery searches for traces based on criteria
type SearchTracesQuery struct {
	StartTime   time.Time
	EndTime     time.Time
	ServiceName string
	Operation   string
	MinDuration time.Duration
	MaxDuration time.Duration
	Tags        map[string]string
	HasError    *bool // filter by error status
	Limit       int
	Offset      int
}

// TracesSearchResult represents trace search results
type TracesSearchResult struct {
	Traces     []TraceSummary
	TotalCount int
	HasMore    bool
}

// TraceSummary represents a summary of a trace
type TraceSummary struct {
	TraceID     string
	RootSpan    string
	ServiceName string
	Duration    time.Duration
	StartTime   time.Time
	SpanCount   int
	ErrorCount  int
}

// ===== HEALTH & STATUS QUERIES =====

// GetHealthQuery checks the health of the SDK connection
type GetHealthQuery struct{}

// HealthQueryResult represents health status
type HealthQueryResult struct {
	Status      string // healthy, degraded, unhealthy
	LastSuccess time.Time
	LastError   error
	Metrics     HealthMetrics
}

// HealthMetrics represents health-related metrics
type HealthMetrics struct {
	TotalRequests   int64
	FailedRequests  int64
	AverageLatency  time.Duration
	ConnectionState string
}

// GetSDKStatusQuery gets the current SDK status
type GetSDKStatusQuery struct{}

// SDKStatusResult represents SDK status
type SDKStatusResult struct {
	Initialized    bool
	Version        string
	EnabledSignals []string
	Config         map[string]interface{}
	Statistics     SDKStatistics
}

// SDKStatistics represents SDK statistics
type SDKStatistics struct {
	MetricsSent int64
	LogsSent    int64
	TracesSent  int64
	ErrorsCount int64
	LastFlush   time.Time
	QueueSize   int
}

// ===== QUERY BUS =====

// QueryBus dispatches queries to handlers
type QueryBus struct {
	handlers map[string]QueryHandler
}

// NewQueryBus creates a new query bus
func NewQueryBus() *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler),
	}
}

// Register registers a query handler
func (b *QueryBus) Register(queryType string, handler QueryHandler) {
	b.handlers[queryType] = handler
}

// Dispatch dispatches a query to its handler
func (b *QueryBus) Dispatch(ctx context.Context, query Query) (interface{}, error) {
	// Implementation will delegate to specific handlers
	// For now, this is the structure
	return nil, nil
}
