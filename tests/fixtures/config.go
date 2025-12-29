// =============================================================================
// Test Fixtures - Configuration
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// =============================================================================

package fixtures

import "time"

// ConfigFixture represents a configuration test fixture
type ConfigFixture struct {
	Name        string
	ServiceName string
	Version     string
	Environment string
	Endpoint    string
	Protocol    string
	Insecure    bool
	Timeout     time.Duration
	BatchSize   int
	Metrics     bool
	Logs        bool
	Traces      bool
	Attributes  map[string]string
	Valid       bool
	ErrorMsg    string
}

// ValidConfigs contains valid configuration scenarios
var ValidConfigs = []ConfigFixture{
	{
		Name:        "minimal_config",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Insecure:    true,
		Timeout:     30 * time.Second,
		BatchSize:   512,
		Metrics:     true,
		Logs:        true,
		Traces:      true,
		Attributes:  nil,
		Valid:       true,
	},
	{
		Name:        "production_config",
		ServiceName: "production-api",
		Version:     "1.1.1",
		Environment: "production",
		Endpoint:    "api.telemetryflow.id:4317",
		Protocol:    "grpc",
		Insecure:    false,
		Timeout:     60 * time.Second,
		BatchSize:   1024,
		Metrics:     true,
		Logs:        true,
		Traces:      true,
		Attributes: map[string]string{
			"team":       "backend",
			"region":     "us-east-1",
			"datacenter": "dc1",
		},
		Valid: true,
	},
	{
		Name:        "http_protocol",
		ServiceName: "http-service",
		Version:     "1.1.1",
		Environment: "staging",
		Endpoint:    "https://api.telemetryflow.id:4318",
		Protocol:    "http",
		Insecure:    false,
		Timeout:     30 * time.Second,
		BatchSize:   256,
		Metrics:     true,
		Logs:        true,
		Traces:      true,
		Attributes:  nil,
		Valid:       true,
	},
	{
		Name:        "metrics_only",
		ServiceName: "metrics-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Insecure:    true,
		Timeout:     30 * time.Second,
		BatchSize:   512,
		Metrics:     true,
		Logs:        false,
		Traces:      false,
		Attributes:  nil,
		Valid:       true,
	},
	{
		Name:        "traces_only",
		ServiceName: "tracing-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Insecure:    true,
		Timeout:     30 * time.Second,
		BatchSize:   512,
		Metrics:     false,
		Logs:        false,
		Traces:      true,
		Attributes:  nil,
		Valid:       true,
	},
}

// InvalidConfigs contains invalid configuration scenarios
var InvalidConfigs = []ConfigFixture{
	{
		Name:        "empty_service_name",
		ServiceName: "",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Valid:       false,
		ErrorMsg:    "service name cannot be empty",
	},
	{
		Name:        "empty_endpoint",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "",
		Protocol:    "grpc",
		Valid:       false,
		ErrorMsg:    "endpoint cannot be empty",
	},
	{
		Name:        "invalid_protocol",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "websocket",
		Valid:       false,
		ErrorMsg:    "protocol must be 'grpc' or 'http'",
	},
	{
		Name:        "negative_timeout",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Timeout:     -1 * time.Second,
		Valid:       false,
		ErrorMsg:    "timeout must be positive",
	},
	{
		Name:        "zero_batch_size",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		BatchSize:   0,
		Valid:       false,
		ErrorMsg:    "batch size must be positive",
	},
	{
		Name:        "no_signals_enabled",
		ServiceName: "test-service",
		Version:     "1.1.1",
		Environment: "development",
		Endpoint:    "localhost:4317",
		Protocol:    "grpc",
		Metrics:     false,
		Logs:        false,
		Traces:      false,
		Valid:       false,
		ErrorMsg:    "at least one signal must be enabled",
	},
}

// EnvConfigFile represents a .env configuration file content
var EnvConfigFile = `# TelemetryFlow Configuration
TELEMETRYFLOW_API_KEY_ID=tfk_test_key_12345
TELEMETRYFLOW_API_KEY_SECRET=tfs_test_secret_67890
TELEMETRYFLOW_ENDPOINT=api.telemetryflow.id:4317
TELEMETRYFLOW_SERVICE_NAME=test-service
TELEMETRYFLOW_SERVICE_VERSION=1.1.1
TELEMETRYFLOW_ENVIRONMENT=development
TELEMETRYFLOW_PROTOCOL=grpc
TELEMETRYFLOW_INSECURE=false
TELEMETRYFLOW_TIMEOUT=30s
TELEMETRYFLOW_BATCH_SIZE=512
TELEMETRYFLOW_ENABLE_METRICS=true
TELEMETRYFLOW_ENABLE_LOGS=true
TELEMETRYFLOW_ENABLE_TRACES=true
`

// YAMLConfigFile represents a YAML configuration file content
var YAMLConfigFile = `# TelemetryFlow Configuration
telemetryflow:
  api_key_id: tfk_test_key_12345
  api_key_secret: tfs_test_secret_67890
  endpoint: api.telemetryflow.id:4317
  service:
    name: test-service
    Version: 1.1.1
  environment: development
  protocol: grpc
  insecure: false
  timeout: 30s
  batch_size: 512
  signals:
    metrics: true
    logs: true
    traces: true
  attributes:
    team: backend
    region: us-east-1
`
