// Package mocks provides mock implementations for testing.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// MockLogger is a mock implementation of zap.Logger
type MockLogger struct {
	mock.Mock
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Debug mocks debug logging
func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

// Info mocks info logging
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

// Warn mocks warn logging
func (m *MockLogger) Warn(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

// Error mocks error logging
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	args := []interface{}{msg}
	for _, field := range fields {
		args = append(args, field)
	}
	m.Called(args...)
}

// With mocks creating a child logger with additional fields
func (m *MockLogger) With(fields ...zap.Field) *zap.Logger {
	args := m.Called(fields)
	return args.Get(0).(*zap.Logger)
}

// Check mocks checking if a log level is enabled
func (m *MockLogger) Check(level zapcore.Level, msg string) *zapcore.CheckedEntry {
	args := m.Called(level, msg)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*zapcore.CheckedEntry)
}
