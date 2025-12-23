# TelemetryFlow Go SDK - Architecture

This document explains the architectural decisions and patterns used in the TelemetryFlow Go SDK.

## Table of Contents

- [Overview](#overview)
- [High-Level Architecture](#high-level-architecture)
- [Domain-Driven Design (DDD)](#domain-driven-design-ddd)
- [CQRS Pattern](#cqrs-pattern)
- [Layer Architecture](#layer-architecture)
- [Design Principles](#design-principles)
- [Data Flow](#data-flow)
- [Error Handling](#error-handling)
- [Testing Strategy](#testing-strategy)
- [Extension Points](#extension-points)
- [Performance Considerations](#performance-considerations)

## Overview

The TelemetryFlow Go SDK is built using **Domain-Driven Design (DDD)** and **Command Query Responsibility Segregation (CQRS)** patterns. This architecture provides:

- Clear separation of concerns
- Maintainable and testable code
- Easy to extend for new features
- Type-safe with compile-time guarantees
- Production-ready with comprehensive error handling

## High-Level Architecture

```mermaid
graph TB
    subgraph "User Application"
        APP[Application Code]
    end

    subgraph "TelemetryFlow Go SDK"
        subgraph "Interface Layer"
            CLIENT[Client]
            BUILDER[Builder]
        end

        subgraph "Application Layer"
            CMD[Commands]
            QRY[Queries]
            CMDBUS[Command Bus]
            QRYBUS[Query Bus]
        end

        subgraph "Domain Layer"
            CONFIG[TelemetryConfig]
            CREDS[Credentials]
            PROTO[Protocol]
            SIGNAL[SignalType]
        end

        subgraph "Infrastructure Layer"
            HANDLER[Command Handler]
            EXPORTER[OTLP Exporters]
            OTEL[OpenTelemetry SDK]
        end
    end

    subgraph "TelemetryFlow Backend"
        BACKEND[TelemetryFlow API]
    end

    APP --> CLIENT
    APP --> BUILDER
    BUILDER --> CONFIG
    CLIENT --> CMDBUS
    CLIENT --> QRYBUS
    CMDBUS --> HANDLER
    QRYBUS --> HANDLER
    HANDLER --> OTEL
    OTEL --> EXPORTER
    EXPORTER -->|gRPC/HTTP| BACKEND
    CONFIG --> CREDS
    CONFIG --> PROTO
    CONFIG --> SIGNAL
```

## Domain-Driven Design (DDD)

### Bounded Contexts

The SDK has a single bounded context: **Telemetry**

```mermaid
graph LR
    subgraph "Telemetry Bounded Context"
        direction TB
        A[Configuration Management]
        B[Credential Handling]
        C[Signal Type Definitions]
        D[Protocol Abstraction]
    end

    A --- B
    A --- C
    A --- D
```

This context contains all domain logic related to:
- Configuration management
- Credential handling
- Signal type definitions (metrics, logs, traces)

### Entities and Value Objects

```mermaid
classDiagram
    class TelemetryConfig {
        <<Entity - Aggregate Root>>
        -credentials *Credentials
        -endpoint string
        -serviceName string
        -serviceVersion string
        -environment string
        -protocol Protocol
        -enabledSignals []SignalType
        +WithEndpoint(string) *TelemetryConfig
        +WithServiceName(string) *TelemetryConfig
        +WithProtocol(Protocol) *TelemetryConfig
        +EnableSignals(...SignalType) *TelemetryConfig
    }

    class Credentials {
        <<Value Object>>
        -keyID string
        -keySecret string
        +NewCredentials(keyID, keySecret) (*Credentials, error)
        +KeyID() string
        +KeySecret() string
    }

    class Protocol {
        <<Value Object>>
        +ProtocolGRPC
        +ProtocolHTTP
    }

    class SignalType {
        <<Value Object>>
        +SignalMetrics
        +SignalLogs
        +SignalTraces
    }

    TelemetryConfig *-- Credentials : contains
    TelemetryConfig --> Protocol : uses
    TelemetryConfig --> SignalType : enables
```

#### Value Objects (Immutable)

**Credentials** (`domain/credentials.go`)
```go
type Credentials struct {
    keyID     string
    keySecret string
}
```

Value Object characteristics:
- Immutable after creation
- Validates itself on construction
- Equality based on content
- No identity separate from its attributes

**Why Value Object?**
- API credentials are immutable by nature
- Should be validated immediately
- No lifecycle management needed

#### Entities (Mutable)

**TelemetryConfig** (`domain/config.go`)
```go
type TelemetryConfig struct {
    credentials      *Credentials
    endpoint         string
    // ... other configuration
}
```

Entity characteristics:
- Has identity (tied to credentials)
- Can be modified through builder methods
- Represents the aggregate root

**Why Entity?**
- Configuration evolves during setup
- Builder pattern requires mutability
- Acts as aggregate root for telemetry context

### Domain Services

Domain services encapsulate domain logic that doesn't naturally fit into entities or value objects.

The SDK keeps domain services minimal, with most logic in entities and value objects.

## CQRS Pattern

```mermaid
graph TB
    subgraph "User Application"
        UC[User Code]
    end

    subgraph "CQRS Architecture"
        direction TB

        subgraph "Write Side"
            CMD1[RecordMetricCommand]
            CMD2[EmitLogCommand]
            CMD3[StartSpanCommand]
            CMDBUS[Command Bus]
            CMDH[Command Handler]
        end

        subgraph "Read Side"
            QRY1[GetMetricQuery]
            QRY2[GetLogsQuery]
            QRY3[GetTraceQuery]
            QRYBUS[Query Bus]
            QRYH[Query Handler]
        end
    end

    subgraph "Backend"
        OTEL[OpenTelemetry SDK]
        API[TelemetryFlow API]
    end

    UC -->|Write| CMD1 & CMD2 & CMD3
    UC -->|Read| QRY1 & QRY2 & QRY3

    CMD1 & CMD2 & CMD3 --> CMDBUS
    CMDBUS --> CMDH
    CMDH --> OTEL
    OTEL -->|Export| API

    QRY1 & QRY2 & QRY3 --> QRYBUS
    QRYBUS --> QRYH
    QRYH -->|Fetch| API
```

### Why CQRS?

CQRS separates **Commands** (write operations) from **Queries** (read operations).

**Benefits:**
1. **Clear Intent**: Commands express what should happen
2. **Testability**: Commands can be tested in isolation
3. **Extensibility**: New commands don't affect queries
4. **Scalability**: Different optimization strategies for reads/writes

### Command Side

Commands represent intentions to change state:

```mermaid
classDiagram
    class Command {
        <<interface>>
        +Execute(ctx) error
    }

    class RecordMetricCommand {
        +Name string
        +Value float64
        +Unit string
        +Attributes map
        +Timestamp time
    }

    class RecordCounterCommand {
        +Name string
        +Value int64
        +Attributes map
    }

    class RecordGaugeCommand {
        +Name string
        +Value float64
        +Attributes map
    }

    class RecordHistogramCommand {
        +Name string
        +Value float64
        +Buckets array
        +Attributes map
    }

    class EmitLogCommand {
        +Severity string
        +Message string
        +Attributes map
        +Timestamp time
    }

    class StartSpanCommand {
        +Name string
        +Kind string
        +Attributes map
    }

    class EndSpanCommand {
        +SpanID string
        +Status string
        +Error error
    }

    Command <|.. RecordMetricCommand
    Command <|.. RecordCounterCommand
    Command <|.. RecordGaugeCommand
    Command <|.. RecordHistogramCommand
    Command <|.. EmitLogCommand
    Command <|.. StartSpanCommand
    Command <|.. EndSpanCommand
```

```go
type RecordMetricCommand struct {
    Name       string
    Value      float64
    Unit       string
    Attributes map[string]interface{}
    Timestamp  time.Time
}

type EmitLogCommand struct {
    Severity   string
    Message    string
    Attributes map[string]interface{}
    Timestamp  time.Time
}

type StartSpanCommand struct {
    Name       string
    Kind       string
    Attributes map[string]interface{}
}
```

**Command Handler:**
```go
type TelemetryCommandHandler struct {
    config         *domain.TelemetryConfig
    tracerProvider *sdktrace.TracerProvider
    meterProvider  *metric.MeterProvider
}

func (h *TelemetryCommandHandler) Handle(ctx context.Context, cmd Command) error {
    // Dispatch to specific handler
}
```

### Query Side

Queries represent requests for data:

```go
type GetMetricQuery struct {
    Name      string
    StartTime time.Time
    EndTime   time.Time
}

type GetLogsQuery struct {
    StartTime  time.Time
    EndTime    time.Time
    Severity   []string
}

type GetTraceQuery struct {
    TraceID string
}
```

**Query Handler:**
```go
type TelemetryQueryHandler struct {
    // HTTP client or gRPC client for TelemetryFlow API
}

func (h *TelemetryQueryHandler) Handle(ctx context.Context, query Query) (interface{}, error) {
    // Fetch and return data
}
```

### Command/Query Buses

The buses route commands and queries to appropriate handlers:

```go
type CommandBus struct {
    handlers map[string]CommandHandler
}

func (b *CommandBus) Register(commandType string, handler CommandHandler) {
    b.handlers[commandType] = handler
}

func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
    // Find and invoke handler
}
```

## Layer Architecture

The SDK follows a clean architecture with four layers:

```mermaid
graph TB
    subgraph "Layer Dependencies"
        direction TB

        IL[Interface Layer<br/>pkg/telemetryflow/]
        AL[Application Layer<br/>pkg/telemetryflow/application/]
        DL[Domain Layer<br/>pkg/telemetryflow/domain/]
        INF[Infrastructure Layer<br/>pkg/telemetryflow/infrastructure/]

        IL --> AL
        IL --> DL
        IL --> INF
        AL --> DL
        INF --> DL
        INF --> AL
    end

    subgraph "External Dependencies"
        OTEL[OpenTelemetry SDK]
        GRPC[gRPC]
        HTTP[HTTP Client]
    end

    INF --> OTEL
    INF --> GRPC
    INF --> HTTP
    DL -.->|No External Dependencies| DL
```

```mermaid
graph LR
    subgraph "Package Structure"
        direction TB
        ROOT[pkg/telemetryflow/]
        CLIENT[client.go]
        BUILDER[builder.go]

        subgraph "domain/"
            CREDS[credentials.go]
            CONFIG[config.go]
        end

        subgraph "application/"
            CMDS[commands.go]
            QRYS[queries.go]
        end

        subgraph "infrastructure/"
            HANDLERS[handlers.go]
            EXPORTERS[exporters.go]
        end
    end

    ROOT --> CLIENT
    ROOT --> BUILDER
    ROOT --> domain/
    ROOT --> application/
    ROOT --> infrastructure/
```

### 1. Domain Layer (`pkg/telemetryflow/domain/`)

**Responsibility:** Core business logic and rules

**Contains:**
- Entities: `TelemetryConfig`
- Value Objects: `Credentials`, `Protocol`, `SignalType`
- Domain rules and validation

**Dependencies:** None (pure Go)

**Example:**
```go
// Value object with self-validation
func NewCredentials(keyID, keySecret string) (*Credentials, error) {
    if !strings.HasPrefix(keyID, "tfk_") {
        return nil, errors.New("invalid key ID format")
    }
    return &Credentials{keyID, keySecret}, nil
}
```

### 2. Application Layer (`pkg/telemetryflow/application/`)

**Responsibility:** Orchestrate domain objects and implement use cases

**Contains:**
- Commands and Queries (CQRS)
- Command/Query buses
- Application services

**Dependencies:** Domain layer only

**Example:**
```go
// Command represents a use case
type RecordMetricCommand struct {
    Name       string
    Value      float64
    Attributes map[string]interface{}
}
```

### 3. Infrastructure Layer (`pkg/telemetryflow/infrastructure/`)

**Responsibility:** Technical implementation details

**Contains:**
- OTLP exporters (gRPC/HTTP)
- Command handlers
- External integrations (OpenTelemetry SDK)

**Dependencies:** Domain and Application layers

**Example:**
```go
// Infrastructure creates and manages OTLP exporters
type OTLPExporterFactory struct {
    config *domain.TelemetryConfig
}

func (f *OTLPExporterFactory) CreateTraceExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
    // Create gRPC or HTTP exporter based on config
}
```

### 4. Interface Layer (`pkg/telemetryflow/`)

**Responsibility:** Public API that users interact with

**Contains:**
- `Client`: Main SDK client
- `Builder`: Fluent configuration builder
- Public API methods

**Dependencies:** All layers

**Example:**
```go
// Public API
type Client struct {
    config         *domain.TelemetryConfig
    commandHandler *infrastructure.TelemetryCommandHandler
}

func (c *Client) IncrementCounter(ctx context.Context, name string, value int64, attrs map[string]interface{}) error {
    cmd := &application.RecordCounterCommand{Name: name, Value: value}
    return c.commandHandler.Handle(ctx, cmd)
}
```

## Design Principles

### 1. Dependency Inversion

Higher-level modules don't depend on lower-level modules. Both depend on abstractions.

```go
// Bad: Direct dependency on infrastructure
type Client struct {
    grpcExporter *otlptracegrpc.Exporter
}

// Good: Depend on abstraction
type Client struct {
    commandHandler CommandHandler
}
```

### 2. Single Responsibility

Each type has one reason to change.

```go
// Credentials: Only responsible for API key validation
type Credentials struct { /* ... */ }

// TelemetryConfig: Only responsible for configuration
type TelemetryConfig struct { /* ... */ }

// Client: Only responsible for public API
type Client struct { /* ... */ }
```

### 3. Open/Closed Principle

Open for extension, closed for modification.

```go
// Add new commands without modifying existing code
type NewCustomCommand struct {
    // New fields
}

// Register handler
commandBus.Register("new_custom", newCustomHandler)
```

### 4. Interface Segregation

Clients shouldn't depend on interfaces they don't use.

```go
// Separate interfaces for different concerns
type CommandHandler interface {
    Handle(ctx context.Context, cmd Command) error
}

type QueryHandler interface {
    Handle(ctx context.Context, query Query) (interface{}, error)
}
```

### 5. Liskov Substitution

Subtypes must be substitutable for their base types.

```go
// All commands implement Command interface
type Command interface {
    Execute(ctx context.Context) error
}

// Any command can be used where Command is expected
func processCommand(cmd Command) error {
    return cmd.Execute(context.Background())
}
```

## Data Flow

### Write Path (Commands)

```mermaid
sequenceDiagram
    participant App as User Application
    participant Client as SDK Client
    participant Cmd as Command
    participant Bus as Command Bus
    participant Handler as Command Handler
    participant OTel as OpenTelemetry SDK
    participant Exporter as OTLP Exporter
    participant Backend as TelemetryFlow Backend

    App->>Client: IncrementCounter(ctx, "requests", 1)
    Client->>Cmd: Create RecordCounterCommand
    Client->>Bus: Dispatch(cmd)
    Bus->>Handler: Handle(ctx, cmd)
    Handler->>OTel: Record metric via Meter
    OTel->>OTel: Batch & aggregate
    OTel->>Exporter: Export batch
    Exporter->>Backend: OTLP/gRPC or OTLP/HTTP
    Backend-->>Exporter: ACK
    Exporter-->>Handler: Success
    Handler-->>Client: nil (success)
    Client-->>App: nil
```

### Read Path (Queries)

```mermaid
sequenceDiagram
    participant App as User Application
    participant Client as SDK Client
    participant Query as Query
    participant Bus as Query Bus
    participant Handler as Query Handler
    participant API as TelemetryFlow API

    App->>Client: GetMetrics(ctx, query)
    Client->>Query: Create GetMetricQuery
    Client->>Bus: Dispatch(query)
    Bus->>Handler: Handle(ctx, query)
    Handler->>API: HTTP GET /api/v1/metrics
    API-->>Handler: MetricQueryResult
    Handler-->>Bus: Result
    Bus-->>Client: Result
    Client-->>App: []MetricDataPoint
```

### Signal-Specific Data Flows

#### Metrics Flow

```mermaid
flowchart LR
    subgraph "SDK"
        A[RecordMetric] --> B[Counter/Gauge/Histogram]
        B --> C[Meter Provider]
        C --> D[Batch Processor]
    end

    subgraph "Export"
        D --> E{Protocol}
        E -->|gRPC| F[OTLP gRPC Exporter]
        E -->|HTTP| G[OTLP HTTP Exporter]
    end

    subgraph "Backend"
        F --> H[TelemetryFlow Metrics API]
        G --> H
    end
```

#### Traces Flow

```mermaid
flowchart LR
    subgraph "SDK"
        A[StartSpan] --> B[Tracer]
        B --> C[Span]
        C --> D[AddEvent/SetAttribute]
        D --> E[EndSpan]
        E --> F[Span Processor]
    end

    subgraph "Export"
        F --> G{Protocol}
        G -->|gRPC| H[OTLP gRPC Exporter]
        G -->|HTTP| I[OTLP HTTP Exporter]
    end

    subgraph "Backend"
        H --> J[TelemetryFlow Traces API]
        I --> J
    end
```

#### Logs Flow

```mermaid
flowchart LR
    subgraph "SDK"
        A[LogInfo/LogWarn/LogError] --> B[Log Record]
        B --> C[Logger Provider]
        C --> D[Batch Processor]
    end

    subgraph "Export"
        D --> E{Protocol}
        E -->|gRPC| F[OTLP gRPC Exporter]
        E -->|HTTP| G[OTLP HTTP Exporter]
    end

    subgraph "Backend"
        F --> H[TelemetryFlow Logs API]
        G --> H
    end
```

## Error Handling

```mermaid
flowchart TB
    subgraph "Error Propagation"
        direction TB

        subgraph "Domain Layer"
            D1[Validation Error] --> D2[Domain Error]
        end

        subgraph "Application Layer"
            D2 --> A1[Wrap with Context]
            A1 --> A2[Application Error]
        end

        subgraph "Infrastructure Layer"
            A2 --> I1[Add Technical Details]
            EXT[External Service Error] --> I1
            I1 --> I2[Infrastructure Error]
        end

        subgraph "Interface Layer"
            I2 --> IF1[Format for User]
            IF1 --> IF2[Return to Caller]
        end
    end
```

### Domain Layer Errors

Return domain-specific errors:

```go
func NewCredentials(keyID, keySecret string) (*Credentials, error) {
    if keyID == "" {
        return nil, errors.New("API key ID cannot be empty")
    }
    // ...
}
```

### Application Layer Errors

Wrap domain errors with context:

```go
func (h *Handler) Handle(ctx context.Context, cmd Command) error {
    if err := validate(cmd); err != nil {
        return fmt.Errorf("command validation failed: %w", err)
    }
    // ...
}
```

### Infrastructure Layer Errors

Handle technical errors gracefully:

```go
func (f *Factory) CreateExporter(ctx context.Context) (Exporter, error) {
    exporter, err := otlpgrpc.New(ctx, options...)
    if err != nil {
        return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
    }
    return exporter, nil
}
```

## Testing Strategy

```mermaid
graph TB
    subgraph "Testing Pyramid"
        direction TB

        E2E[End-to-End Tests<br/>Full integration with backend]
        INT[Integration Tests<br/>Cross-layer interactions]
        UNIT[Unit Tests<br/>Individual components]
    end

    subgraph "Test Coverage by Layer"
        direction LR

        D[Domain Layer<br/>- Credentials validation<br/>- Config building]
        A[Application Layer<br/>- Command structure<br/>- Query structure]
        I[Infrastructure Layer<br/>- Handler dispatch<br/>- Exporter creation]
        IF[Interface Layer<br/>- Client API<br/>- Builder methods]
    end

    UNIT --> D
    UNIT --> A
    INT --> I
    INT --> IF
    E2E --> IF
```

### Unit Tests

Test each layer independently:

```go
// Domain layer tests
func TestCredentials_Validation(t *testing.T) {
    _, err := domain.NewCredentials("invalid", "secret")
    assert.Error(t, err)
}

// Application layer tests
func TestCommandHandler_RecordMetric(t *testing.T) {
    handler := NewMockHandler()
    cmd := &RecordMetricCommand{Name: "test"}
    err := handler.Handle(context.Background(), cmd)
    assert.NoError(t, err)
}
```

### Integration Tests

Test interactions between layers:

```go
func TestClient_EndToEnd(t *testing.T) {
    client := setupTestClient()
    err := client.IncrementCounter(ctx, "test", 1, nil)
    assert.NoError(t, err)
}
```

## Extension Points

```mermaid
flowchart TB
    subgraph "Extension Points"
        direction TB

        subgraph "New Commands"
            NC1[1. Define Command in application/]
            NC2[2. Add Handler in infrastructure/]
            NC3[3. Wire up in Client API]
            NC1 --> NC2 --> NC3
        end

        subgraph "New Signal Types"
            NS1[1. Add SignalType constant]
            NS2[2. Add Configuration support]
            NS3[3. Implement Command handlers]
            NS4[4. Expose in public API]
            NS1 --> NS2 --> NS3 --> NS4
        end

        subgraph "New Protocols"
            NP1[1. Add Protocol constant]
            NP2[2. Implement Exporter factory]
            NP3[3. Register in handler]
            NP1 --> NP2 --> NP3
        end
    end
```

### Adding New Commands

1. Define command in `application/commands.go`:
```go
type NewFeatureCommand struct {
    Field1 string
    Field2 int
}
```

2. Add handler in infrastructure layer
3. Wire up in `Client` public API

### Adding New Signal Types

1. Add constant in `domain/config.go`:
```go
const SignalEvents SignalType = "events"
```

2. Add configuration support
3. Implement command handlers
4. Expose in public API

## Performance Considerations

```mermaid
graph TB
    subgraph "Performance Optimizations"
        direction TB

        subgraph "Batching"
            B1[Telemetry Data] --> B2[Batch Processor]
            B2 -->|Timeout or Size| B3[Export Batch]
            B2 -->|Continue collecting| B2
        end

        subgraph "Connection Pooling"
            C1[Multiple Requests] --> C2[gRPC Connection Pool]
            C2 --> C3[Reused Connections]
        end

        subgraph "Concurrency"
            CO1[Goroutine 1] --> CO4[Thread-safe Client]
            CO2[Goroutine 2] --> CO4
            CO3[Goroutine N] --> CO4
            CO4 --> CO5[Synchronized Access]
        end
    end
```

### Batching

Commands are batched by OpenTelemetry SDK:

```go
config.WithBatchSettings(
    10 * time.Second,  // batch timeout
    512,               // max batch size
)
```

### Concurrency

The SDK is thread-safe:

```go
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        client.IncrementCounter(ctx, "concurrent", 1, nil)
    }()
}
wg.Wait()
```

### Memory Management

- Credentials are immutable (no copying)
- Commands are short-lived
- Exporters reuse connections

## Conclusion

```mermaid
mindmap
    root((TelemetryFlow SDK))
        Architecture Benefits
            Maintainability
                Clear layer separation
                Single responsibility
            Testability
                Mockable interfaces
                Isolated components
            Extensibility
                Plugin-style commands
                Open/closed principle
        Technical Excellence
            Performance
                Batching
                Connection pooling
                Thread safety
            Type Safety
                Compile-time guarantees
                Strong typing
            Production Ready
                Error handling
                Graceful shutdown
                Retry mechanisms
        Design Patterns
            DDD
                Bounded context
                Entities & Value Objects
                Aggregate roots
            CQRS
                Command separation
                Query separation
                Event-driven ready
```

The TelemetryFlow Go SDK's architecture provides:

| Feature              | Description                                                    |
| -------------------- | -------------------------------------------------------------- |
| **Maintainability**  | Clear layer separation with single responsibility              |
| **Testability**      | Easy to mock and test each layer independently                 |
| **Extensibility**    | Simple to add new commands, signals, and protocols             |
| **Performance**      | Efficient batching, connection pooling, and concurrency        |
| **Type Safety**      | Compile-time guarantees with strong Go typing                  |
| **Production Ready** | Comprehensive error handling, retries, and graceful shutdown   |

This architecture serves as both a fully functional SDK and a reference implementation for applying DDD and CQRS patterns in Go.

---

## Quick Reference

### SDK Initialization Flow

```mermaid
sequenceDiagram
    participant App as Application
    participant Builder as Builder
    participant Config as TelemetryConfig
    participant Client as Client
    participant Handler as CommandHandler

    App->>Builder: New(keyID, keySecret)
    Builder->>Config: Create with defaults
    App->>Builder: WithEndpoint(), WithServiceName()
    Builder->>Config: Apply configurations
    App->>Builder: Build()
    Builder->>Client: Create Client
    Client->>Handler: Initialize handlers
    Handler->>Handler: Setup OpenTelemetry providers
    Client-->>App: Ready to use
```

### Shutdown Flow

```mermaid
sequenceDiagram
    participant App as Application
    participant Client as Client
    participant Handler as CommandHandler
    participant OTel as OpenTelemetry
    participant Exporter as Exporters

    App->>Client: Shutdown(ctx)
    Client->>Handler: Flush pending data
    Handler->>OTel: ForceFlush()
    OTel->>Exporter: Export remaining batches
    Exporter-->>OTel: Complete
    Handler->>OTel: Shutdown providers
    OTel->>Exporter: Close connections
    Exporter-->>Handler: Closed
    Handler-->>Client: Shutdown complete
    Client-->>App: nil (success)
```
