// Package telemetryflow provides the main SDK interface for TelemetryFlow.
// This package offers a fluent builder pattern for creating telemetry clients
// that export traces, metrics, and logs to TelemetryFlow collectors using OTLP.
package telemetryflow

import (
	"fmt"
	"os"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

// Builder provides a fluent interface for creating TelemetryFlow clients
type Builder struct {
	apiKeyID         string
	apiKeySecret     string
	collectorID      string
	endpoint         string
	serviceName      string
	serviceNamespace string
	serviceVersion   string
	environment      string
	datacenter       string
	protocol         domain.Protocol
	insecure         bool
	timeout          time.Duration
	enableMetrics    bool
	enableLogs       bool
	enableTraces     bool
	enableExemplars  bool
	customAttrs      map[string]string
	errors           []error

	// TFO v2 API settings (aligned with tfoexporter)
	useV2API        bool
	v2Only          bool
	tracesEndpoint  string
	metricsEndpoint string
	logsEndpoint    string

	// Collector Identity (aligned with tfoidentityextension)
	collectorName        string
	collectorDescription string
	collectorHostname    string
	collectorTags        map[string]string
	enrichResources      bool
}

// NewBuilder creates a new SDK builder
func NewBuilder() *Builder {
	return &Builder{
		protocol:         domain.ProtocolGRPC,
		insecure:         false,
		timeout:          30 * time.Second,
		enableMetrics:    true,
		enableLogs:       true,
		enableTraces:     true,
		enableExemplars:  true, // enabled by default for metrics-to-traces correlation
		serviceNamespace: "telemetryflow",
		datacenter:       "default",
		customAttrs:      make(map[string]string),
		errors:           make([]error, 0),
		// TFO v2 API settings (aligned with tfoexporter)
		useV2API:        true,  // v2 API enabled by default for TFO Platform
		v2Only:          false, // allow both v1 and v2 by default
		tracesEndpoint:  "",    // use defaults
		metricsEndpoint: "",    // use defaults
		logsEndpoint:    "",    // use defaults
		// Collector Identity (aligned with tfoidentityextension)
		collectorName:        "",
		collectorDescription: "",
		collectorHostname:    "",
		collectorTags:        make(map[string]string),
		enrichResources:      true, // enabled by default
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

// WithCollectorID sets the collector identifier for TelemetryFlow headers
func (b *Builder) WithCollectorID(id string) *Builder {
	b.collectorID = id
	return b
}

// WithCollectorIDFromEnv reads collector ID from environment variable
func (b *Builder) WithCollectorIDFromEnv() *Builder {
	b.collectorID = os.Getenv("TELEMETRYFLOW_COLLECTOR_ID")
	return b
}

// WithServiceNamespace sets the service namespace
func (b *Builder) WithServiceNamespace(namespace string) *Builder {
	b.serviceNamespace = namespace
	return b
}

// WithServiceNamespaceFromEnv reads service namespace from environment variable
func (b *Builder) WithServiceNamespaceFromEnv() *Builder {
	namespace := os.Getenv("TELEMETRYFLOW_SERVICE_NAMESPACE")
	if namespace == "" {
		namespace = "telemetryflow"
	}
	b.serviceNamespace = namespace
	return b
}

// WithExemplars enables/disables exemplars for metrics-to-traces correlation
func (b *Builder) WithExemplars(enabled bool) *Builder {
	b.enableExemplars = enabled
	return b
}

// WithV2API enables/disables TFO Platform v2 API endpoints (aligned with tfoexporter)
func (b *Builder) WithV2API(enabled bool) *Builder {
	b.useV2API = enabled
	return b
}

// WithV2Only sets v2-only mode - only v2 endpoints will be used
func (b *Builder) WithV2Only() *Builder {
	b.v2Only = true
	b.useV2API = true // v2Only implies useV2API
	return b
}

// WithTracesEndpoint sets a custom traces endpoint path
func (b *Builder) WithTracesEndpoint(endpoint string) *Builder {
	b.tracesEndpoint = endpoint
	return b
}

// WithMetricsEndpoint sets a custom metrics endpoint path
func (b *Builder) WithMetricsEndpoint(endpoint string) *Builder {
	b.metricsEndpoint = endpoint
	return b
}

// WithLogsEndpoint sets a custom logs endpoint path
func (b *Builder) WithLogsEndpoint(endpoint string) *Builder {
	b.logsEndpoint = endpoint
	return b
}

// WithCollectorName sets the human-readable collector name (aligned with tfoidentityextension)
func (b *Builder) WithCollectorName(name string) *Builder {
	b.collectorName = name
	return b
}

// WithCollectorNameFromEnv reads collector name from environment variable
func (b *Builder) WithCollectorNameFromEnv() *Builder {
	b.collectorName = os.Getenv("TELEMETRYFLOW_COLLECTOR_NAME")
	return b
}

// WithCollectorDescription sets the collector description (aligned with tfoidentityextension)
func (b *Builder) WithCollectorDescription(description string) *Builder {
	b.collectorDescription = description
	return b
}

// WithCollectorHostname sets the collector hostname (aligned with tfoidentityextension)
func (b *Builder) WithCollectorHostname(hostname string) *Builder {
	b.collectorHostname = hostname
	return b
}

// WithCollectorTag adds a single collector tag (aligned with tfoidentityextension)
func (b *Builder) WithCollectorTag(key, value string) *Builder {
	b.collectorTags[key] = value
	return b
}

// WithCollectorTags sets all collector tags (aligned with tfoidentityextension)
func (b *Builder) WithCollectorTags(tags map[string]string) *Builder {
	b.collectorTags = tags
	return b
}

// WithEnrichResources enables/disables collector identity enrichment (aligned with tfoidentityextension)
func (b *Builder) WithEnrichResources(enabled bool) *Builder {
	b.enrichResources = enabled
	return b
}

// WithDatacenter sets the datacenter identifier
func (b *Builder) WithDatacenter(datacenter string) *Builder {
	b.datacenter = datacenter
	return b
}

// WithDatacenterFromEnv reads datacenter from environment variable
func (b *Builder) WithDatacenterFromEnv() *Builder {
	dc := os.Getenv("TELEMETRYFLOW_DATACENTER")
	if dc == "" {
		dc = "default"
	}
	b.datacenter = dc
	return b
}

// WithAutoConfiguration attempts to configure from environment variables
func (b *Builder) WithAutoConfiguration() *Builder {
	return b.
		WithAPIKeyFromEnv().
		WithEndpointFromEnv().
		WithServiceFromEnv().
		WithServiceNamespaceFromEnv().
		WithCollectorIDFromEnv().
		WithCollectorNameFromEnv().
		WithDatacenterFromEnv().
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
		WithServiceNamespace(b.serviceNamespace).
		WithEnvironment(b.environment).
		WithDatacenter(b.datacenter).
		WithExemplars(b.enableExemplars).
		WithV2API(b.useV2API).
		WithV2Only(b.v2Only).
		WithEnrichResources(b.enrichResources)

	// Set collector ID if provided
	if b.collectorID != "" {
		config.WithCollectorID(b.collectorID)
	}

	// Set collector identity fields (aligned with tfoidentityextension)
	if b.collectorName != "" {
		config.WithCollectorName(b.collectorName)
	}
	if b.collectorDescription != "" {
		config.WithCollectorDescription(b.collectorDescription)
	}
	if b.collectorHostname != "" {
		config.WithCollectorHostname(b.collectorHostname)
	}

	// Set collector tags
	for key, value := range b.collectorTags {
		config.WithCollectorTag(key, value)
	}

	// Set custom endpoint paths (aligned with tfoexporter)
	if b.tracesEndpoint != "" {
		config.WithTracesEndpoint(b.tracesEndpoint)
	}
	if b.metricsEndpoint != "" {
		config.WithMetricsEndpoint(b.metricsEndpoint)
	}
	if b.logsEndpoint != "" {
		config.WithLogsEndpoint(b.logsEndpoint)
	}

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
