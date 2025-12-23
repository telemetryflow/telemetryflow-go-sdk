// Package mocks provides mock implementations for testing.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package mocks

import (
	"context"
	"sync"

	"github.com/stretchr/testify/mock"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow/application"
)

// MockCommandHandler is a mock implementation of the CommandHandler interface
type MockCommandHandler struct {
	mock.Mock
	mu sync.RWMutex

	// Track calls for assertions
	Commands       []application.Command
	initialized    bool
	metricsCount   int64
	logsCount      int64
	tracesCount    int64
	flushCount     int
	shutdownCalled bool
}

// NewMockCommandHandler creates a new mock command handler
func NewMockCommandHandler() *MockCommandHandler {
	return &MockCommandHandler{
		Commands: make([]application.Command, 0),
	}
}

// Handle mocks handling a command
func (m *MockCommandHandler) Handle(ctx context.Context, cmd application.Command) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Commands = append(m.Commands, cmd)

	// Track by command type
	switch cmd.(type) {
	case *application.InitializeSDKCommand:
		m.initialized = true
	case *application.ShutdownSDKCommand:
		m.shutdownCalled = true
		m.initialized = false
	case *application.FlushTelemetryCommand:
		m.flushCount++
	case *application.RecordMetricCommand, *application.RecordCounterCommand,
		*application.RecordGaugeCommand, *application.RecordHistogramCommand:
		m.metricsCount++
	case *application.EmitLogCommand:
		m.logsCount++
	case *application.StartSpanCommand, *application.EndSpanCommand, *application.AddSpanEventCommand:
		m.tracesCount++
	}

	args := m.Called(ctx, cmd)
	return args.Error(0)
}

// IsInitialized returns whether the handler is initialized
func (m *MockCommandHandler) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.initialized
}

// MetricsCount returns the number of metric commands handled
func (m *MockCommandHandler) MetricsCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.metricsCount
}

// LogsCount returns the number of log commands handled
func (m *MockCommandHandler) LogsCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.logsCount
}

// TracesCount returns the number of trace commands handled
func (m *MockCommandHandler) TracesCount() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.tracesCount
}

// FlushCount returns the number of flush commands handled
func (m *MockCommandHandler) FlushCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.flushCount
}

// ShutdownCalled returns whether shutdown was called
func (m *MockCommandHandler) ShutdownCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shutdownCalled
}

// Reset clears all tracked data
func (m *MockCommandHandler) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Commands = make([]application.Command, 0)
	m.initialized = false
	m.metricsCount = 0
	m.logsCount = 0
	m.tracesCount = 0
	m.flushCount = 0
	m.shutdownCalled = false
}

// GetCommandsByType returns all commands of a specific type
func (m *MockCommandHandler) GetCommandsByType(cmdType string) []application.Command {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []application.Command
	for _, cmd := range m.Commands {
		switch cmdType {
		case "metric":
			if _, ok := cmd.(*application.RecordMetricCommand); ok {
				result = append(result, cmd)
			}
		case "counter":
			if _, ok := cmd.(*application.RecordCounterCommand); ok {
				result = append(result, cmd)
			}
		case "gauge":
			if _, ok := cmd.(*application.RecordGaugeCommand); ok {
				result = append(result, cmd)
			}
		case "histogram":
			if _, ok := cmd.(*application.RecordHistogramCommand); ok {
				result = append(result, cmd)
			}
		case "log":
			if _, ok := cmd.(*application.EmitLogCommand); ok {
				result = append(result, cmd)
			}
		case "span_start":
			if _, ok := cmd.(*application.StartSpanCommand); ok {
				result = append(result, cmd)
			}
		case "span_end":
			if _, ok := cmd.(*application.EndSpanCommand); ok {
				result = append(result, cmd)
			}
		}
	}
	return result
}
