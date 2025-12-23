// =============================================================================
// Test Fixtures - Telemetry Data
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// =============================================================================

package fixtures

import "time"

// MetricFixture represents a metric test fixture
type MetricFixture struct {
	Name       string
	Type       string // counter, gauge, histogram
	Value      float64
	Unit       string
	Attributes map[string]interface{}
}

// SampleMetrics contains sample metric fixtures
var SampleMetrics = []MetricFixture{
	{
		Name:  "http.requests.total",
		Type:  "counter",
		Value: 1,
		Unit:  "",
		Attributes: map[string]interface{}{
			"method":      "GET",
			"path":        "/api/users",
			"status_code": 200,
		},
	},
	{
		Name:  "http.request.duration",
		Type:  "histogram",
		Value: 0.125,
		Unit:  "s",
		Attributes: map[string]interface{}{
			"method": "POST",
			"path":   "/api/orders",
		},
	},
	{
		Name:  "system.memory.usage",
		Type:  "gauge",
		Value: 512.5,
		Unit:  "MB",
		Attributes: map[string]interface{}{
			"host": "server-01",
		},
	},
	{
		Name:  "db.connections.active",
		Type:  "gauge",
		Value: 25,
		Unit:  "",
		Attributes: map[string]interface{}{
			"database": "postgres",
			"pool":     "primary",
		},
	},
	{
		Name:  "cache.hits",
		Type:  "counter",
		Value: 100,
		Unit:  "",
		Attributes: map[string]interface{}{
			"cache_type": "redis",
			"operation":  "get",
		},
	},
	{
		Name:  "queue.messages.processed",
		Type:  "counter",
		Value: 50,
		Unit:  "",
		Attributes: map[string]interface{}{
			"queue":    "orders",
			"consumer": "worker-1",
		},
	},
}

// LogFixture represents a log entry test fixture
type LogFixture struct {
	Level      string
	Message    string
	Attributes map[string]interface{}
}

// SampleLogs contains sample log fixtures
var SampleLogs = []LogFixture{
	{
		Level:   "info",
		Message: "Application started successfully",
		Attributes: map[string]interface{}{
			"version":    "1.0.0",
			"port":       8080,
			"start_time": time.Now().Format(time.RFC3339),
		},
	},
	{
		Level:   "debug",
		Message: "Processing request",
		Attributes: map[string]interface{}{
			"request_id": "req-12345",
			"method":     "GET",
			"path":       "/api/users",
		},
	},
	{
		Level:   "warn",
		Message: "High memory usage detected",
		Attributes: map[string]interface{}{
			"usage_mb":     850,
			"threshold_mb": 800,
			"host":         "server-01",
		},
	},
	{
		Level:   "error",
		Message: "Database connection failed",
		Attributes: map[string]interface{}{
			"error":        "connection timeout",
			"host":         "db.example.com",
			"port":         5432,
			"retry_count":  3,
			"last_attempt": time.Now().Add(-5 * time.Second).Format(time.RFC3339),
		},
	},
	{
		Level:   "info",
		Message: "User authenticated",
		Attributes: map[string]interface{}{
			"user_id":    "user-67890",
			"auth_type":  "jwt",
			"ip_address": "192.168.1.100",
		},
	},
	{
		Level:   "info",
		Message: "Order processed successfully",
		Attributes: map[string]interface{}{
			"order_id": "order-abc123",
			"amount":   99.99,
			"currency": "USD",
			"items":    3,
		},
	},
}

// SpanFixture represents a trace span test fixture
type SpanFixture struct {
	Name       string
	Kind       string // internal, server, client, producer, consumer
	Attributes map[string]interface{}
	Events     []SpanEventFixture
	Duration   time.Duration
	HasError   bool
	ErrorMsg   string
}

// SpanEventFixture represents a span event test fixture
type SpanEventFixture struct {
	Name       string
	Attributes map[string]interface{}
}

