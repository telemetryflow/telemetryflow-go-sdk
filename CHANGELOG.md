<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

  <h3>TelemetryFlow GO SDK</h3>

[![Version](https://img.shields.io/badge/Version-1.1.2-orange.svg)](CHANGELOG.md)
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

# Changelog

All notable changes to the TelemetryFlow Go SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.2] - 2026-01-04

### Added

#### TFO v2 API Support (aligned with tfoexporter)

- **V2 Endpoint Configuration**: New methods for TFO v2 API endpoints
  - `WithV2API(bool)` - Enable/disable v2 API (default: true)
  - `WithV2Only()` - Enable v2-only mode (disables v1 endpoints)
  - `WithTracesEndpoint(string)` - Custom traces endpoint path
  - `WithMetricsEndpoint(string)` - Custom metrics endpoint path
  - `WithLogsEndpoint(string)` - Custom logs endpoint path
  - Default endpoints: `/v2/traces`, `/v2/metrics`, `/v2/logs`

- **Config Getters**: New configuration accessors
  - `UseV2API()` - Check if v2 API is enabled
  - `IsV2Only()` - Check if v2-only mode is enabled
  - `TracesEndpoint()` - Get traces endpoint path
  - `MetricsEndpoint()` - Get metrics endpoint path
  - `LogsEndpoint()` - Get logs endpoint path

#### Collector Identity (aligned with tfoidentityextension)

- **Identity Configuration**: New methods for collector identity
  - `WithCollectorName(string)` - Set collector name
  - `WithCollectorNameFromEnv()` - Read from `TELEMETRYFLOW_COLLECTOR_NAME`
  - `WithCollectorDescription(string)` - Set collector description
  - `WithCollectorHostname(string)` - Set collector hostname
  - `WithCollectorTag(key, value)` - Add single tag
  - `WithCollectorTags(map[string]string)` - Set all tags
  - `WithEnrichResources(bool)` - Enable/disable resource enrichment
  - `WithDatacenter(string)` - Set datacenter identifier
  - `WithDatacenterFromEnv()` - Read from `TELEMETRYFLOW_DATACENTER`

- **Identity Getters**: New configuration accessors
  - `CollectorName()`, `CollectorDescription()`, `CollectorHostname()`
  - `CollectorTags()`, `IsEnrichResourcesEnabled()`, `Datacenter()`

#### Enhanced Headers (aligned with tfoauthextension)

- **HTTP Headers**: New TelemetryFlow-specific headers
  - `X-TelemetryFlow-Collector-Name` - Collector name
  - `X-TelemetryFlow-Collector-Hostname` - Collector hostname
  - `X-TelemetryFlow-Environment` - Deployment environment
  - `X-TelemetryFlow-Datacenter` - Datacenter identifier
  - `X-TelemetryFlow-API-Version` - API version indicator (v2)

- **gRPC Metadata**: Same headers added to gRPC interceptor
  - All headers in lowercase for gRPC metadata compliance

#### TFO-Collector v1.1.2 Integration

- **Unified TFO Configuration**: Added `configs/tfo-collector-unified.yaml` for full TFO v2-only mode
  - v2 endpoints only with mandatory authentication (`v2_only: true`)
  - TFO exporter enabled with dedicated instances for traces, metrics, and logs
  - Auto-injected authentication via `tfoauth` extension
  - Collector identity enrichment via `tfoidentity` extension

- **TFO Custom Components Support**: Updated configurations for TFO-Collector OCB-native components
  - `tfootlp` receiver - OTLP with dual v1/v2 endpoint support
  - `tfo` exporter - Auto TFO authentication header injection
  - `tfoauth` extension - Centralized API key management
  - `tfoidentity` extension - Collector identity and resource enrichment

- **Enhanced TFO Exporter Configuration**: Multiple dedicated TFO exporter instances
  - `tfo/traces` - Trace-specific settings with queue size 2000
  - `tfo/metrics` - Metrics-specific settings with queue size 1000
  - `tfo/logs` - Logs-specific settings with queue size 1500
  - All with gzip compression, retry policies, and sending queues

#### SDK Configuration Files

- **Default SDK Configuration**: Added `configs/sdk-default.yaml`
  - Full configuration with all TFO v2 API settings
  - Environment variable substitution support
  - Collector identity, batching, retry, and gRPC settings
  - Production-ready defaults with sensible values

- **V2-Only SDK Configuration**: Added `configs/sdk-v2-only.yaml`
  - TFO v2-only mode with mandatory authentication
  - Production-optimized settings (compression, larger batches)
  - Enhanced gRPC buffer and message sizes
  - Complete collector identity configuration

- **Minimal SDK Configuration**: Added `configs/sdk-minimal.yaml`
  - Quick-start configuration with essential settings only
  - Service name, credentials, and endpoint
  - All signals enabled with TFO v2 API

- **Updated .env.example**: Enhanced with TFO v2 API environment variables
  - TFO v2 API settings: `TELEMETRYFLOW_USE_V2_API`, `TELEMETRYFLOW_V2_ONLY`
  - Collector identity: `TELEMETRYFLOW_COLLECTOR_NAME`, `TELEMETRYFLOW_DATACENTER`
  - Protocol and retry settings
  - Batch and signal configuration

### Changed

- **HTTP Exporters**: Updated to use configurable URL paths
  - `otlptracehttp.WithURLPath(config.TracesEndpoint())`
  - `otlpmetrichttp.WithURLPath(config.MetricsEndpoint())`

- **Builder Defaults**: Updated default values
  - `useV2API: true` - v2 API enabled by default
  - `enrichResources: true` - Resource enrichment enabled by default
  - `datacenter: "default"` - Default datacenter value

- **AutoConfiguration**: Extended to include new environment variables
  - `WithCollectorNameFromEnv()` added to auto-configuration chain
  - `WithDatacenterFromEnv()` added to auto-configuration chain

- **Docker Compose**: Updated for TFO-Collector v1.1.2 OCB-native architecture
  - Unified binary name: `tfo-collector` (no more `-ocb` suffix)
  - Updated configuration mounts for new TFO config files
  - Short flag support: `-c` for config, `-s` for set, `-f` for feature-gates

- **OTEL Collector Configs**: Updated to align with OTEL 0.142.0
  - Updated `configs/otel-collector.yaml` with latest best practices
  - Added `configs/otel-collector-minimal.yaml` for quick start scenarios

### Tests

- **Domain Config Tests**: Added TFO v2 API test cases
  - `TestTelemetryConfig_V2APIDefaults`
  - `TestTelemetryConfig_WithV2API`
  - `TestTelemetryConfig_WithV2Only`
  - `TestTelemetryConfig_WithCustomEndpoints`
  - `TestTelemetryConfig_CollectorIdentityDefaults`
  - `TestTelemetryConfig_WithCollectorName`
  - `TestTelemetryConfig_WithCollectorDescription`
  - `TestTelemetryConfig_WithCollectorHostname`
  - `TestTelemetryConfig_WithCollectorTags`
  - `TestTelemetryConfig_WithEnrichResources`
  - `TestTelemetryConfig_WithDatacenter`
  - `TestTelemetryConfig_TFOv2FullChaining`

- **Builder Tests**: Added TFO v2 API builder test cases
  - `TestBuilder_WithV2API`
  - `TestBuilder_WithV2Only`
  - `TestBuilder_WithCustomEndpoints`
  - `TestBuilder_WithCollectorName`
  - `TestBuilder_WithCollectorNameFromEnv`
  - `TestBuilder_WithCollectorDescription`
  - `TestBuilder_WithCollectorHostname`
  - `TestBuilder_WithCollectorTags`
  - `TestBuilder_WithEnrichResources`
  - `TestBuilder_WithDatacenter`
  - `TestBuilder_WithDatacenterFromEnv`
  - `TestBuilder_TFOv2FullChaining`

- **Benchmarks**: Added TFO v2 benchmark
  - `BenchmarkBuilder_BuildWithTFOv2`

### Documentation

- Updated configuration documentation for TFO-Collector v1.1.2
- Added TFO v2 endpoint usage examples with authentication headers
- Added TFO custom components reference table

### Compatibility

- TFO-Collector: v1.1.2 (OCB-native)
- OpenTelemetry Collector: v0.142.0
- OpenTelemetry Go SDK: v1.39.0

---

## [1.1.1] - 2024-12-30

### Added

- **Dual Endpoint Ingestion Support**: Updated docker-compose and OTEL collector configs for TFO-Collector dual ingestion
  - v1 endpoints: Standard OTEL community format (`/v1/traces`, `/v1/metrics`, `/v1/logs`)
  - v2 endpoints: TelemetryFlow enhanced format (`/v2/traces`, `/v2/metrics`, `/v2/logs`)
  - gRPC endpoint: Same port (4317) for both v1 and v2
- **TFO-Collector as Default**: Docker-compose now uses `telemetryflow/telemetryflow-collector` as default image
  - Commented alternatives for TFO-Collector-OCB and OTEL Collector Contrib
  - Separate volume mounts for each collector type
- **Enhanced Port Configuration**: Added additional ports for observability
  - zPages (55679) for debugging
  - pprof (1777) for profiling
  - Prometheus exporter (8889)
- **Connectors for Exemplars**: Added spanmetrics and servicegraph connectors
  - Metrics-to-traces correlation with exemplars enabled
  - Service dependency graph generation
- **Comprehensive Unit Tests**: Added unit tests for all DDD layers following external test package pattern
  - `tests/unit/domain/` - Credentials and Config domain tests
  - `tests/unit/application/` - Command and Query tests
  - `tests/unit/infrastructure/` - Template, HTTP, Database tests
  - `tests/unit/client/` - Client and Builder tests
  - `tests/unit/generator/` - Generator and REST API generator tests
  - `tests/unit/version/` - Version information tests
  - `tests/unit/banner/` - Banner display tests

### Changed

- **Documentation**: Updated docs structure to match TelemetryFlow Collector format
  - Added `docs/README.md` - Documentation index
  - Added `docs/TESTING.md` - Comprehensive testing guide
  - Added `docs/BUILD-SYSTEM.md` - Build system documentation
  - Updated `docs/ARCHITECTURE.md` - Added version header and DDD layer details

### Fixed

- **Linting**: Fixed ST1000 package comment errors for test packages
- **Test Utils**: Changed from `os.Setenv`/`os.Unsetenv` to `t.Setenv` for proper test cleanup

---

## [1.1.0] - 2024-12-27

### Added

#### Go SDK Enhancements

- **Exemplars Support**: Enabled by default for metrics-to-traces correlation
  - New `WithExemplars(bool)` builder method
  - New `IsExemplarsEnabled()` config getter
  - Compatible with OTEL Collector spanmetrics connector

- **Service Namespace**: Support for multi-tenant service organization
  - New `WithServiceNamespace(string)` builder and config method
  - New `WithServiceNamespaceFromEnv()` builder method
  - New `ServiceNamespace()` config getter
  - Environment variable: `TELEMETRYFLOW_SERVICE_NAMESPACE`
  - Default value: `telemetryflow`
  - Added to OTLP resource using semconv `service.namespace`

- **Collector ID**: Unique identifier for TelemetryFlow collector authentication
  - New `WithCollectorID(string)` builder and config method
  - New `WithCollectorIDFromEnv()` builder method
  - New `CollectorID()` config getter
  - Environment variable: `TELEMETRYFLOW_COLLECTOR_ID`
  - Added `X-TelemetryFlow-Collector-ID` header to OTLP exports

- **TelemetryFlow Custom Headers**: Enhanced authentication headers
  - `X-TelemetryFlow-Key-ID`: API Key ID header
  - `X-TelemetryFlow-Key-Secret`: API Key Secret header
  - `X-TelemetryFlow-Collector-ID`: Collector identifier header
  - Headers added to both gRPC metadata and HTTP headers

- **Enhanced gRPC Configuration**:
  - New `GRPCKeepaliveConfig` struct with Time, Timeout, PermitWithoutStream
  - `WithGRPCKeepalive(time, timeout, permitWithoutStream)` config method
  - `WithGRPCBufferSizes(readSize, writeSize)` config method
  - `WithGRPCMessageSizes(recvSize, sendSize)` config method
  - New getters: `GRPCKeepalive()`, `GRPCMaxRecvMsgSize()`, `GRPCMaxSendMsgSize()`, `GRPCReadBufferSize()`, `GRPCWriteBufferSize()`
  - Default settings aligned with OTEL Collector configuration

#### OTEL Collector Configuration Enhancements

- **Connectors Support**: Added `spanmetrics` and `servicegraph` connectors for:
  - Deriving metrics from traces with exemplars support
  - Building service dependency graphs automatically
  - Metrics-to-traces correlation for drill-down analysis
- **Environment Variable Support**: Full support for `TELEMETRYFLOW_*` environment variables in configuration templates
- **TelemetryFlow Extensions**: Custom `telemetryflow:` and `collector:` sections for standalone build authentication

#### Configuration Templates

- New `otel-collector.yaml.tpl` template replacing `ocb-collector.yaml.tpl`
- Updated `tfo-collector.yaml.tpl` with standard OTEL format (removed custom `enabled` flags)
- Added connector pipelines: `metrics/spanmetrics` and `metrics/servicegraph`
- OpenMetrics support with `enable_open_metrics: true` for exemplars in Prometheus

### Changed

#### SDK Default Configuration

- **Default Configuration**: Updated defaults to align with TelemetryFlow Collector
  - gRPC keepalive: 10s time, 5s timeout
  - gRPC buffer sizes: 512KB read/write
  - gRPC message sizes: 4 MiB max recv/send
  - Exemplars enabled by default
  - Service namespace defaults to "telemetryflow"

- **Resource Attributes**: Service namespace now added to OTLP resource
  - Uses OpenTelemetry semconv `service.namespace` attribute

- **AutoConfiguration**: Extended to include new environment variables
  - Now reads `TELEMETRYFLOW_SERVICE_NAMESPACE`
  - Now reads `TELEMETRYFLOW_COLLECTOR_ID`

#### Breaking Changes

- **Configuration Format**: Migrated from custom format with `enabled` flags to standard OTEL format
  - Old: `enabled: true/false` flags throughout config
  - New: Comment out sections to disable (standard OTEL approach)
- **Template Renamed**: `ocb-collector.yaml.tpl` → `otel-collector.yaml.tpl`
- **Exporter Renamed**: `logging` exporter → `debug` exporter (standard OTEL naming)
- **Pipeline Structure**: Moved from `pipelines:` at root to `service.pipelines:`

#### Configuration Updates

- Updated `configs/otel-collector.yaml` with connectors and new telemetry format
- Updated `tests/e2e/testdata/otel-collector.yaml` with memory limiter and new format
- Updated batch processor settings: `timeout: 200ms`, `send_batch_size: 8192`
- Updated telemetry format to use `service.telemetry.metrics.readers`

### Fixed

- **Trace Exporter**: Fixed issue where spans were not being properly exported
  - StartSpanDirect now correctly returns OpenTelemetry span IDs
  - Resolved "span-id-placeholder" issue in v1.0.x

### Dependencies

- OpenTelemetry Go SDK: v1.39.0 (latest stable)
- OpenTelemetry Exporters: v1.39.0
- gRPC: v1.77.0
- Compatible with OTEL Collector Contrib v0.142.0

### Removed

- Removed `ocb-collector.yaml.tpl` (replaced by `otel-collector.yaml.tpl`)
- Removed `otel-collector-config.yaml.tpl` (duplicate file)
- Removed custom `enabled` flags from all configuration templates

### Migration Guide

#### From 1.0.x to 1.1.0

**Configuration Format Changes:**

```yaml
# Old format (1.0.x) - Custom with enabled flags
receivers:
  otlp:
    enabled: true  # No longer supported
    protocols:
      grpc:
        enabled: true  # No longer supported
        endpoint: "0.0.0.0:4317"

exporters:
  logging:  # Renamed
    enabled: true
    loglevel: "info"

pipelines:  # Moved under service
  metrics:
    receivers: [otlp]
    exporters: [logging]

# New format (1.1.0) - Standard OTEL
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"

exporters:
  debug:  # Standard OTEL naming
    verbosity: detailed

service:
  pipelines:  # Under service section
    metrics:
      receivers: [otlp]
      exporters: [debug]
```

**New Connectors (Optional):**

```yaml
connectors:
  spanmetrics:
    exemplars:
      enabled: true
  servicegraph:
    store:
      ttl: 2s

service:
  pipelines:
    traces:
      exporters: [debug, spanmetrics, servicegraph]
    metrics/spanmetrics:
      receivers: [spanmetrics]
      exporters: [prometheus]
```

---

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
- Fixed test assertion for default endpoint in `builder_test.go` (`.io` → `.id`)

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

[1.1.2]: https://github.com/telemetryflow/telemetryflow-go-sdk/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/telemetryflow/telemetryflow-go-sdk/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/telemetryflow/telemetryflow-go-sdk/compare/v1.0.1...v1.1.0
[1.0.1]: https://github.com/telemetryflow/telemetryflow-go-sdk/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/telemetryflow/telemetryflow-go-sdk/releases/tag/v1.0.0
