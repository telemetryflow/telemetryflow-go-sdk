// Package domain provides core domain types for the TelemetryFlow SDK.
package domain

import (
	"errors"
	"fmt"
	"time"
)

// Protocol represents the OTLP protocol type
type Protocol string

const (
	ProtocolGRPC Protocol = "grpc"
	ProtocolHTTP Protocol = "http"
)

// SignalType represents the type of telemetry signal
type SignalType string

const (
	SignalMetrics SignalType = "metrics"
	SignalLogs    SignalType = "logs"
	SignalTraces  SignalType = "traces"
)

// GRPCKeepaliveConfig holds gRPC keepalive settings
type GRPCKeepaliveConfig struct {
	Time                time.Duration
	Timeout             time.Duration
	PermitWithoutStream bool
}

// TelemetryConfig is an aggregate root that contains all configuration
// This is the main Domain Entity in our DDD design
type TelemetryConfig struct {
	// Identity
	credentials *Credentials
	collectorID string // Unique collector identifier for TelemetryFlow headers

	// Collector Identity (aligned with tfoidentityextension)
	collectorName        string            // Human-readable collector name
	collectorDescription string            // Human-readable collector description
	collectorHostname    string            // Collector hostname (auto-detected if empty)
	collectorTags        map[string]string // Custom key-value pairs for labeling
	enrichResources      bool              // Add collector identity to all telemetry resources

	// Connection settings
	endpoint        string
	protocol        Protocol
	insecure        bool
	timeout         time.Duration
	retryEnabled    bool
	maxRetries      int
	retryBackoff    time.Duration
	compressionGzip bool

	// TFO API Version settings (aligned with tfoexporter)
	useV2API        bool   // Use v2 API endpoints (/v2/traces, /v2/metrics, /v2/logs)
	v2Only          bool   // v2-only mode - reject v1 endpoints
	tracesEndpoint  string // Custom traces endpoint path override
	metricsEndpoint string // Custom metrics endpoint path override
	logsEndpoint    string // Custom logs endpoint path override

	// gRPC specific settings
	grpcKeepalive       *GRPCKeepaliveConfig
	grpcMaxRecvMsgSize  int // in MiB
	grpcMaxSendMsgSize  int // in MiB
	grpcReadBufferSize  int // in bytes
	grpcWriteBufferSize int // in bytes

	// Signal configuration
	enabledSignals map[SignalType]bool

	// Resource attributes
	serviceName      string
	serviceNamespace string
	serviceVersion   string
	environment      string
	datacenter       string // Datacenter identifier
	customAttributes map[string]string

	// Batch settings
	batchTimeout time.Duration
	batchMaxSize int

	// Rate limiting (client-side)
	rateLimit int // requests per minute

	// Exemplars support (for metrics-to-traces correlation)
	exemplarsEnabled bool
}

// NewTelemetryConfig creates a new configuration with required fields
func NewTelemetryConfig(credentials *Credentials, endpoint string, serviceName string) (*TelemetryConfig, error) {
	if credentials == nil {
		return nil, errors.New("credentials cannot be nil")
	}
	if endpoint == "" {
		return nil, errors.New("endpoint cannot be empty")
	}
	if serviceName == "" {
		return nil, errors.New("service name cannot be empty")
	}

	return &TelemetryConfig{
		credentials:     credentials,
		collectorID:     "", // auto-generated if empty
		endpoint:        endpoint,
		protocol:        ProtocolGRPC, // default
		insecure:        false,
		timeout:         30 * time.Second,
		retryEnabled:    true,
		maxRetries:      3,
		retryBackoff:    5 * time.Second,
		compressionGzip: true,
		// TFO API Version settings (aligned with tfoexporter)
		useV2API:        true, // v2 API enabled by default for TFO Platform
		v2Only:          false,
		tracesEndpoint:  "", // use defaults: /v2/traces or /v1/traces
		metricsEndpoint: "", // use defaults: /v2/metrics or /v1/metrics
		logsEndpoint:    "", // use defaults: /v2/logs or /v1/logs
		// Collector Identity (aligned with tfoidentityextension)
		collectorName:        "",
		collectorDescription: "",
		collectorHostname:    "", // auto-detected
		collectorTags:        make(map[string]string),
		enrichResources:      true, // enabled by default
		// gRPC settings aligned with OTEL Collector config
		grpcKeepalive: &GRPCKeepaliveConfig{
			Time:                10 * time.Second,
			Timeout:             5 * time.Second,
			PermitWithoutStream: true,
		},
		grpcMaxRecvMsgSize:  4,      // 4 MiB
		grpcMaxSendMsgSize:  4,      // 4 MiB
		grpcReadBufferSize:  524288, // 512 KB
		grpcWriteBufferSize: 524288, // 512 KB
		enabledSignals: map[SignalType]bool{
			SignalMetrics: true,
			SignalLogs:    true,
			SignalTraces:  true,
		},
		serviceName:      serviceName,
		serviceNamespace: "telemetryflow", // default namespace
		serviceVersion:   "1.0.0",
		environment:      "production",
		datacenter:       "default",
		customAttributes: make(map[string]string),
		batchTimeout:     10 * time.Second,
		batchMaxSize:     512,
		rateLimit:        1000,
		exemplarsEnabled: true, // enabled by default for metrics-to-traces correlation
	}, nil
}