// SampleSpans contains sample span fixtures
var SampleSpans = []SpanFixture{
	{
		Name: "http.request",
		Kind: "server",
		Attributes: map[string]interface{}{
			"http.method":      "GET",
			"http.url":         "/api/users/123",
			"http.status_code": 200,
			"http.user_agent":  "Mozilla/5.0",
		},
		Events: []SpanEventFixture{
			{
				Name: "request.received",
				Attributes: map[string]interface{}{
					"content_length": 0,
				},
			},
			{
				Name: "response.sent",
				Attributes: map[string]interface{}{
					"content_length": 256,
				},
			},
		},
		Duration: 125 * time.Millisecond,
		HasError: false,
	},
	{
		Name: "db.query",
		Kind: "client",
		Attributes: map[string]interface{}{
			"db.system":    "postgresql",
			"db.name":      "users",
			"db.operation": "SELECT",
			"db.statement": "SELECT * FROM users WHERE id = $1",
		},
		Events: []SpanEventFixture{
			{
				Name: "query.executed",
				Attributes: map[string]interface{}{
					"rows_affected": 1,
				},
			},
		},
		Duration: 15 * time.Millisecond,
		HasError: false,
	},
	{
		Name: "cache.get",
		Kind: "client",
		Attributes: map[string]interface{}{
			"cache.type": "redis",
			"cache.key":  "user:123",
			"cache.hit":  true,
		},
		Events:   nil,
		Duration: 2 * time.Millisecond,
		HasError: false,
	},
	{
		Name: "process.order",
		Kind: "internal",
		Attributes: map[string]interface{}{
			"order.id":     "order-xyz789",
			"order.amount": 149.99,
			"order.items":  5,
		},
		Events: []SpanEventFixture{
			{
				Name: "validation.complete",
				Attributes: map[string]interface{}{
					"valid": true,
				},
			},
			{
				Name: "payment.processed",
				Attributes: map[string]interface{}{
					"payment_id": "pay-abc123",
				},
			},
		},
		Duration: 500 * time.Millisecond,
		HasError: false,
	},
	{
		Name: "external.api.call",
		Kind: "client",
		Attributes: map[string]interface{}{
			"http.method": "POST",
			"http.url":    "https://api.payment.com/charge",
			"peer.service": "payment-gateway",
		},
		Events:   nil,
		Duration: 1500 * time.Millisecond,
		HasError: true,
		ErrorMsg: "connection timeout",
	},
	{
		Name: "message.publish",
		Kind: "producer",
		Attributes: map[string]interface{}{
			"messaging.system":      "kafka",
			"messaging.destination": "orders",
			"messaging.message_id":  "msg-12345",
		},
		Events:   nil,
		Duration: 10 * time.Millisecond,
		HasError: false,
	},
	{
		Name: "message.consume",
		Kind: "consumer",
		Attributes: map[string]interface{}{
			"messaging.system":      "kafka",
			"messaging.destination": "orders",
			"messaging.message_id":  "msg-12345",
			"messaging.consumer_id": "consumer-1",
		},
		Events:   nil,
		Duration: 50 * time.Millisecond,
		HasError: false,
	},
}

// TraceContext represents trace context for propagation tests
type TraceContext struct {
	TraceID    string
	SpanID     string
	TraceFlags byte
	TraceState string
}

// SampleTraceContexts contains sample trace context fixtures
var SampleTraceContexts = []TraceContext{
	{
		TraceID:    "0af7651916cd43dd8448eb211c80319c",
		SpanID:     "b7ad6b7169203331",
		TraceFlags: 0x01, // sampled
		TraceState: "",
	},
	{
		TraceID:    "80f198ee56343ba864fe8b2a57d3eff7",
		SpanID:     "e457b5a2e4d86bd1",
		TraceFlags: 0x01,
		TraceState: "congo=t61rcWkgMzE",
	},
	{
		TraceID:    "ff000000000000000000000000000041",
		SpanID:     "ff00000000000041",
		TraceFlags: 0x00, // not sampled
		TraceState: "",
	},
}
