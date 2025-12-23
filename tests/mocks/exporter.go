// Package mocks provides mock implementations for testing.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package mocks

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"
)

// ExportedData represents data that was exported
type ExportedData struct {
	Type       string // metrics, logs, traces
	Data       interface{}
	Attributes map[string]interface{}
}

// MockExporter is a mock implementation of an OTLP exporter
type MockExporter struct {
	mock.Mock
	mu sync.RWMutex

	name         string
	running      bool
	exportedData []ExportedData
	flushCount   int
}

// NewMockExporter creates a new mock exporter
func NewMockExporter(name string) *MockExporter {
	return &MockExporter{
		name:         name,
		exportedData: make([]ExportedData, 0),
	}
}

// Name returns the exporter name
func (m *MockExporter) Name() string {
	return m.name
}

// Start mocks starting the exporter
func (m *MockExporter) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	args := m.Called(ctx)
	if args.Error(0) == nil {
		m.running = true
	}
	return args.Error(0)
}

// Stop mocks stopping the exporter
func (m *MockExporter) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	args := m.Called(ctx)
	m.running = false
	return args.Error(0)
}

// ExportMetrics mocks exporting metrics
func (m *MockExporter) ExportMetrics(ctx context.Context, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exportedData = append(m.exportedData, ExportedData{
		Type: "metrics",
		Data: data,
	})

	args := m.Called(ctx, data)
	return args.Error(0)
}

// ExportLogs mocks exporting logs
func (m *MockExporter) ExportLogs(ctx context.Context, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exportedData = append(m.exportedData, ExportedData{
		Type: "logs",
		Data: data,
	})

	args := m.Called(ctx, data)
	return args.Error(0)
}

// ExportTraces mocks exporting traces
func (m *MockExporter) ExportTraces(ctx context.Context, data interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.exportedData = append(m.exportedData, ExportedData{
		Type: "traces",
		Data: data,
	})

	args := m.Called(ctx, data)
	return args.Error(0)
}

// Flush mocks flushing pending data
func (m *MockExporter) Flush(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.flushCount++

	args := m.Called(ctx)
	return args.Error(0)
}

// IsRunning returns whether the exporter is running
func (m *MockExporter) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.running
}

// GetExportedData returns all exported data
func (m *MockExporter) GetExportedData() []ExportedData {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.exportedData
}

// GetExportedDataByType returns exported data filtered by type
func (m *MockExporter) GetExportedDataByType(dataType string) []ExportedData {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []ExportedData
	for _, d := range m.exportedData {
		if d.Type == dataType {
			result = append(result, d)
		}
	}
	return result
}

// FlushCount returns the number of flush calls
func (m *MockExporter) FlushCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.flushCount
}

// Reset clears all tracked data
func (m *MockExporter) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.exportedData = make([]ExportedData, 0)
	m.flushCount = 0
}