// Credentials returns the API credentials for TelemetryFlow authentication.
func (c *TelemetryConfig) Credentials() *Credentials { return c.credentials }

// CollectorID returns the unique identifier for this collector instance.
func (c *TelemetryConfig) CollectorID() string { return c.collectorID }

// Endpoint returns the OTLP collector endpoint address.
func (c *TelemetryConfig) Endpoint() string { return c.endpoint }

// Protocol returns the OTLP protocol type (gRPC or HTTP).
func (c *TelemetryConfig) Protocol() Protocol { return c.protocol }

// IsInsecure returns true if TLS verification is disabled.
func (c *TelemetryConfig) IsInsecure() bool { return c.insecure }

// Timeout returns the connection timeout duration.
func (c *TelemetryConfig) Timeout() time.Duration { return c.timeout }

// IsRetryEnabled returns true if automatic retries are enabled.
func (c *TelemetryConfig) IsRetryEnabled() bool { return c.retryEnabled }

// MaxRetries returns the maximum number of retry attempts.
func (c *TelemetryConfig) MaxRetries() int { return c.maxRetries }

// RetryBackoff returns the backoff duration between retries.
func (c *TelemetryConfig) RetryBackoff() time.Duration { return c.retryBackoff }

// IsCompressionEnabled returns true if gzip compression is enabled.
func (c *TelemetryConfig) IsCompressionEnabled() bool { return c.compressionGzip }

// GRPCKeepalive returns the gRPC keepalive configuration.
func (c *TelemetryConfig) GRPCKeepalive() *GRPCKeepaliveConfig { return c.grpcKeepalive }

// GRPCMaxRecvMsgSize returns the maximum gRPC receive message size in bytes.
func (c *TelemetryConfig) GRPCMaxRecvMsgSize() int { return c.grpcMaxRecvMsgSize }

// GRPCMaxSendMsgSize returns the maximum gRPC send message size in bytes.
func (c *TelemetryConfig) GRPCMaxSendMsgSize() int { return c.grpcMaxSendMsgSize }

// GRPCReadBufferSize returns the gRPC read buffer size in bytes.
func (c *TelemetryConfig) GRPCReadBufferSize() int { return c.grpcReadBufferSize }

// GRPCWriteBufferSize returns the gRPC write buffer size in bytes.
func (c *TelemetryConfig) GRPCWriteBufferSize() int { return c.grpcWriteBufferSize }

// ServiceName returns the name of the instrumented service.
func (c *TelemetryConfig) ServiceName() string { return c.serviceName }

// ServiceNamespace returns the namespace of the instrumented service.
func (c *TelemetryConfig) ServiceNamespace() string { return c.serviceNamespace }

// ServiceVersion returns the version of the instrumented service.
func (c *TelemetryConfig) ServiceVersion() string { return c.serviceVersion }

// Environment returns the deployment environment (e.g., production, staging).
func (c *TelemetryConfig) Environment() string { return c.environment }

// CustomAttributes returns the custom attributes to add to all telemetry data.
func (c *TelemetryConfig) CustomAttributes() map[string]string { return c.customAttributes }

// BatchTimeout returns the batch export timeout duration.
func (c *TelemetryConfig) BatchTimeout() time.Duration { return c.batchTimeout }

// BatchMaxSize returns the maximum batch size for export.
func (c *TelemetryConfig) BatchMaxSize() int { return c.batchMaxSize }

// RateLimit returns the rate limit for telemetry data export.
func (c *TelemetryConfig) RateLimit() int { return c.rateLimit }

// IsExemplarsEnabled returns true if exemplars are enabled for metrics-to-traces correlation.
func (c *TelemetryConfig) IsExemplarsEnabled() bool { return c.exemplarsEnabled }

// UseV2API returns true if v2 API endpoints are enabled.
func (c *TelemetryConfig) UseV2API() bool { return c.useV2API }

// IsV2Only returns true if only v2 endpoints should be used.
func (c *TelemetryConfig) IsV2Only() bool { return c.v2Only }

