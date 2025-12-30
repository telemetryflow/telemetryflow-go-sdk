<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

  <h3>TelemetryFlow GO SDK</h3>

[![Version](https://img.shields.io/badge/Version-1.1.1-orange.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![OTEL SDK](https://img.shields.io/badge/OpenTelemetry_SDK-1.39.0-blueviolet)](https://opentelemetry.io/)
[![OpenTelemetry](https://img.shields.io/badge/OTLP-100%25%20Compliant-success?logo=opentelemetry)](https://opentelemetry.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/telemetryflow/telemetryflow-sdk)

</div>

<p align="center">
  Enterprise-grade Go SDK for <a href="https://telemetryflow.id">TelemetryFlow</a> - the observability platform that provides unified metrics, logs, and traces collection following OpenTelemetry standards.
</p>

---

## 🌟 Features

- **100% OTLP Compliant**: Full OpenTelemetry Protocol implementation
- **DDD Architecture**: Domain-Driven Design with clear bounded contexts
- **CQRS Pattern**: Separate command and query responsibilities
- **Easy Integration**: Builder pattern and environment-based configuration
- **Exemplars Support**: Metrics-to-traces correlation for powerful debugging
- **Service Graph Ready**: Compatible with OTEL Collector servicegraph connector
- **Enhanced gRPC Settings**: Configurable keepalive, buffer sizes, and message limits
- **TelemetryFlow Headers**: Custom headers for collector authentication
- **Code Generators**: CLI tools for quick project setup
  - `telemetryflow-gen`: SDK integration generator
  - `telemetryflow-restapi`: DDD + CQRS RESTful API generator
- **Docker Ready**: Multi-stage Dockerfile and Docker Compose configurations
- **Production Ready**: Comprehensive error handling, retries, and batching
- **Type Safe**: Strong typing with compile-time safety
- **Zero Dependencies** (core): Minimal external dependencies in core SDK
- **E2E Testing**: Complete end-to-end testing infrastructure
- **Extensible**: Easy to extend for custom use cases

## 📦 Installation

```bash
go get github.com/telemetryflow/telemetryflow-go-sdk
```

## 🚀 Quick Start

### 1. Environment Variables Setup

Create a `.env` file:

```bash
TELEMETRYFLOW_API_KEY_ID=tfk_your_key_id_here
TELEMETRYFLOW_API_KEY_SECRET=tfs_your_secret_here
TELEMETRYFLOW_ENDPOINT=api.telemetryflow.id:4317
TELEMETRYFLOW_SERVICE_NAME=my-go-service
TELEMETRYFLOW_SERVICE_VERSION=1.1.1
TELEMETRYFLOW_SERVICE_NAMESPACE=telemetryflow
TELEMETRYFLOW_COLLECTOR_ID=my-collector-id
ENV=production
```

### 2. Initialize SDK

```go
package main

import (
    "context"
    "log"

    "github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

func main() {
    // Create client from environment variables
    client, err := telemetryflow.NewFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize the SDK
    ctx := context.Background()
    if err := client.Initialize(ctx); err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown(ctx)

    // Your application code...
}
```

### 3. Send Telemetry

```go
// Send a metric
client.IncrementCounter(ctx, "api.requests", 1, map[string]interface{}{
    "method": "GET",
    "status": 200,
})

// Send a log
client.LogInfo(ctx, "User logged in", map[string]interface{}{
    "user_id": "12345",
})

// Create a trace span
spanID, _ := client.StartSpan(ctx, "process-order", "internal", map[string]interface{}{
    "order_id": "67890",
})
defer client.EndSpan(ctx, spanID, nil)
```

## 🏗️ Architecture

The SDK follows Domain-Driven Design (DDD) and CQRS patterns:

```
pkg/telemetryflow/
├── domain/          # Domain layer (entities, value objects)
│   ├── credentials.go
│   └── config.go
├── application/     # Application layer (use cases, CQRS)
│   ├── commands.go
│   └── queries.go
├── infrastructure/  # Infrastructure layer (OTLP exporters)
│   ├── exporters.go
│   └── handlers.go
└── client.go        # Public API (interface layer)
```

### Domain Layer

Core business entities and value objects:
- `Credentials`: Immutable API key pair with validation
- `TelemetryConfig`: Aggregate root containing all configuration

### Application Layer

CQRS implementation:
- **Commands**: `RecordMetricCommand`, `EmitLogCommand`, `StartSpanCommand`, etc.
- **Queries**: `GetMetricQuery`, `GetLogsQuery`, `GetTraceQuery`, etc.
- **Command/Query Buses**: Route requests to appropriate handlers

### Infrastructure Layer

Technical implementations:
- `OTLPExporterFactory`: Creates OTLP exporters (gRPC/HTTP)
- `TelemetryCommandHandler`: Handles telemetry commands
- OpenTelemetry SDK integration

## 📚 Usage Examples

### Builder Pattern

```go
client := telemetryflow.NewBuilder().
    WithAPIKey("tfk_...", "tfs_...").
    WithEndpoint("api.telemetryflow.id:4317").
    WithService("my-service", "1.0.0").
    WithEnvironment("production").
    WithGRPC().
    WithSignals(true, true, true). // metrics, logs, traces
    WithCustomAttribute("team", "backend").
    MustBuild()
```

### Metrics

```go
// Counter
client.IncrementCounter(ctx, "requests.total", 1, map[string]interface{}{
    "method": "POST",
    "endpoint": "/api/users",
})

// Gauge
client.RecordGauge(ctx, "memory.usage", 512.0, map[string]interface{}{
    "unit": "MB",
})

// Histogram
client.RecordHistogram(ctx, "request.duration", 0.25, "s", map[string]interface{}{
    "endpoint": "/api/orders",
})
```

### Logs

```go
// Info level
client.LogInfo(ctx, "Application started", map[string]interface{}{
    "version": "1.0.0",
    "port": 8080,
})

// Warning level
client.LogWarn(ctx, "High memory usage", map[string]interface{}{
    "usage_mb": 512,
    "threshold_mb": 400,
})

// Error level
client.LogError(ctx, "Database connection failed", map[string]interface{}{
    "error": "timeout",
    "host": "db.example.com",
})
```

### Traces

```go
// Start a span
spanID, err := client.StartSpan(ctx, "process-payment", "internal", map[string]interface{}{
    "payment_method": "credit_card",
    "amount": 99.99,
})

// Add events to span
client.AddSpanEvent(ctx, spanID, "validation.complete", map[string]interface{}{
    "valid": true,
})

// End span (with optional error)
client.EndSpan(ctx, spanID, nil)
```

## 🛠️ Code Generators

The SDK includes powerful code generators to bootstrap your integration.

### SDK Generator (`telemetryflow-gen`)

Generates TelemetryFlow SDK integration code for existing projects.

```bash
# Install
go install github.com/telemetryflow/telemetryflow-go-sdk/cmd/generator@latest

# Initialize integration
telemetryflow-gen init \
    --project "my-app" \
    --service "my-service" \
    --key-id "tfk_..." \
    --key-secret "tfs_..."

# Generate examples
telemetryflow-gen example http-server
telemetryflow-gen example worker
```

### RESTful API Generator (`telemetryflow-restapi`)

Generates complete DDD + CQRS RESTful API projects with Echo framework.

```bash
# Install
go install github.com/telemetryflow/telemetryflow-go-sdk/cmd/generator-restfulapi@latest

# Create new project
telemetryflow-restapi new \
    --name order-service \
    --module github.com/myorg/order-service \
    --db-driver postgres

# Add entities
telemetryflow-restapi entity \
    --name Order \
    --fields 'customer_id:uuid,total:decimal,status:string'
```

See [Generator Documentation](docs/GENERATOR.md) and [RESTful API Generator](docs/GENERATOR_RESTAPI.md) for details.

## 🐳 Docker

### Using Docker

```bash
# Run SDK generator
docker run --rm -v $(pwd):/workspace telemetryflow/telemetryflow-sdk:latest \
    init --project myapp --output /workspace

# Run RESTful API generator
docker run --rm -v $(pwd):/workspace telemetryflow/telemetryflow-sdk:latest \
    telemetryflow-restapi new --name myapi --output /workspace
```

### Docker Compose

Development environment with full observability stack:

```bash
# Start development environment
docker-compose up -d

# Start with observability tools (Jaeger, Prometheus, Grafana)
docker-compose --profile full up -d

# Run E2E tests
docker-compose -f docker-compose.e2e.yml up --abort-on-container-exit
```

### Generated Structure

```
your-project/
├── telemetry/
│   ├── init.go           # SDK initialization
│   ├── metrics/
│   │   └── metrics.go    # Metrics helpers
│   ├── logs/
│   │   └── logs.go       # Logging helpers
│   ├── traces/
│   │   └── traces.go     # Tracing helpers
│   └── README.md
├── .env.telemetryflow    # Configuration template
└── example_*.go          # Generated examples
```

## 🎯 Advanced Configuration

### Protocol Selection

```go
// gRPC (default, recommended)
config.WithProtocol(domain.ProtocolGRPC)

// HTTP
config.WithProtocol(domain.ProtocolHTTP)
```

### Retry Configuration

```go
config.WithRetry(
    true,                    // enabled
    3,                       // max retries
    5 * time.Second,        // backoff duration
)
```

### Batch Settings

```go
config.WithBatchSettings(
    10 * time.Second,       // batch timeout
    512,                     // max batch size
)
```

### Signal Control

```go
// Enable specific signals only
config.WithSignals(
    true,   // metrics
    false,  // logs
    true,   // traces
)

// Or use convenience methods
config.WithMetricsOnly()
config.WithTracesOnly()
```

### Exemplars Support (v1.1.1+)

Exemplars enable metrics-to-traces correlation for powerful debugging:

```go
// Enabled by default, disable if not needed
client := telemetryflow.NewBuilder().
    WithAutoConfiguration().
    WithExemplars(false).  // Disable exemplars
    MustBuild()
```

### Service Namespace (v1.1.1+)

```go
// Set service namespace for multi-tenant environments
client := telemetryflow.NewBuilder().
    WithAutoConfiguration().
    WithServiceNamespace("my-namespace").
    MustBuild()
```

### Collector ID (v1.1.1+)

```go
// Set collector ID for TelemetryFlow backend authentication
client := telemetryflow.NewBuilder().
    WithAutoConfiguration().
    WithCollectorID("my-unique-collector-id").
    MustBuild()
```

### Custom Resource Attributes

```go
config.
    WithCustomAttribute("team", "backend").
    WithCustomAttribute("region", "us-east-1").
    WithCustomAttribute("datacenter", "dc1")
```

## 🔒 Security

### API Key Management

API keys should **never** be hardcoded. Use environment variables or secure secret management:

```go
// ✅ Good: From environment
client, _ := telemetryflow.NewFromEnv()

// ✅ Good: From secret manager (example)
apiKey := secretManager.GetSecret("telemetryflow/api-key")
client, _ := telemetryflow.NewBuilder().
    WithAPIKey(apiKey.KeyID, apiKey.KeySecret).
    Build()

// ❌ Bad: Hardcoded
client, _ := telemetryflow.NewBuilder().
    WithAPIKey("tfk_hardcoded", "tfs_secret"). // DON'T DO THIS
    Build()
```

### TLS/SSL

By default, the SDK uses secure connections. Only disable for testing:

```go
// Production (default)
config.WithInsecure(false)

// Testing only
config.WithInsecure(true)
```

## 📊 Best Practices

1. **Initialize Once**: Create one client instance per application
2. **Defer Shutdown**: Always defer `client.Shutdown(ctx)`
3. **Context Propagation**: Pass context through your call chain
4. **Attribute Cardinality**: Limit high-cardinality attributes
5. **Batch Configuration**: Tune batch settings for your workload
6. **Error Handling**: Always check errors from telemetry calls
7. **Graceful Shutdown**: Allow time for final flush

```go
func main() {
    client, _ := telemetryflow.NewFromEnv()
    ctx := context.Background()

    // Initialize
    if err := client.Initialize(ctx); err != nil {
        log.Fatal(err)
    }

    // Ensure graceful shutdown
    defer func() {
        shutdownCtx, cancel := context.WithTimeout(
            context.Background(),
            10*time.Second,
        )
        defer cancel()

        client.Flush(shutdownCtx)
        client.Shutdown(shutdownCtx)
    }()

    // Application code...
}
```

## 🧪 Testing

The SDK includes comprehensive tests organized by type:

```text
tests/
├── unit/              # Unit tests for individual components
│   ├── domain/        # Credentials, Config tests
│   ├── application/   # Command/Query tests
│   ├── infrastructure/# Template, HTTP, Database tests
│   └── client/        # Client, Builder tests
├── integration/       # Cross-layer integration tests
├── e2e/               # End-to-end pipeline tests
│   └── testdata/      # E2E test fixtures (OTEL config, nginx)
├── mocks/             # Mock implementations
└── fixtures/          # Test fixtures
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run unit tests only
go test ./tests/unit/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run short tests (skip integration)
go test -short ./...

# Run e2e tests (requires environment setup)
TELEMETRYFLOW_E2E=true go test ./tests/e2e/...
```

### Test Coverage Goals

| Layer          | Target |
| -------------- | ------ |
| Domain         | 90%+   |
| Application    | 85%+   |
| Infrastructure | 80%+   |
| Client         | 85%+   |

## 📖 Documentation

- [Quickstart Guide](docs/QUICKSTART.md) - Get started in 5 minutes
- [Architecture Guide](docs/ARCHITECTURE.md) - SDK architecture with DDD/CQRS patterns
- [API Reference](docs/API_REFERENCE.md) - Complete API documentation
- [SDK Generator](docs/GENERATOR.md) - telemetryflow-gen CLI documentation
- [RESTful API Generator](docs/GENERATOR_RESTAPI.md) - telemetryflow-restapi CLI documentation
- [Examples](examples/) - Sample applications and patterns
- [Changelog](CHANGELOG.md) - Version history and release notes
- [TelemetryFlow Platform](https://docs.telemetryflow.id) - Platform documentation

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- 📧 Email: support@telemetryflow.id
- 💬 Slack: [TelemetryFlow Community](https://telemetryflow.slack.com)
- 📚 Docs: [https://docs.telemetryflow.id](https://docs.telemetryflow.id)
- 🐛 Issues: [GitHub Issues](https://github.com/telemetryflow/telemetryflow-go-sdk/issues)

## 🙏 Acknowledgments

- [OpenTelemetry](https://opentelemetry.io/) for the excellent instrumentation standard
- All contributors who have helped shape this SDK

---

Built with ❤️ by the **DevOpsCorner Indonesia**
