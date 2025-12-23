# TelemetryFlow Integration

This directory contains the TelemetryFlow observability integration for {{.ProjectName}}.

## Setup

1. Install dependencies:
   ```bash
   go get github.com/telemetryflow/telemetryflow-go-sdk
   ```

2. Configure your environment:
   Copy .env.telemetryflow to .env and fill in your API credentials.

3. Initialize in your application:
   ```go
   package main

   import (
       "log"
       "{{.ModulePath}}/telemetry"
   )

   func main() {
       if err := telemetry.Init(); err != nil {
           log.Fatal(err)
       }
       defer telemetry.Shutdown()

       // Your application code...
   }
   ```

## Usage Examples

### Metrics

```go
import "{{.ModulePath}}/telemetry/metrics"

metrics.IncrementCounter("orders.created", 1, map[string]interface{}{
    "customer_type": "premium",
})

metrics.RecordGauge("active_connections", 42, nil)

metrics.RecordHistogram("request.duration", 0.125, "s", map[string]interface{}{
    "endpoint": "/api/users",
})
```

### Logs

```go
import "{{.ModulePath}}/telemetry/logs"

logs.Info("User logged in", map[string]interface{}{
    "user_id": "123",
    "ip": "192.168.1.1",
})

logs.Error("Payment failed", map[string]interface{}{
    "order_id": "ord_456",
    "reason": "insufficient_funds",
})
```

### Traces

```go
import "{{.ModulePath}}/telemetry/traces"

spanID, _ := traces.StartSpan(ctx, "process-payment", map[string]interface{}{
    "amount": 100.50,
})
defer traces.EndSpan(ctx, spanID, nil)

// Add events during processing
traces.AddEvent(ctx, spanID, "payment.validated", nil)
```

## Generated Files

- `init.go` - SDK initialization and shutdown
- `metrics/` - Metrics helpers
- `logs/` - Logging helpers
- `traces/` - Tracing helpers
- `.env.telemetryflow` - Configuration template

## Signals Enabled

{{if .EnableMetrics}}- Metrics{{end}}
{{if .EnableLogs}}- Logs{{end}}
{{if .EnableTraces}}- Traces{{end}}

## Documentation

For more information, visit: https://docs.telemetryflow.id
