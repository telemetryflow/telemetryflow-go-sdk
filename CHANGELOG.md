<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="docs/assets/tfo-logo-sdk-light.svg">
    <img alt="TelemetryFlow SDK Logo" src="docs/assets/tfo-logo-sdk-light.svg" width="80%">
  </picture>
</p>

<p align="center">
  <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go" alt="Go Version"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-Apache%202.0-blue.svg" alt="License"></a>
  <a href="https://opentelemetry.io/"><img src="https://img.shields.io/badge/OTLP-100%25%20Compliant-green" alt="OTLP Compliant"></a>
  <a href="https://hub.docker.com/r/telemetryflow/telemetryflow-sdk"><img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker" alt="Docker"></a>
  <a href="CHANGELOG.md"><img src="https://img.shields.io/badge/Version-1.0.1-blue.svg" alt="Version"></a>
</p>

<p align="center">
  Enterprise-grade Go SDK for <a href="https://telemetryflow.id">TelemetryFlow</a> - the observability platform that provides unified metrics, logs, and traces collection following OpenTelemetry standards.
</p>

---

# Changelog

All notable changes to the TelemetryFlow Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2024-12-23

### Added

#### Docker Support
- Multi-stage `Dockerfile` for building SDK generators with minimal image size
- `docker-compose.yml` for development environment with:
  - OpenTelemetry Collector
  - Jaeger (tracing)
  - Prometheus (metrics)
  - Grafana (visualization)
- `docker-compose.e2e.yml` for end-to-end testing infrastructure with:
  - OTEL Collector for telemetry
  - Mock backend (nginx)
  - Jaeger for trace visualization
  - PostgreSQL for integration tests
  - Redis for caching tests

#### Documentation Enhancements
- Mermaid diagrams added to all documentation files:
  - `ARCHITECTURE.md`: SDK layer diagram, telemetry flow, and component relationships
  - `API_REFERENCE.md`: Package overview diagram
  - `QUICKSTART.md`: SDK initialization sequence diagram
  - `GENERATOR.md`: Command flow and telemetry integration diagrams
  - `GENERATOR_RESTAPI.md`: DDD architecture and CQRS pattern diagrams
- New documentation file: `GENERATOR_RESTAPI.md` - comprehensive RESTful API generator docs

#### Testing Infrastructure
- Unit tests for infrastructure layer:
  - `template_test.go`: Template functions and execution tests
  - `config_test.go`: Configuration file handling tests
  - `http_test.go`: HTTP middleware and response pattern tests
  - `database_test.go`: Database configuration and query pattern tests
- E2E test data files:
  - `otel-collector.yaml`: OpenTelemetry Collector configuration
  - `nginx.conf`: Mock backend configuration
  - `mock-responses/health.json`: Health check response fixture
- Comprehensive test fixtures package (`tests/fixtures/`):
  - `credentials.go`: Valid/invalid API key scenarios
  - `config.go`: SDK configuration test data
  - `telemetry.go`: Sample metrics, logs, and traces fixtures
  - `templates.go`: Code generator template test data
  - `responses.go`: Mock HTTP/OTLP response fixtures
  - `fixtures.go`: Helper functions for loading fixtures

#### SDK Logo

- Added TelemetryFlow SDK logos with dark/light theme support:
  - `docs/assets/tfo-logo-sdk-dark.svg`
  - `docs/assets/tfo-logo-sdk-light.svg`

#### Configuration Files
- `configs/otel-collector.yaml`: Production-ready OTEL Collector configuration
- `configs/prometheus.yaml`: Prometheus scrape configuration

### Changed

- Updated Go version requirement from 1.21 to 1.24 in README badges
- Expanded README features list to include:
  - Docker support
  - RESTful API generator
  - E2E testing infrastructure
- Enhanced README with Docker and Docker Compose usage sections
- Updated documentation links in README

### Fixed

- Improved ASCII diagrams converted to Mermaid for better rendering on GitHub

## [1.0.0] - 2024-12-01

### Added

#### Core SDK
- `telemetryflow.Client`: Main SDK client for telemetry operations
- `telemetryflow.NewFromEnv()`: Environment-based client initialization
- `telemetryflow.NewBuilder()`: Builder pattern for client configuration

