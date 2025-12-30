# {{.ProjectName}}

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-dark.svg">
    <source media="(prefers-color-scheme: light)" srcset="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg">
    <img src="https://github.com/telemetryflow/.github/raw/main/docs/assets/tfo-logo-sdk-light.svg" alt="TelemetryFlow Logo" width="80%">
  </picture>

[![Version](https://img.shields.io/badge/Version-1.1.1-orange.svg)](CHANGELOG.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/)
[![OpenTelemetry](https://img.shields.io/badge/OTLP-100%25%20Compliant-success?logo=opentelemetry)](https://opentelemetry.io/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://hub.docker.com/r/telemetryflow/telemetryflow-sdk)

</div>

<p align="center">
<strong>[GENERATED TelemetryFlow SDK]</strong> Order-Service - {{.ServiceName}} - RESTful API with DDD + CQRS Pattern
</p>

## Architecture

This project follows **Domain-Driven Design (DDD)** with **CQRS (Command Query Responsibility Segregation)** pattern.

```
{{.ProjectName}}/
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain Layer (Core Business Logic)
│   │   ├── entity/             # Domain entities
│   │   ├── repository/         # Repository interfaces
│   │   └── valueobject/        # Value objects
│   ├── application/            # Application Layer (Use Cases)
│   │   ├── command/            # Commands (write operations)
│   │   ├── query/              # Queries (read operations)
│   │   ├── handler/            # Command & Query handlers
│   │   └── dto/                # Data Transfer Objects
│   └── infrastructure/         # Infrastructure Layer
│       ├── persistence/        # Database implementations
│       ├── http/               # HTTP server & handlers
│       └── config/             # Configuration
├── pkg/                        # Shared packages
├── telemetry/                  # TelemetryFlow integration
├── config/                     # Service configurations
│   └── otel/                   # OpenTelemetry Collector config
├── docs/                       # Documentation
│   ├── api/                    # OpenAPI/Swagger specs
│   ├── diagrams/               # ERD, DFD diagrams
│   └── postman/                # Postman collections
├── migrations/                 # Database migrations
├── configs/                    # Application configuration files
└── tests/                      # Tests
    ├── unit/
    ├── integration/
    └── e2e/
```

## Quick Start

### Prerequisites

- Go 1.24+
{{- if eq .DBDriver "postgres"}}
- PostgreSQL 16+
{{- else if eq .DBDriver "mysql"}}
- MySQL 8+
{{- else}}
- SQLite 3+
{{- end}}
- Docker & Docker Compose (recommended)

### Setup

1. Clone the repository

2. Copy environment file:
   ```bash
   cp .env.example .env
   ```

3. Edit `.env` with your configuration

4. Install dependencies:
   ```bash
   make deps
   ```

5. Run migrations:
   ```bash
   make migrate-up
   ```

6. Start the server:
   ```bash
   make run
   ```

## Docker Compose

The easiest way to run the service with all dependencies:

```bash
# Start all services ({{if eq .DBDriver "postgres"}}PostgreSQL{{else if eq .DBDriver "mysql"}}MySQL{{end}} + API{{if .EnableTelemetry}} + OpenTelemetry Collector{{end}})
make docker-compose-up

# Or use profiles for selective startup
docker compose --profile all up -d        # Start everything
docker compose --profile db up -d         # Start only database
docker compose --profile app up -d        # Start only API
{{- if .EnableTelemetry}}
docker compose --profile monitoring up -d # Start only OTEL Collector
{{- end}}

# Stop all services
make docker-compose-down

# View logs
docker logs -f {{.ProjectName | lower}}_api
docker logs -f {{.ProjectName | lower}}_{{if eq .DBDriver "postgres"}}postgres{{else if eq .DBDriver "mysql"}}mysql{{end}}
{{- if .EnableTelemetry}}
docker logs -f {{.ProjectName | lower}}_otel
{{- end}}
```

### Services

| Service | Container | Port | Description |
|---------|-----------|------|-------------|
{{- if eq .DBDriver "postgres"}}
| PostgreSQL | `{{.ProjectName | lower}}_postgres` | {{.DBPort}} | Database |
{{- else if eq .DBDriver "mysql"}}
| MySQL | `{{.ProjectName | lower}}_mysql` | {{.DBPort}} | Database |
{{- end}}
| API | `{{.ProjectName | lower}}_api` | {{.ServerPort}} | RESTful API |
{{- if .EnableTelemetry}}
| OTEL Collector | `{{.ProjectName | lower}}_otel` | 4317, 4318, 8889, 13133 | OpenTelemetry Collector |
| Jaeger | `{{.ProjectName | lower}}_jaeger` | 16686 | Distributed Tracing UI |
{{- end}}

## Development

### Running locally

```bash
# Build and run
make run

# Run with hot reload
make dev

# Run tests
make test

# Build binary
make build
```

### Adding a new entity

Use the TelemetryFlow RESTful API Generator:

```bash
telemetryflow-restapi entity -n Product -f 'name:string,price:float64,stock:int'
```

This generates:
- Domain entity
- Repository interface & implementation
- CQRS commands & queries
- HTTP handlers
- Database migration

## API Documentation

| Documentation | Location |
|--------------|----------|
| OpenAPI Spec | `docs/api/openapi.yaml` |
| Swagger JSON | `docs/api/swagger.json` |
| ERD Diagram | `docs/diagrams/ERD.md` |
| DFD Diagram | `docs/diagrams/DFD.md` |
| Postman Collection | `docs/postman/collection.json` |

### API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/api/v1/docs` | Swagger UI |

## Configuration

Configuration is loaded from environment variables and `.env` file.

### Application Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `{{.ServerPort}}` |
| `SERVER_READ_TIMEOUT` | Read timeout | `15s` |
| `SERVER_WRITE_TIMEOUT` | Write timeout | `15s` |
| `ENV` | Environment (development/production) | `{{.Environment}}` |

### Database Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_DRIVER` | Database driver | `{{.DBDriver}}` |
| `DB_HOST` | Database host | `{{.DBHost}}` |
| `DB_PORT` | Database port | `{{.DBPort}}` |
| `DB_NAME` | Database name | `{{.DBName}}` |
| `DB_USER` | Database user | `{{.DBUser}}` |
| `DB_PASSWORD` | Database password | - |
| `DB_SSL_MODE` | SSL mode | `disable` |
| `DB_MAX_OPEN_CONNS` | Max open connections | `25` |
| `DB_MAX_IDLE_CONNS` | Max idle connections | `5` |
| `DB_CONN_MAX_LIFETIME` | Connection max lifetime | `5m` |

{{- if .EnableAuth}}

### JWT Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing secret | - |
| `JWT_REFRESH_SECRET` | JWT refresh secret | - |
| `JWT_EXPIRATION` | Token expiration | `24h` |
| `JWT_REFRESH_EXPIRATION` | Refresh token expiration | `168h` |
{{- end}}

{{- if .EnableTelemetry}}

### TelemetryFlow / OpenTelemetry Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `TELEMETRYFLOW_API_KEY_ID` | TelemetryFlow API Key ID | - |
| `TELEMETRYFLOW_API_KEY_SECRET` | TelemetryFlow API Key Secret | - |
| `TELEMETRYFLOW_ENDPOINT` | OTLP endpoint | `localhost:4317` |
| `TELEMETRYFLOW_SERVICE_NAME` | Service name | `{{.ServiceName}}` |
| `TELEMETRYFLOW_SERVICE_VERSION` | Service version | `{{.ServiceVersion}}` |
{{- end}}

### Docker Compose Configuration

| Variable | Description | Default |
|----------|-------------|---------|
{{- if eq .DBDriver "postgres"}}
| `POSTGRES_VERSION` | PostgreSQL image version | `16-alpine` |
{{- else if eq .DBDriver "mysql"}}
| `MYSQL_VERSION` | MySQL image version | `8.0` |
{{- end}}
{{- if .EnableTelemetry}}
| `OTEL_VERSION` | OTEL Collector image version | `0.142.0` |
| `JAEGER_VERSION` | Jaeger image version | `2.2.0` |
{{- end}}
| `CONTAINER_API` | API container name | `{{.ProjectName | lower}}_api` |

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Run E2E tests
make test-e2e

# Generate coverage report
make test-coverage
```

## Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Start all services (app + database{{if .EnableTelemetry}} + monitoring{{end}})
make docker-compose-up

# Stop all services
make docker-compose-down
```

{{- if .EnableTelemetry}}

## Observability

The service is instrumented with OpenTelemetry for:
- **Traces**: Distributed tracing for request flows
- **Metrics**: Application and runtime metrics
- **Logs**: Structured logging

The OpenTelemetry Collector receives telemetry data and can export to:
- Jaeger (tracing)
- Prometheus (metrics)
- Any OTLP-compatible backend

### Prometheus Metrics

Access metrics at: `http://localhost:8889/metrics`

### Jaeger UI

Access Jaeger UI at: `http://localhost:16686`
{{- end}}

## License

Copyright (c) 2024-2026 {{.ProjectName}}. All rights reserved.