// TracesEndpoint returns the custom traces endpoint path, or default based on v2 setting.
func (c *TelemetryConfig) TracesEndpoint() string {
	if c.tracesEndpoint != "" {
		return c.tracesEndpoint
	}
	if c.useV2API {
		return "/v2/traces"
	}
	return "/v1/traces"
}

// MetricsEndpoint returns the custom metrics endpoint path, or default based on v2 setting.
func (c *TelemetryConfig) MetricsEndpoint() string {
	if c.metricsEndpoint != "" {
		return c.metricsEndpoint
	}
	if c.useV2API {
		return "/v2/metrics"
	}
	return "/v1/metrics"
}

// LogsEndpoint returns the custom logs endpoint path, or default based on v2 setting.
func (c *TelemetryConfig) LogsEndpoint() string {
	if c.logsEndpoint != "" {
		return c.logsEndpoint
	}
	if c.useV2API {
		return "/v2/logs"
	}
	return "/v1/logs"
}

// CollectorName returns the human-readable collector name.
func (c *TelemetryConfig) CollectorName() string { return c.collectorName }

// CollectorDescription returns the collector description.
func (c *TelemetryConfig) CollectorDescription() string { return c.collectorDescription }

// CollectorHostname returns the collector hostname.
func (c *TelemetryConfig) CollectorHostname() string { return c.collectorHostname }

// CollectorTags returns the collector tags for labeling and filtering.
func (c *TelemetryConfig) CollectorTags() map[string]string { return c.collectorTags }

// IsEnrichResourcesEnabled returns true if collector identity enrichment is enabled.
func (c *TelemetryConfig) IsEnrichResourcesEnabled() bool { return c.enrichResources }

// Datacenter returns the datacenter identifier.
func (c *TelemetryConfig) Datacenter() string { return c.datacenter }

// IsSignalEnabled checks if a signal type is enabled
func (c *TelemetryConfig) IsSignalEnabled(signal SignalType) bool {
	return c.enabledSignals[signal]
}

// Builder pattern methods for configuration

// WithProtocol sets the protocol
func (c *TelemetryConfig) WithProtocol(protocol Protocol) *TelemetryConfig {
	c.protocol = protocol
	return c
}

// WithInsecure sets insecure connection
func (c *TelemetryConfig) WithInsecure(insecure bool) *TelemetryConfig {
	c.insecure = insecure
	return c
}

// WithTimeout sets connection timeout
func (c *TelemetryConfig) WithTimeout(timeout time.Duration) *TelemetryConfig {
	c.timeout = timeout
	return c
}

// WithRetry configures retry behavior
func (c *TelemetryConfig) WithRetry(enabled bool, maxRetries int, backoff time.Duration) *TelemetryConfig {
	c.retryEnabled = enabled
	c.maxRetries = maxRetries
	c.retryBackoff = backoff
	return c
}

// WithCompression enables/disables gzip compression
func (c *TelemetryConfig) WithCompression(enabled bool) *TelemetryConfig {
	c.compressionGzip = enabled
	return c
}

// WithSignals configures which signals to enable
func (c *TelemetryConfig) WithSignals(metrics, logs, traces bool) *TelemetryConfig {
	c.enabledSignals[SignalMetrics] = metrics
	c.enabledSignals[SignalLogs] = logs
	c.enabledSignals[SignalTraces] = traces
	return c
}

// WithServiceVersion sets the service version
func (c *TelemetryConfig) WithServiceVersion(version string) *TelemetryConfig {
	c.serviceVersion = version
	return c
}

// WithEnvironment sets the deployment environment
func (c *TelemetryConfig) WithEnvironment(env string) *TelemetryConfig {
	c.environment = env
	return c
}

// WithCustomAttribute adds a custom resource attribute
func (c *TelemetryConfig) WithCustomAttribute(key, value string) *TelemetryConfig {
	c.customAttributes[key] = value
	return c
}

// WithBatchSettings configures batch export settings
func (c *TelemetryConfig) WithBatchSettings(timeout time.Duration, maxSize int) *TelemetryConfig {
	c.batchTimeout = timeout
	c.batchMaxSize = maxSize
	return c
}

// WithRateLimit sets client-side rate limit (requests per minute)
func (c *TelemetryConfig) WithRateLimit(limit int) *TelemetryConfig {
	c.rateLimit = limit
	return c
}

// WithCollectorID sets the collector identifier for TelemetryFlow headers
func (c *TelemetryConfig) WithCollectorID(id string) *TelemetryConfig {
	c.collectorID = id
	return c
}

// WithServiceNamespace sets the service namespace
func (c *TelemetryConfig) WithServiceNamespace(namespace string) *TelemetryConfig {
	c.serviceNamespace = namespace
	return c
}

