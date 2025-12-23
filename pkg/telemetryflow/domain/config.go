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

// TelemetryConfig is an aggregate root that contains all configuration
// This is the main Domain Entity in our DDD design
type TelemetryConfig struct {
	// Identity
	credentials *Credentials

	// Connection settings
	endpoint         string
	protocol         Protocol
	insecure         bool
	timeout          time.Duration
	retryEnabled     bool
	maxRetries       int
	retryBackoff     time.Duration
	compressionGzip  bool

	// Signal configuration
	enabledSignals   map[SignalType]bool

	// Resource attributes
	serviceName      string
	serviceVersion   string
	environment      string
	customAttributes map[string]string

	// Batch settings
	batchTimeout     time.Duration
	batchMaxSize     int

	// Rate limiting (client-side)
	rateLimit        int // requests per minute
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
		credentials:      credentials,
		endpoint:         endpoint,
		protocol:         ProtocolGRPC, // default
		insecure:         false,
		timeout:          30 * time.Second,
		retryEnabled:     true,
		maxRetries:       3,
		retryBackoff:     5 * time.Second,
		compressionGzip:  true,
		enabledSignals: map[SignalType]bool{
			SignalMetrics: true,
			SignalLogs:    true,
			SignalTraces:  true,
		},
		serviceName:      serviceName,
		serviceVersion:   "1.0.0",
		environment:      "production",
		customAttributes: make(map[string]string),
		batchTimeout:     10 * time.Second,
		batchMaxSize:     512,
		rateLimit:        1000,
	}, nil
}

// Getters for immutability
func (c *TelemetryConfig) Credentials() *Credentials        { return c.credentials }
func (c *TelemetryConfig) Endpoint() string                 { return c.endpoint }
func (c *TelemetryConfig) Protocol() Protocol               { return c.protocol }
func (c *TelemetryConfig) IsInsecure() bool                 { return c.insecure }
func (c *TelemetryConfig) Timeout() time.Duration           { return c.timeout }
func (c *TelemetryConfig) IsRetryEnabled() bool             { return c.retryEnabled }
func (c *TelemetryConfig) MaxRetries() int                  { return c.maxRetries }
func (c *TelemetryConfig) RetryBackoff() time.Duration      { return c.retryBackoff }
func (c *TelemetryConfig) IsCompressionEnabled() bool       { return c.compressionGzip }
func (c *TelemetryConfig) ServiceName() string              { return c.serviceName }
func (c *TelemetryConfig) ServiceVersion() string           { return c.serviceVersion }
func (c *TelemetryConfig) Environment() string              { return c.environment }
func (c *TelemetryConfig) CustomAttributes() map[string]string { return c.customAttributes }
func (c *TelemetryConfig) BatchTimeout() time.Duration      { return c.batchTimeout }
func (c *TelemetryConfig) BatchMaxSize() int                { return c.batchMaxSize }
func (c *TelemetryConfig) RateLimit() int                   { return c.rateLimit }

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
