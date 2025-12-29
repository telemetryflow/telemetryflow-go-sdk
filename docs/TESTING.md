# Testing Documentation

Testing guide for the TelemetryFlow Go SDK.

- **Version:** 1.1.1
- **Last Updated:** December 2025

---

## Overview

TelemetryFlow Go SDK uses a comprehensive testing strategy with three levels of testing:

1. **Unit Tests** - Test individual components in isolation
2. **Integration Tests** - Test component interactions
3. **End-to-End Tests** - Test complete data flows

All tests follow Domain-Driven Design (DDD) patterns for organization.

---

## Test Structure

```text
tests/
├── unit/                       # Unit tests (DDD organized)
│   ├── domain/                 # Core business logic tests
│   │   ├── config/             # TelemetryConfig entity tests
│   │   └── credentials/        # Credentials value object tests
│   ├── application/            # Use case tests
│   │   └── commands/           # CQRS command tests
│   ├── infrastructure/         # External adapter tests
│   │   ├── config/             # Configuration file handling
│   │   ├── template/           # Template processing tests
│   │   ├── http/               # HTTP middleware tests
│   │   └── database/           # Database pattern tests
│   ├── presentation/           # UI/Output tests
│   │   ├── client/             # Client API tests
│   │   ├── builder/            # Builder pattern tests
│   │   ├── version/            # Version display tests
│   │   └── banner/             # Banner display tests
│   └── generator/              # Code generator tests
│       ├── generator/          # telemetryflow-gen tests
│       └── restapi/            # telemetryflow-restapi tests
├── integration/                # Integration tests
├── e2e/                        # End-to-end tests
│   └── testdata/               # E2E test fixtures
├── fixtures/                   # Test data and fixtures
└── mocks/                      # Mock implementations
```

---

## Running Tests

### All Tests

```bash
# Run all tests
make test

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...

# Run with coverage
go test -cover ./...
```

### Unit Tests

```bash
# All unit tests
go test ./tests/unit/...

# By DDD layer
go test ./tests/unit/domain/...
go test ./tests/unit/application/...
go test ./tests/unit/infrastructure/...
go test ./tests/unit/presentation/...

# Specific package
go test ./tests/unit/domain/config/...
go test ./tests/unit/presentation/client/...
```

### Integration Tests

```bash
# All integration tests
go test ./tests/integration/...
```

### End-to-End Tests

```bash
# All e2e tests (requires environment setup)
TELEMETRYFLOW_E2E=true go test ./tests/e2e/...

# Run with timeout (e2e tests may take longer)
go test -timeout 5m ./tests/e2e/...
```

---

## Coverage

### Generate Coverage Report

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

# Open HTML report (macOS)
open coverage.html
```

### Coverage Targets

| Layer          | Package       | Target | Description                    |
|----------------|---------------|--------|--------------------------------|
| Domain         | config        | 90%    | TelemetryConfig entity         |
| Domain         | credentials   | 95%    | Credentials value object       |
| Application    | commands      | 85%    | CQRS commands                  |
| Infrastructure | exporters     | 90%    | OTLP exporters                 |
| Infrastructure | handlers      | 90%    | Command handlers               |
| Presentation   | client        | 90%    | Client public API              |
| Presentation   | builder       | 95%    | Builder pattern                |
| Presentation   | version       | 100%   | Version information            |
| Presentation   | banner        | 90%    | Banner display                 |

---

## Writing Tests

### Test File Naming

- Unit tests: `*_test.go` in `tests/unit/<layer>/<package>/`
- Integration tests: `*_test.go` in `tests/integration/<package>/`
- E2E tests: `*_test.go` in `tests/e2e/`

### Test Package Pattern

Use external test packages for black-box testing:

```go
// Good: External test package
package config_test

