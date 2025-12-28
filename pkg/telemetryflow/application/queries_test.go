// Package application provides unit tests for the application layer queries.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package application

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryBus(t *testing.T) {
	t.Run("should create new query bus", func(t *testing.T) {
		bus := NewQueryBus()

		assert.NotNil(t, bus)
		assert.NotNil(t, bus.handlers)
	})

	t.Run("should register handler", func(t *testing.T) {
		bus := NewQueryBus()
		handler := &mockQueryHandler{}

		bus.Register("test", handler)

		assert.Equal(t, handler, bus.handlers["test"])
	})

	t.Run("should dispatch query", func(t *testing.T) {
		bus := NewQueryBus()

		result, err := bus.Dispatch(context.Background(), nil)

		assert.NoError(t, err)
		assert.Nil(t, result)
	})
}

type mockQueryHandler struct{}

func (h *mockQueryHandler) Handle(ctx context.Context, query Query) (interface{}, error) {
	return nil, nil
}

func TestGetMetricQuery(t *testing.T) {
	t.Run("should create metric query", func(t *testing.T) {
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now()

		query := &GetMetricQuery{
			Name:      "cpu.usage",
			StartTime: start,
			EndTime:   end,
			GroupBy:   []string{"host", "region"},
			Filters:   map[string]string{"env": "prod"},
		}

		assert.Equal(t, "cpu.usage", query.Name)
		assert.Equal(t, start, query.StartTime)
		assert.Equal(t, end, query.EndTime)
		assert.Equal(t, []string{"host", "region"}, query.GroupBy)
		assert.Equal(t, "prod", query.Filters["env"])
	})
}

func TestMetricQueryResult(t *testing.T) {
	t.Run("should create metric query result", func(t *testing.T) {
		now := time.Now()
		result := &MetricQueryResult{
			Name: "cpu.usage",
			DataPoints: []MetricDataPoint{
				{
					Timestamp:  now,
					Value:      85.5,
					Attributes: map[string]interface{}{"host": "server1"},
				},
			},
		}

		assert.Equal(t, "cpu.usage", result.Name)
		assert.Len(t, result.DataPoints, 1)
		assert.Equal(t, 85.5, result.DataPoints[0].Value)
	})
}

func TestAggregateMetricsQuery(t *testing.T) {
	t.Run("should create aggregate metrics query", func(t *testing.T) {
		query := &AggregateMetricsQuery{
			Name:        "request.latency",
			Aggregation: "avg",
			GroupBy:     []string{"service"},
		}

		assert.Equal(t, "request.latency", query.Name)
		assert.Equal(t, "avg", query.Aggregation)
		assert.Equal(t, []string{"service"}, query.GroupBy)
	})
}

func TestGetLogsQuery(t *testing.T) {
	t.Run("should create logs query", func(t *testing.T) {
		start := time.Now().Add(-1 * time.Hour)
		end := time.Now()

		query := &GetLogsQuery{
			StartTime:  start,
			EndTime:    end,
			Severity:   []string{"error", "warn"},
			SearchText: "connection refused",
			Limit:      100,
			Offset:     0,
		}

		assert.Equal(t, start, query.StartTime)
		assert.Equal(t, end, query.EndTime)
		assert.Equal(t, []string{"error", "warn"}, query.Severity)
		assert.Equal(t, "connection refused", query.SearchText)
		assert.Equal(t, 100, query.Limit)
	})
}

func TestLogsQueryResult(t *testing.T) {
	t.Run("should create logs query result", func(t *testing.T) {
		now := time.Now()
		result := &LogsQueryResult{
			Logs: []LogEntry{
				{
					Timestamp:  now,
					Severity:   "error",
					Message:    "Connection failed",
					Attributes: map[string]interface{}{"host": "db-1"},
					TraceID:    "trace-123",
					SpanID:     "span-456",
				},
			},
			TotalCount: 50,
			HasMore:    true,
		}

		assert.Len(t, result.Logs, 1)
		assert.Equal(t, "error", result.Logs[0].Severity)
		assert.Equal(t, 50, result.TotalCount)
		assert.True(t, result.HasMore)
	})
}

