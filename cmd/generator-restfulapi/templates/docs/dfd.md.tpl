# Data Flow Diagram (DFD)

## {{.ProjectName}}

This document describes the data flow for {{.ServiceName}}.

## Level 0 - Context Diagram

```mermaid
graph LR
    subgraph External
        Client[Client Application]
{{- if .EnableTelemetry}}
        TF[TelemetryFlow Platform]
{{- end}}
    end

    subgraph System
        API[{{.ProjectName}} API]
        DB[(Database)]
    end

    Client -->|HTTP Requests| API
    API -->|HTTP Responses| Client
    API -->|CRUD Operations| DB
    DB -->|Query Results| API
{{- if .EnableTelemetry}}
    API -->|Telemetry Data| TF
{{- end}}
```

## Level 1 - System Diagram

```mermaid
graph TB
    subgraph Client Layer
        WEB[Web Client]
        MOBILE[Mobile App]
        CLI[CLI Tool]
    end

    subgraph API Gateway
        LB[Load Balancer]
        RATE[Rate Limiter]
{{- if .EnableAuth}}
        AUTH[Auth Middleware]
{{- end}}
    end

    subgraph Application Layer
        subgraph Presentation
            HANDLER[HTTP Handlers]
        end

        subgraph Application
            CMD[Command Handlers]
            QRY[Query Handlers]
        end

        subgraph Domain
            ENT[Entities]
            REPO[Repository Interfaces]
        end

        subgraph Infrastructure
            PERSIST[Persistence Layer]
            CACHE[Cache Layer]
        end
    end

    subgraph Data Layer
        DB[({{if eq .DBDriver "postgres"}}PostgreSQL{{else if eq .DBDriver "mysql"}}MySQL{{else}}SQLite{{end}})]
        REDIS[(Redis Cache)]
    end

{{- if .EnableTelemetry}}
    subgraph Observability
        TF[TelemetryFlow]
    end
{{- end}}

    WEB --> LB
    MOBILE --> LB
    CLI --> LB
    LB --> RATE
{{- if .EnableAuth}}
    RATE --> AUTH
    AUTH --> HANDLER
{{- else}}
    RATE --> HANDLER
{{- end}}

    HANDLER --> CMD
    HANDLER --> QRY
    CMD --> REPO
    QRY --> REPO
    REPO --> PERSIST
    PERSIST --> DB
    PERSIST -.-> CACHE
    CACHE -.-> REDIS

{{- if .EnableTelemetry}}
    HANDLER --> TF
    CMD --> TF
    PERSIST --> TF
{{- end}}
```

## Level 2 - CQRS Flow

### Command Flow (Write Operations)

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Handler
    participant V as Validator
    participant CH as Command Handler
    participant R as Repository
    participant DB as Database
{{- if .EnableTelemetry}}
    participant T as Telemetry
{{- end}}

    C->>H: POST /api/v1/resource
    H->>V: Validate Request
    V-->>H: Validation Result

    alt Validation Failed
        H-->>C: 400 Bad Request
    else Validation Passed
        H->>CH: Execute Command
{{- if .EnableTelemetry}}
        CH->>T: Start Span
{{- end}}
        CH->>R: Create/Update Entity
        R->>DB: INSERT/UPDATE
        DB-->>R: Result
        R-->>CH: Entity
{{- if .EnableTelemetry}}
        CH->>T: End Span
        CH->>T: Log & Metrics
{{- end}}
        CH-->>H: Command Result
        H-->>C: 201 Created / 200 OK
    end
```

### Query Flow (Read Operations)

```mermaid
sequenceDiagram
    participant C as Client
    participant H as HTTP Handler
    participant QH as Query Handler
    participant R as Repository
    participant CACHE as Cache
    participant DB as Database
{{- if .EnableTelemetry}}
    participant T as Telemetry
{{- end}}

    C->>H: GET /api/v1/resource
    H->>QH: Execute Query
{{- if .EnableTelemetry}}
    QH->>T: Start Span
{{- end}}
    QH->>R: Find Entity
    R->>CACHE: Check Cache

    alt Cache Hit
        CACHE-->>R: Cached Data
    else Cache Miss
        R->>DB: SELECT
        DB-->>R: Query Result
        R->>CACHE: Store in Cache
    end

    R-->>QH: Entity/List
{{- if .EnableTelemetry}}
    QH->>T: End Span
{{- end}}
    QH-->>H: Query Result
    H-->>C: 200 OK
```

## Data Transformations

| Layer | Input | Output | Transformation |
|-------|-------|--------|----------------|
| HTTP Handler | HTTP Request | DTO | Parse JSON, Validate |
| Command Handler | DTO | Domain Event | Apply Business Rules |
| Repository | Entity | DB Row | Serialize to DB Format |
| Query Handler | Query Params | DTO | Projection, Filtering |

## Error Handling Flow

```mermaid
graph TD
    REQ[Request] --> VALID{Validation}
    VALID -->|Invalid| E400[400 Bad Request]
    VALID -->|Valid| AUTH{Authorization}
    AUTH -->|Unauthorized| E401[401 Unauthorized]
    AUTH -->|Forbidden| E403[403 Forbidden]
    AUTH -->|OK| BIZ{Business Logic}
    BIZ -->|Not Found| E404[404 Not Found]
    BIZ -->|Conflict| E409[409 Conflict]
    BIZ -->|Error| E500[500 Internal Error]
    BIZ -->|Success| OK[200/201 Success]
```

## Notes

1. All requests go through rate limiting
2. Commands modify state, Queries are read-only (CQRS)
3. Cache is checked before database on read operations
4. Telemetry captures metrics, logs, and traces at each layer