#### Domain Layer (DDD)
- `Credentials`: Immutable API key pair value object with validation
- `TelemetryConfig`: Aggregate root for SDK configuration
- Support for gRPC and HTTP OTLP protocols

#### Application Layer (CQRS)
- Commands: `RecordMetricCommand`, `EmitLogCommand`, `StartSpanCommand`, `EndSpanCommand`
- Queries: `GetMetricQuery`, `GetLogsQuery`, `GetTraceQuery`
- Command and Query bus implementations

#### Infrastructure Layer
- `OTLPExporterFactory`: Creates OTLP exporters for metrics, logs, and traces
- `TelemetryCommandHandler`: Handles telemetry commands
- OpenTelemetry SDK integration with automatic batching

#### Metrics
- `IncrementCounter()`: Counter metric recording
- `RecordGauge()`: Gauge metric recording
- `RecordHistogram()`: Histogram metric recording with units

#### Logs
- `LogInfo()`, `LogWarn()`, `LogError()`, `LogDebug()`: Structured logging at various levels
- Attribute support for log enrichment

#### Traces
- `StartSpan()`: Create new spans with attributes
- `EndSpan()`: Complete spans with optional error status
- `AddSpanEvent()`: Add events to active spans
- Context propagation support

#### Code Generator (`telemetryflow-gen`)
- `init`: Initialize TelemetryFlow integration in projects
- `example`: Generate example code (http-server, worker, grpc-server, basic)
- `config`: Generate configuration files
- `version`: Display version information
- Custom template support with `--template-dir`

#### RESTful API Generator (`telemetryflow-restapi`)
- `new`: Create complete DDD + CQRS RESTful API projects
- `entity`: Generate entities with repository, handlers, and migrations
- `docs`: Regenerate API documentation
- Echo framework integration
- OpenAPI/Swagger documentation generation
- Postman collection generation
- JWT authentication middleware
- CORS and rate limiting middleware
- PostgreSQL, MySQL, and SQLite support

#### Configuration
- Environment variable support (`TELEMETRYFLOW_*`)
- `.env.telemetryflow` configuration file support
- Protocol selection (gRPC/HTTP)
- Retry configuration with exponential backoff
- Batch settings customization
- Signal control (metrics, logs, traces)
- Custom resource attributes

### Documentation
- `ARCHITECTURE.md`: SDK architecture with DDD/CQRS patterns
- `API_REFERENCE.md`: Complete API documentation
- `QUICKSTART.md`: Getting started guide
- `GENERATOR.md`: SDK generator documentation
- Example applications in `examples/` directory

---

## Release Notes

### Upgrade Guide

#### From 0.x to 1.0.0

The 1.0.0 release is a complete rewrite with breaking changes:

1. **Client Initialization**:
   ```go
   // Old (0.x)
   client := telemetryflow.NewClient(config)

   // New (1.0.0+)
   client, err := telemetryflow.NewFromEnv()
   // or
   client := telemetryflow.NewBuilder().
       WithAPIKey("key", "secret").
       MustBuild()
   ```

2. **Metric Recording**:
   ```go
   // Old (0.x)
   client.RecordMetric("name", value, tags)

   // New (1.0.0+)
   client.IncrementCounter(ctx, "name", value, attributes)
   client.RecordHistogram(ctx, "name", value, "unit", attributes)
   ```

3. **Logging**:
   ```go
   // Old (0.x)
   client.Log(level, message, fields)

   // New (1.0.0+)
   client.LogInfo(ctx, message, attributes)
   client.LogError(ctx, message, attributes)
   ```

4. **Tracing**:
   ```go
   // Old (0.x)
   span := client.StartTrace("name")
   defer span.End()

   // New (1.0.0+)
   spanID, _ := client.StartSpan(ctx, "name", "kind", attrs)
   defer client.EndSpan(ctx, spanID, nil)
   ```

### Known Issues

None at this time.

### Contributors

- DevOpsCorner Indonesia Team

---

[1.0.1]: https://github.com/telemetryflow/telemetryflow-go-sdk/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/telemetryflow/telemetryflow-go-sdk/releases/tag/v1.0.0
