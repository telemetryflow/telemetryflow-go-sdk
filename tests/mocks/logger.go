// Package mocks provides mock implementations for testing.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform
// Copyright (c) 2024-2026 Telemetri Data Indonesia. All rights reserved.
// Open Source Software built by Telemetri Data Indonesia.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