// WithGRPCKeepalive configures gRPC keepalive settings
func (c *TelemetryConfig) WithGRPCKeepalive(time, timeout time.Duration, permitWithoutStream bool) *TelemetryConfig {
	c.grpcKeepalive = &GRPCKeepaliveConfig{
		Time:                time,
		Timeout:             timeout,
		PermitWithoutStream: permitWithoutStream,
	}
	return c
}

// WithGRPCBufferSizes sets gRPC read/write buffer sizes
func (c *TelemetryConfig) WithGRPCBufferSizes(readSize, writeSize int) *TelemetryConfig {
	c.grpcReadBufferSize = readSize
	c.grpcWriteBufferSize = writeSize
	return c
}

// WithGRPCMessageSizes sets gRPC max recv/send message sizes in MiB
func (c *TelemetryConfig) WithGRPCMessageSizes(recvSize, sendSize int) *TelemetryConfig {
	c.grpcMaxRecvMsgSize = recvSize
	c.grpcMaxSendMsgSize = sendSize
	return c
}

// WithExemplars enables/disables exemplars for metrics-to-traces correlation
func (c *TelemetryConfig) WithExemplars(enabled bool) *TelemetryConfig {
	c.exemplarsEnabled = enabled
	return c
}

// WithV2API enables/disables v2 API endpoints (aligned with tfoexporter)
func (c *TelemetryConfig) WithV2API(enabled bool) *TelemetryConfig {
	c.useV2API = enabled
	return c
}

// WithV2Only sets v2-only mode (rejects v1 endpoints)
func (c *TelemetryConfig) WithV2Only(enabled bool) *TelemetryConfig {
	c.v2Only = enabled
	if enabled {
		c.useV2API = true // v2Only implies useV2API
	}
	return c
}

// WithTracesEndpoint sets a custom traces endpoint path
func (c *TelemetryConfig) WithTracesEndpoint(endpoint string) *TelemetryConfig {
	c.tracesEndpoint = endpoint
	return c
}

// WithMetricsEndpoint sets a custom metrics endpoint path
func (c *TelemetryConfig) WithMetricsEndpoint(endpoint string) *TelemetryConfig {
	c.metricsEndpoint = endpoint
	return c
}

// WithLogsEndpoint sets a custom logs endpoint path
func (c *TelemetryConfig) WithLogsEndpoint(endpoint string) *TelemetryConfig {
	c.logsEndpoint = endpoint
	return c
}

// WithCollectorName sets the human-readable collector name (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithCollectorName(name string) *TelemetryConfig {
	c.collectorName = name
	return c
}

// WithCollectorDescription sets the collector description (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithCollectorDescription(description string) *TelemetryConfig {
	c.collectorDescription = description
	return c
}

// WithCollectorHostname sets the collector hostname (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithCollectorHostname(hostname string) *TelemetryConfig {
	c.collectorHostname = hostname
	return c
}

// WithCollectorTag adds a single collector tag (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithCollectorTag(key, value string) *TelemetryConfig {
	c.collectorTags[key] = value
	return c
}

// WithCollectorTags sets all collector tags (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithCollectorTags(tags map[string]string) *TelemetryConfig {
	c.collectorTags = tags
	return c
}

// WithEnrichResources enables/disables collector identity enrichment (aligned with tfoidentityextension)
func (c *TelemetryConfig) WithEnrichResources(enabled bool) *TelemetryConfig {
	c.enrichResources = enabled
	return c
}

// WithDatacenter sets the datacenter identifier
func (c *TelemetryConfig) WithDatacenter(datacenter string) *TelemetryConfig {
	c.datacenter = datacenter
	return c
}

// Validate ensures the configuration is valid
func (c *TelemetryConfig) Validate() error {
	if c.credentials == nil {
		return errors.New("credentials cannot be nil")
	}
	if c.endpoint == "" {
		return errors.New("endpoint cannot be empty")
	}
	if c.serviceName == "" {
		return errors.New("service name cannot be empty")
	}
	if c.timeout <= 0 {
		return errors.New("timeout must be positive")
	}
	if c.maxRetries < 0 {
		return errors.New("max retries cannot be negative")
	}
	if c.batchMaxSize <= 0 {
		return errors.New("batch max size must be positive")
	}
	if c.rateLimit < 0 {
		return errors.New("rate limit cannot be negative")
	}
	return nil
}

// String returns a string representation of the configuration
func (c *TelemetryConfig) String() string {
	return fmt.Sprintf(
		"TelemetryConfig{endpoint: %s, protocol: %s, service: %s, env: %s}",
		c.endpoint, c.protocol, c.serviceName, c.environment,
	)
}
