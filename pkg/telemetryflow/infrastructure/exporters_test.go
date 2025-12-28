// Package infrastructure provides unit tests for the infrastructure layer.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/domain"
)

func createValidConfig(t *testing.T) *domain.TelemetryConfig {
	creds, err := domain.NewCredentials("tfk_test", "tfs_secret")
	require.NoError(t, err)

	config, err := domain.NewTelemetryConfig(creds, "localhost:4317", "test-service")
	require.NoError(t, err)

	return config
}

func TestNewOTLPExporterFactory(t *testing.T) {
	t.Run("should create factory with config", func(t *testing.T) {
		config := createValidConfig(t)

		factory := NewOTLPExporterFactory(config)

		assert.NotNil(t, factory)
		assert.Equal(t, config, factory.config)
	})
}

func TestOTLPExporterFactory_CreateResource(t *testing.T) {
	t.Run("should create resource with service info", func(t *testing.T) {
		config := createValidConfig(t)
		factory := NewOTLPExporterFactory(config)

		resource, err := factory.CreateResource(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, resource)
	})

	t.Run("should include custom attributes in resource", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithCustomAttribute("custom.key", "custom.value")
		factory := NewOTLPExporterFactory(config)

		resource, err := factory.CreateResource(context.Background())

		require.NoError(t, err)
		assert.NotNil(t, resource)
	})
}

func TestOTLPExporterFactory_CreateTraceExporter(t *testing.T) {
	t.Run("should fail when traces not enabled", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithSignals(true, true, false) // disable traces
		factory := NewOTLPExporterFactory(config)

		_, err := factory.CreateTraceExporter(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enabled")
	})

	t.Run("should fail with unsupported protocol", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithProtocol(domain.Protocol("unsupported"))
		factory := NewOTLPExporterFactory(config)

		_, err := factory.CreateTraceExporter(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported protocol")
	})
}

func TestOTLPExporterFactory_CreateMetricExporter(t *testing.T) {
	t.Run("should fail when metrics not enabled", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithSignals(false, true, true) // disable metrics
		factory := NewOTLPExporterFactory(config)

		_, err := factory.CreateMetricExporter(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not enabled")
	})

	t.Run("should fail with unsupported protocol", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithProtocol(domain.Protocol("unsupported"))
		factory := NewOTLPExporterFactory(config)

		_, err := factory.CreateMetricExporter(context.Background())

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unsupported protocol")
	})
}

func TestOTLPExporterFactory_GetAuthHeaders(t *testing.T) {
	t.Run("should include required headers", func(t *testing.T) {
		config := createValidConfig(t)
		factory := NewOTLPExporterFactory(config)

		headers := factory.getAuthHeaders()

		assert.NotEmpty(t, headers["authorization"])
		assert.Equal(t, "application/x-protobuf", headers["content-type"])
		assert.Equal(t, "tfk_test", headers["X-TelemetryFlow-Key-ID"])
		assert.Equal(t, "tfs_secret", headers["X-TelemetryFlow-Key-Secret"])
	})

	t.Run("should include collector ID when configured", func(t *testing.T) {
		config := createValidConfig(t)
		config.WithCollectorID("collector-123")
		factory := NewOTLPExporterFactory(config)

		headers := factory.getAuthHeaders()

		assert.Equal(t, "collector-123", headers["X-TelemetryFlow-Collector-ID"])
	})

	t.Run("should not include collector ID when not configured", func(t *testing.T) {
		config := createValidConfig(t)
		factory := NewOTLPExporterFactory(config)

		headers := factory.getAuthHeaders()

		_, exists := headers["X-TelemetryFlow-Collector-ID"]
		assert.False(t, exists)
	})
}

func TestOTLPExporterFactory_AuthInterceptor(t *testing.T) {
	t.Run("should create auth interceptor", func(t *testing.T) {
		config := createValidConfig(t)
		factory := NewOTLPExporterFactory(config)

		interceptor := factory.authInterceptor()

		assert.NotNil(t, interceptor)
	})
}
