// Package middleware provides HTTP middleware.
package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
{{- if .EnableTelemetry}}
	"{{.ModulePath}}/telemetry/logs"
	"{{.ModulePath}}/telemetry/metrics"
{{- end}}
)

// Logger returns a logging middleware
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)
			res := c.Response()

{{- if .EnableTelemetry}}
			// Record metrics
			metrics.RecordHTTPRequest(req.Method, req.URL.Path, res.Status, duration.Seconds())

			// Log request
			logs.Info("HTTP request", map[string]interface{}{
				"method":     req.Method,
				"path":       req.URL.Path,
				"status":     res.Status,
				"duration":   duration.String(),
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				"user_agent": req.UserAgent(),
				"remote_ip":  c.RealIP(),
			})
{{- end}}

			return err
		}
	}
}
