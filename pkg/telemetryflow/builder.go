package telemetryflow

import (
	"fmt"
	"os"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

// Builder provides a fluent interface for creating TelemetryFlow clients
type Builder struct {
	apiKeyID       string
	apiKeySecret   string
	endpoint       string
	serviceName    string
	serviceVersion string
	environment    string
	protocol       domain.Protocol
	insecure       bool
	timeout        time.Duration
	enableMetrics  bool
	enableLogs     bool
	enableTraces   bool
	customAttrs    map[string]string
	errors         []error
}

// NewBuilder creates a new SDK builder
func NewBuilder() *Builder {
	return &Builder{
		protocol:      domain.ProtocolGRPC,
		insecure:      false,
		timeout:       30 * time.Second,
		enableMetrics: true,
		enableLogs:    true,
		enableTraces:  true,
		customAttrs:   make(map[string]string),
		errors:        make([]error, 0),
	}
}

// WithAPIKey sets the API credentials
func (b *Builder) WithAPIKey(keyID, keySecret string) *Builder {
	b.apiKeyID = keyID
	b.apiKeySecret = keySecret
	return b
}

// WithAPIKeyFromEnv reads API credentials from environment variables
func (b *Builder) WithAPIKeyFromEnv() *Builder {
	b.apiKeyID = os.Getenv("TELEMETRYFLOW_API_KEY_ID")
	b.apiKeySecret = os.Getenv("TELEMETRYFLOW_API_KEY_SECRET")

	if b.apiKeyID == "" {
		b.errors = append(b.errors, fmt.Errorf("TELEMETRYFLOW_API_KEY_ID environment variable not set"))
	}
	if b.apiKeySecret == "" {
		b.errors = append(b.errors, fmt.Errorf("TELEMETRYFLOW_API_KEY_SECRET environment variable not set"))
	}

	return b
}

// WithEndpoint sets the OTLP endpoint
func (b *Builder) WithEndpoint(endpoint string) *Builder {
	b.endpoint = endpoint
	return b
}

// WithEndpointFromEnv reads endpoint from environment variable
func (b *Builder) WithEndpointFromEnv() *Builder {
	endpoint := os.Getenv("TELEMETRYFLOW_ENDPOINT")
	if endpoint == "" {
		endpoint = "api.telemetryflow.id:4317" // default
	}
	b.endpoint = endpoint
	return b
}

// WithService sets the service name and version
func (b *Builder) WithService(name, version string) *Builder {
	b.serviceName = name
	b.serviceVersion = version
	return b
}

// WithServiceFromEnv reads service info from environment variables
func (b *Builder) WithServiceFromEnv() *Builder {
	b.serviceName = os.Getenv("TELEMETRYFLOW_SERVICE_NAME")
	b.serviceVersion = os.Getenv("TELEMETRYFLOW_SERVICE_VERSION")

	if b.serviceName == "" {
		b.serviceName = "unknown-service"
	}
	if b.serviceVersion == "" {
		b.serviceVersion = "1.0.0"
	}

	return b
}

// WithEnvironment sets the deployment environment
func (b *Builder) WithEnvironment(env string) *Builder {
	b.environment = env
	return b
}

// WithEnvironmentFromEnv reads environment from ENV variable
func (b *Builder) WithEnvironmentFromEnv() *Builder {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
	}
	if env == "" {
		env = "production"
	}
	b.environment = env
	return b
}

// WithProtocol sets the OTLP protocol (grpc or http)
func (b *Builder) WithProtocol(protocol domain.Protocol) *Builder {
	b.protocol = protocol
	return b
}

// WithGRPC sets the protocol to gRPC
func (b *Builder) WithGRPC() *Builder {
	b.protocol = domain.ProtocolGRPC
	return b
}

// WithHTTP sets the protocol to HTTP
func (b *Builder) WithHTTP() *Builder {
	b.protocol = domain.ProtocolHTTP
	return b
}

