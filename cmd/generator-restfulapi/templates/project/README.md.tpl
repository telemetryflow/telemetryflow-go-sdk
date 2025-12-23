# {{.ProjectName}}

{{.ServiceName}} - RESTful API with DDD + CQRS Pattern

## Architecture

This project follows **Domain-Driven Design (DDD)** with **CQRS (Command Query Responsibility Segregation)** pattern.

```
{{.ProjectName}}/
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── domain/                 # Domain Layer (Core Business Logic)
│   │   ├── entity/            # Domain entities
│   │   ├── repository/        # Repository interfaces
│   │   └── valueobject/       # Value objects
│   ├── application/            # Application Layer (Use Cases)
│   │   ├── command/           # Commands (write operations)
│   │   ├── query/             # Queries (read operations)
│   │   ├── handler/           # Command & Query handlers
│   │   └── dto/               # Data Transfer Objects
│   └── infrastructure/         # Infrastructure Layer
│       ├── persistence/       # Database implementations
│       ├── http/              # HTTP server & handlers
│       └── config/            # Configuration
├── pkg/                        # Shared packages
├── telemetry/                  # TelemetryFlow integration
├── docs/                       # Documentation
│   ├── api/                   # OpenAPI/Swagger specs
│   ├── diagrams/              # ERD, DFD diagrams
│   └── postman/               # Postman collections
├── migrations/                 # Database migrations
├── configs/                    # Configuration files
└── tests/                      # Tests
    ├── unit/
    ├── integration/
    └── e2e/
```

## Quick Start

### Prerequisites

- Go 1.22+
- {{if eq .DBDriver "postgres"}}PostgreSQL 14+{{else if eq .DBDriver "mysql"}}MySQL 8+{{else}}SQLite 3+{{end}}
- Docker (optional)

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

## Development

### Running locally

```bash
# Build and run
make run

# Run with hot reload
make dev

# Run tests
make test
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

Configuration is loaded from environment variables and `configs/config.yaml`.

| Variable | Description | Default |
|----------|-------------|---------|
| `SERVER_PORT` | HTTP server port | `{{.ServerPort}}` |
| `DB_DRIVER` | Database driver | `{{.DBDriver}}` |
| `DB_HOST` | Database host | `{{.DBHost}}` |
| `DB_PORT` | Database port | `{{.DBPort}}` |
| `DB_NAME` | Database name | `{{.DBName}}` |
| `DB_USER` | Database user | `{{.DBUser}}` |
| `DB_PASSWORD` | Database password | - |
{{- if .EnableTelemetry}}
| `TELEMETRYFLOW_API_KEY_ID` | TelemetryFlow API Key ID | - |
| `TELEMETRYFLOW_API_KEY_SECRET` | TelemetryFlow API Key Secret | - |
| `TELEMETRYFLOW_ENDPOINT` | TelemetryFlow endpoint | `api.telemetryflow.id:4317` |
{{- end}}

## Testing

```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make test-coverage
```

## Docker

```bash
# Build image
make docker-build

# Run container
make docker-run

# Start all services (app + database)
make docker-compose-up
```

## License

Copyright (c) 2024-2026 {{.ProjectName}}. All rights reserved.