func TestGetTraceQuery(t *testing.T) {
	t.Run("should create trace query", func(t *testing.T) {
		query := &GetTraceQuery{
			TraceID: "trace-abc-123",
		}

		assert.Equal(t, "trace-abc-123", query.TraceID)
	})
}

func TestTraceQueryResult(t *testing.T) {
	t.Run("should create trace query result", func(t *testing.T) {
		now := time.Now()
		result := &TraceQueryResult{
			TraceID:   "trace-123",
			Duration:  500 * time.Millisecond,
			StartTime: now,
			EndTime:   now.Add(500 * time.Millisecond),
			Spans: []SpanEntry{
				{
					SpanID:    "span-1",
					ParentID:  "",
					Name:      "root-span",
					Kind:      "server",
					StartTime: now,
					Duration:  500 * time.Millisecond,
					Status:    SpanStatus{Code: "OK", Message: ""},
				},
			},
		}

		assert.Equal(t, "trace-123", result.TraceID)
		assert.Equal(t, 500*time.Millisecond, result.Duration)
		assert.Len(t, result.Spans, 1)
		assert.Equal(t, "root-span", result.Spans[0].Name)
	})
}

func TestSearchTracesQuery(t *testing.T) {
	t.Run("should create search traces query", func(t *testing.T) {
		hasError := true
		query := &SearchTracesQuery{
			ServiceName: "api-gateway",
			Operation:   "GET /users",
			MinDuration: 100 * time.Millisecond,
			MaxDuration: 5 * time.Second,
			Tags:        map[string]string{"env": "prod"},
			HasError:    &hasError,
			Limit:       50,
		}

		assert.Equal(t, "api-gateway", query.ServiceName)
		assert.Equal(t, "GET /users", query.Operation)
		assert.Equal(t, 100*time.Millisecond, query.MinDuration)
		assert.True(t, *query.HasError)
		assert.Equal(t, 50, query.Limit)
	})
}

func TestGetHealthQuery(t *testing.T) {
	t.Run("should create health query", func(t *testing.T) {
		query := &GetHealthQuery{}

		assert.NotNil(t, query)
	})
}

func TestHealthQueryResult(t *testing.T) {
	t.Run("should create health query result", func(t *testing.T) {
		now := time.Now()
		result := &HealthQueryResult{
			Status:      "healthy",
			LastSuccess: now,
			LastError:   nil,
			Metrics: HealthMetrics{
				TotalRequests:   1000,
				FailedRequests:  5,
				AverageLatency:  50 * time.Millisecond,
				ConnectionState: "connected",
			},
		}

		assert.Equal(t, "healthy", result.Status)
		assert.Equal(t, int64(1000), result.Metrics.TotalRequests)
		assert.Equal(t, int64(5), result.Metrics.FailedRequests)
	})
}

func TestGetSDKStatusQuery(t *testing.T) {
	t.Run("should create SDK status query", func(t *testing.T) {
		query := &GetSDKStatusQuery{}

		assert.NotNil(t, query)
	})
}

func TestSDKStatusResult(t *testing.T) {
	t.Run("should create SDK status result", func(t *testing.T) {
		now := time.Now()
		result := &SDKStatusResult{
			Initialized:    true,
			Version:        "1.1.0",
			EnabledSignals: []string{"traces", "metrics", "logs"},
			Statistics: SDKStatistics{
				MetricsSent: 1000,
				LogsSent:    500,
				TracesSent:  200,
				ErrorsCount: 2,
				LastFlush:   now,
				QueueSize:   10,
			},
		}

		assert.True(t, result.Initialized)
		assert.Equal(t, "1.1.0", result.Version)
		assert.Len(t, result.EnabledSignals, 3)
		assert.Equal(t, int64(1000), result.Statistics.MetricsSent)
	})
}
