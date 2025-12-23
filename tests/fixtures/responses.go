// =============================================================================
// Test Fixtures - API Responses
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// =============================================================================

package fixtures

// HTTPResponse represents an HTTP response fixture
type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// HealthResponses contains health check response fixtures
var HealthResponses = map[string]HTTPResponse{
	"healthy": {
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "status": "healthy",
  "service": "telemetryflow-backend",
  "version": "1.0.0",
  "timestamp": "2024-01-01T00:00:00Z"
}`,
	},
	"unhealthy": {
		StatusCode: 503,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "status": "unhealthy",
  "service": "telemetryflow-backend",
  "error": "database connection failed",
  "timestamp": "2024-01-01T00:00:00Z"
}`,
	},
	"degraded": {
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "status": "degraded",
  "service": "telemetryflow-backend",
  "message": "high latency detected",
  "timestamp": "2024-01-01T00:00:00Z"
}`,
	},
}

// AuthResponses contains authentication response fixtures
var AuthResponses = map[string]HTTPResponse{
	"valid": {
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "valid": true,
  "key_id": "tfk_test_key",
  "permissions": ["write:telemetry", "read:telemetry"],
  "expires_at": "2025-01-01T00:00:00Z"
}`,
	},
	"invalid": {
		StatusCode: 401,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "valid": false,
  "error": "invalid API key",
  "code": "INVALID_API_KEY"
}`,
	},
	"expired": {
		StatusCode: 401,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "valid": false,
  "error": "API key has expired",
  "code": "EXPIRED_API_KEY"
}`,
	},
	"rate_limited": {
		StatusCode: 429,
		Headers: map[string]string{
			"Content-Type":          "application/json",
			"Retry-After":           "60",
			"X-RateLimit-Limit":     "1000",
			"X-RateLimit-Remaining": "0",
			"X-RateLimit-Reset":     "1704067260",
		},
		Body: `{
  "error": "rate limit exceeded",
  "code": "RATE_LIMIT_EXCEEDED",
  "retry_after": 60
}`,
	},
}

// TelemetryIngestResponses contains telemetry ingest response fixtures
var TelemetryIngestResponses = map[string]HTTPResponse{
	"accepted": {
		StatusCode: 202,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "accepted": true,
  "batch_id": "batch-12345-abcde",
  "items_received": 100,
  "timestamp": "2024-01-01T00:00:00Z"
}`,
	},
	"partial": {
		StatusCode: 207,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "accepted": true,
  "batch_id": "batch-12345-abcde",
  "items_received": 100,
  "items_rejected": 5,
  "errors": [
    {"index": 10, "error": "invalid metric name"},
    {"index": 25, "error": "missing required field"},
    {"index": 30, "error": "invalid timestamp"},
    {"index": 45, "error": "attribute value too long"},
    {"index": 80, "error": "unknown metric type"}
  ],
  "timestamp": "2024-01-01T00:00:00Z"
}`,
	},
	"invalid_payload": {
		StatusCode: 400,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "invalid payload",
  "code": "INVALID_PAYLOAD",
  "details": "failed to parse JSON body"
}`,
	},
	"service_unavailable": {
		StatusCode: 503,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Retry-After":  "30",
		},
		Body: `{
  "error": "service temporarily unavailable",
  "code": "SERVICE_UNAVAILABLE",
  "retry_after": 30
}`,
	},
}

// ErrorResponses contains common error response fixtures
var ErrorResponses = map[string]HTTPResponse{
	"bad_request": {
		StatusCode: 400,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "bad request",
  "code": "BAD_REQUEST",
  "message": "The request was invalid"
}`,
	},
	"unauthorized": {
		StatusCode: 401,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "unauthorized",
  "code": "UNAUTHORIZED",
  "message": "Authentication required"
}`,
	},
	"forbidden": {
		StatusCode: 403,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "forbidden",
  "code": "FORBIDDEN",
  "message": "Access denied"
}`,
	},
	"not_found": {
		StatusCode: 404,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "not found",
  "code": "NOT_FOUND",
  "message": "Resource not found"
}`,
	},
	"internal_server_error": {
		StatusCode: 500,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "internal server error",
  "code": "INTERNAL_ERROR",
  "message": "An unexpected error occurred"
}`,
	},
	"gateway_timeout": {
		StatusCode: 504,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "error": "gateway timeout",
  "code": "GATEWAY_TIMEOUT",
  "message": "Upstream server did not respond in time"
}`,
	},
}

// OTLPExportResponses contains OTLP export response fixtures
var OTLPExportResponses = map[string]HTTPResponse{
	"success": {
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/x-protobuf",
		},
		Body: "", // Empty body for successful OTLP export
	},
	"partial_success": {
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/x-protobuf",
		},
		Body: "", // Would contain ExportMetricsPartialSuccess protobuf
	},
	"invalid_argument": {
		StatusCode: 400,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "code": 3,
  "message": "invalid argument: missing required field"
}`,
	},
	"unauthenticated": {
		StatusCode: 401,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: `{
  "code": 16,
  "message": "unauthenticated: invalid credentials"
}`,
	},
	"resource_exhausted": {
		StatusCode: 429,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Retry-After":  "60",
		},
		Body: `{
  "code": 8,
  "message": "resource exhausted: rate limit exceeded"
}`,
	},
	"unavailable": {
		StatusCode: 503,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"Retry-After":  "30",
		},
		Body: `{
  "code": 14,
  "message": "unavailable: service temporarily unavailable"
}`,
	},
}