// WithInsecure enables insecure connections (no TLS)
func (b *Builder) WithInsecure(insecure bool) *Builder {
	b.insecure = insecure
	return b
}

// WithTimeout sets connection timeout
func (b *Builder) WithTimeout(timeout time.Duration) *Builder {
	b.timeout = timeout
	return b
}

// WithSignals enables/disables specific signals
func (b *Builder) WithSignals(metrics, logs, traces bool) *Builder {
	b.enableMetrics = metrics
	b.enableLogs = logs
	b.enableTraces = traces
	return b
}

// WithMetricsOnly enables only metrics
func (b *Builder) WithMetricsOnly() *Builder {
	return b.WithSignals(true, false, false)
}

// WithLogsOnly enables only logs
func (b *Builder) WithLogsOnly() *Builder {
	return b.WithSignals(false, true, false)
}

// WithTracesOnly enables only traces
func (b *Builder) WithTracesOnly() *Builder {
	return b.WithSignals(false, false, true)
}

// WithCustomAttribute adds a custom resource attribute
func (b *Builder) WithCustomAttribute(key, value string) *Builder {
	b.customAttrs[key] = value
	return b
}

// WithAutoConfiguration attempts to configure from environment variables
func (b *Builder) WithAutoConfiguration() *Builder {
	return b.
		WithAPIKeyFromEnv().
		WithEndpointFromEnv().
		WithServiceFromEnv().
		WithEnvironmentFromEnv()
}

// Build creates the TelemetryFlow client
func (b *Builder) Build() (*Client, error) {
	// Check for errors collected during building
	if len(b.errors) > 0 {
		return nil, fmt.Errorf("builder errors: %v", b.errors)
	}

	// Validate required fields
	if b.apiKeyID == "" || b.apiKeySecret == "" {
		return nil, fmt.Errorf("API credentials are required")
	}
	if b.endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}
	if b.serviceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	// Create credentials
	credentials, err := domain.NewCredentials(b.apiKeyID, b.apiKeySecret)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %w", err)
	}

	// Create config
	config, err := domain.NewTelemetryConfig(credentials, b.endpoint, b.serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	// Apply builder settings
	config.
		WithProtocol(b.protocol).
		WithInsecure(b.insecure).
		WithTimeout(b.timeout).
		WithSignals(b.enableMetrics, b.enableLogs, b.enableTraces).
		WithServiceVersion(b.serviceVersion).
		WithEnvironment(b.environment)

	// Add custom attributes
	for key, value := range b.customAttrs {
		config.WithCustomAttribute(key, value)
	}

	// Create client
	return NewClient(config)
}

// MustBuild builds the client and panics on error (useful for quick setup)
func (b *Builder) MustBuild() *Client {
	client, err := b.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build TelemetryFlow client: %v", err))
	}
	return client
}

// ===== CONVENIENCE CONSTRUCTORS =====

// NewFromEnv creates a client from environment variables
func NewFromEnv() (*Client, error) {
	return NewBuilder().WithAutoConfiguration().Build()
}

// MustNewFromEnv creates a client from environment variables and panics on error
func MustNewFromEnv() *Client {
	return NewBuilder().WithAutoConfiguration().MustBuild()
}

// NewSimple creates a client with minimal configuration
func NewSimple(apiKeyID, apiKeySecret, endpoint, serviceName string) (*Client, error) {
	return NewBuilder().
		WithAPIKey(apiKeyID, apiKeySecret).
		WithEndpoint(endpoint).
		WithService(serviceName, "1.0.0").
		Build()
}

// MustNewSimple creates a client with minimal configuration and panics on error
func MustNewSimple(apiKeyID, apiKeySecret, endpoint, serviceName string) *Client {
	client, err := NewSimple(apiKeyID, apiKeySecret, endpoint, serviceName)
	if err != nil {
		panic(err)
	}
	return client
}