import (
    "github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

// Bad: Same package (white-box testing)
package domain
```

### Test Function Naming

```go
// Test function: TestFunctionName
func TestNewTelemetryConfig(t *testing.T) {
    // Subtests: t.Run("should do something", ...)
    t.Run("should create config with valid inputs", func(t *testing.T) {
        // Test implementation
    })

    t.Run("should reject nil credentials", func(t *testing.T) {
        // Test implementation
    })
}
```

### Example Unit Test

```go
package config_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    "github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func TestNewTelemetryConfig(t *testing.T) {
    t.Run("should create config with valid inputs", func(t *testing.T) {
        // Given
        creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
        require.NoError(t, err)

        // When
        config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

        // Then
        require.NoError(t, err)
        require.NotNil(t, config)
        assert.Equal(t, "localhost:4317", config.Endpoint())
        assert.Equal(t, "my-service", config.ServiceName())
    })

    t.Run("should set correct defaults", func(t *testing.T) {
        // Given
        creds, _ := domain.NewCredentials("tfk_test", "tfs_secret")

        // When
        config, _ := domain.NewTelemetryConfig(creds, "localhost:4317", "my-service")

        // Then
        assert.Equal(t, domain.ProtocolGRPC, config.Protocol())
        assert.Equal(t, 30*time.Second, config.Timeout())
        assert.True(t, config.IsRetryEnabled())
    })
}
```

### Table-Driven Tests

```go
func TestCredentials_Validation(t *testing.T) {
    tests := []struct {
        name      string
        keyID     string
        keySecret string
        wantErr   bool
        errMsg    string
    }{
        {
            name:      "valid credentials",
            keyID:     "tfk_valid",
            keySecret: "tfs_valid",
            wantErr:   false,
        },
        {
            name:      "empty key ID",
            keyID:     "",
            keySecret: "tfs_valid",
            wantErr:   true,
            errMsg:    "key ID",
        },
        {
            name:      "invalid key ID prefix",
            keyID:     "invalid_key",
            keySecret: "tfs_valid",
            wantErr:   true,
            errMsg:    "tfk_",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := domain.NewCredentials(tt.keyID, tt.keySecret)

            if tt.wantErr {
                require.Error(t, err)
                assert.Contains(t, err.Error(), tt.errMsg)
            } else {
                require.NoError(t, err)
            }
        })
    }
}
```

---

## Using Mocks

### Import Mocks

```go
import (
    "github.com/telemetryflow/telemetryflow-go-sdk/tests/mocks"
)
```

### Example with Mock

```go
func TestClientWithMockExporter(t *testing.T) {
    t.Run("should export data successfully", func(t *testing.T) {
        // Given
        mockExporter := mocks.NewMockExporter()
        mockExporter.On("Export", mock.Anything, mock.Anything).Return(nil)

        client := NewClientWithExporter(mockExporter)

        // When
        err := client.IncrementCounter(ctx, "test", 1, nil)

        // Then
        require.NoError(t, err)
        mockExporter.AssertExpectations(t)
    })
}
```

---

## Using Fixtures

### Load Test Fixtures

```go
import (
    "github.com/telemetryflow/telemetryflow-go-sdk/tests/fixtures"
)

func TestWithFixture(t *testing.T) {
    // Get valid credentials fixture
    creds := fixtures.ValidCredentials()

    // Get sample config fixture
    config := fixtures.SampleConfig()

    // Use fixtures in test
    // ...
}
```

---

## Benchmark Tests

### Writing Benchmarks

```go
func BenchmarkNewCredentials(b *testing.B) {
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = domain.NewCredentials("tfk_benchmark_key", "tfs_benchmark_secret")
    }
}

func BenchmarkRecordMetricCommand_Create(b *testing.B) {
    attrs := map[string]interface{}{
        "method": "GET",
        "status": 200,
    }
    now := time.Now()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = &application.RecordMetricCommand{
            Name:       "http.request.duration",
            Value:      0.125,
            Attributes: attrs,
            Timestamp:  now,
        }
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkNewCredentials ./tests/unit/domain/...

# Run with memory allocation stats
go test -bench=. -benchmem ./...

# Run multiple times for accurate results
go test -bench=. -count=5 ./...
```

---

## CI/CD Integration

### GitHub Actions

Tests are automatically run on every push and pull request:

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - name: Run tests
        run: make test
      - name: Run coverage
        run: go test -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
```

---

## Best Practices

1. **Test in isolation**: Use mocks for external dependencies
2. **Use table-driven tests**: For multiple test cases
3. **Test error paths**: Cover both success and failure scenarios
4. **Use testify**: Prefer `assert` and `require` from testify
5. **Follow DDD boundaries**: Keep tests within their architectural layer
6. **Use external packages**: `package <name>_test` for black-box testing
7. **Name tests descriptively**: `t.Run("should do something", ...)`
8. **Keep tests fast**: Unit tests should complete in milliseconds
9. **Clean up resources**: Use `t.Cleanup()` or `defer` for teardown
10. **Avoid test interdependence**: Each test should be independent
11. **Use `t.Setenv`**: For environment variables (auto-cleanup)

---

## Troubleshooting

### Common Issues

**Tests timeout:**
```bash
# Increase timeout
go test -timeout 10m ./...
```

**Race conditions:**
```bash
# Detect races
go test -race ./...
```

**Flaky tests:**
```bash
# Run multiple times to detect flakiness
go test -count=10 ./...
```

**Coverage not accurate:**
```bash
# Include all packages
go test -coverpkg=./... -coverprofile=coverage.out ./...
```

---

## References

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify](https://github.com/stretchr/testify)
- [Test Mocks](../tests/mocks/)
- [Test Fixtures](../tests/fixtures/)
- [DDD Architecture Guide](ARCHITECTURE.md)

---

**Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.**
