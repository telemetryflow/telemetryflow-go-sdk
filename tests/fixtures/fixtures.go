// =============================================================================
// Test Fixtures - Package Documentation & Helpers
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// Package fixtures provides test data and utilities for testing the
// TelemetryFlow Go SDK. It includes sample data for:
//
//   - Credentials: Valid and invalid API key scenarios
//   - Configuration: Various SDK configuration options
//   - Telemetry: Sample metrics, logs, and traces
//   - Templates: Code generator template data
//   - Responses: Mock HTTP/gRPC responses
//
// =============================================================================

package fixtures

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// FixturesDir returns the path to the fixtures directory
func FixturesDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}

// LoadJSONFixture loads a JSON fixture file and unmarshals it into the target
func LoadJSONFixture(filename string, target interface{}) error {
	path := filepath.Join(FixturesDir(), filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// CreateTempEnvFile creates a temporary .env file with the given content
func CreateTempEnvFile(content string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "telemetryflow-test-*")
	if err != nil {
		return "", nil, err
	}

	envPath := filepath.Join(tmpDir, ".env.telemetryflow")
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return envPath, cleanup, nil
}

// CreateTempConfigFile creates a temporary config file with the given content
func CreateTempConfigFile(filename, content string) (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "telemetryflow-test-*")
	if err != nil {
		return "", nil, err
	}

	configPath := filepath.Join(tmpDir, filename)
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return configPath, cleanup, nil
}

// SetEnvVars sets environment variables and returns a cleanup function
func SetEnvVars(vars map[string]string) func() {
	original := make(map[string]string)

	for key, value := range vars {
		original[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	return func() {
		for key, value := range original {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}
}

// DefaultEnvVars returns a map of default environment variables for testing
func DefaultEnvVars() map[string]string {
	return map[string]string{
		"TELEMETRYFLOW_API_KEY_ID":     ValidCredentials.KeyID,
		"TELEMETRYFLOW_API_KEY_SECRET": ValidCredentials.KeySecret,
		"TELEMETRYFLOW_ENDPOINT":       "localhost:4317",
		"TELEMETRYFLOW_SERVICE_NAME":   "test-service",
		"TELEMETRYFLOW_SERVICE_VERSION": "1.0.0",
		"TELEMETRYFLOW_ENVIRONMENT":    "test",
		"TELEMETRYFLOW_PROTOCOL":       "grpc",
		"TELEMETRYFLOW_INSECURE":       "true",
		"TELEMETRYFLOW_ENABLE_METRICS": "true",
		"TELEMETRYFLOW_ENABLE_LOGS":    "true",
		"TELEMETRYFLOW_ENABLE_TRACES":  "true",
	}
}

// GetRandomMetric returns a random metric fixture
func GetRandomMetric() MetricFixture {
	if len(SampleMetrics) == 0 {
		return MetricFixture{}
	}
	return SampleMetrics[0]
}

// GetRandomLog returns a random log fixture
func GetRandomLog() LogFixture {
	if len(SampleLogs) == 0 {
		return LogFixture{}
	}
	return SampleLogs[0]
}

// GetRandomSpan returns a random span fixture
func GetRandomSpan() SpanFixture {
	if len(SampleSpans) == 0 {
		return SpanFixture{}
	}
	return SampleSpans[0]
}

// GetMetricsByType returns metrics filtered by type
func GetMetricsByType(metricType string) []MetricFixture {
	var result []MetricFixture
	for _, m := range SampleMetrics {
		if m.Type == metricType {
			result = append(result, m)
		}
	}
	return result
}

// GetLogsByLevel returns logs filtered by level
func GetLogsByLevel(level string) []LogFixture {
	var result []LogFixture
	for _, l := range SampleLogs {
		if l.Level == level {
			result = append(result, l)
		}
	}
	return result
}

// GetSpansByKind returns spans filtered by kind
func GetSpansByKind(kind string) []SpanFixture {
	var result []SpanFixture
	for _, s := range SampleSpans {
		if s.Kind == kind {
			result = append(result, s)
		}
	}
	return result
}

// GetErrorSpans returns spans that have errors
func GetErrorSpans() []SpanFixture {
	var result []SpanFixture
	for _, s := range SampleSpans {
		if s.HasError {
			result = append(result, s)
		}
	}
	return result
}
