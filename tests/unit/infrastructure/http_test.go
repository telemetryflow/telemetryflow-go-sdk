// Package infrastructure_test provides unit tests for infrastructure components.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package infrastructure_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test HTTP response handling patterns

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func TestHTTPResponsePatterns(t *testing.T) {
	t.Run("should format success response", func(t *testing.T) {
		data := map[string]string{"id": "123", "name": "Test"}
		response := APIResponse{
			Success: true,
			Data:    data,
		}

		jsonBytes, err := json.Marshal(response)
		require.NoError(t, err)

		var parsed APIResponse
		err = json.Unmarshal(jsonBytes, &parsed)
		require.NoError(t, err)

		assert.True(t, parsed.Success)
		assert.Nil(t, parsed.Error)
	})

	t.Run("should format error response", func(t *testing.T) {
		response := APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: "Email is required",
			},
		}

		jsonBytes, err := json.Marshal(response)
		require.NoError(t, err)

		var parsed APIResponse
		err = json.Unmarshal(jsonBytes, &parsed)
		require.NoError(t, err)

		assert.False(t, parsed.Success)
		assert.NotNil(t, parsed.Error)
		assert.Equal(t, "INVALID_INPUT", parsed.Error.Code)
	})

	t.Run("should handle pagination response", func(t *testing.T) {
		type PaginatedResponse struct {
			APIResponse
			Page       int `json:"page"`
			PerPage    int `json:"per_page"`
			TotalItems int `json:"total_items"`
			TotalPages int `json:"total_pages"`
		}

		response := PaginatedResponse{
			APIResponse: APIResponse{
				Success: true,
				Data:    []string{"item1", "item2"},
			},
			Page:       1,
			PerPage:    10,
			TotalItems: 25,
			TotalPages: 3,
		}

		jsonBytes, err := json.Marshal(response)
		require.NoError(t, err)

		var parsed PaginatedResponse
		err = json.Unmarshal(jsonBytes, &parsed)
		require.NoError(t, err)

		assert.Equal(t, 1, parsed.Page)
		assert.Equal(t, 10, parsed.PerPage)
		assert.Equal(t, 3, parsed.TotalPages)
	})
}

func TestHTTPMiddlewarePatterns(t *testing.T) {
	t.Run("should add CORS headers", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		corsMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				next.ServeHTTP(w, r)
			})
		}

		req := httptest.NewRequest("GET", "/test", nil)
		rec := httptest.NewRecorder()

		corsMiddleware(handler).ServeHTTP(rec, req)

		assert.Equal(t, "*", rec.Header().Get("Access-Control-Allow-Origin"))
		assert.Contains(t, rec.Header().Get("Access-Control-Allow-Methods"), "GET")
	})

	t.Run("should handle OPTIONS preflight", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			w.Write([]byte("OK"))
		})

		req := httptest.NewRequest("OPTIONS", "/test", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("should log request details", func(t *testing.T) {
		var logBuffer bytes.Buffer

		loggingMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logBuffer.WriteString("Method: " + r.Method + ", Path: " + r.URL.Path + "\n")
				next.ServeHTTP(w, r)
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		req := httptest.NewRequest("GET", "/api/users", nil)
		rec := httptest.NewRecorder()

		loggingMiddleware(handler).ServeHTTP(rec, req)

		logOutput := logBuffer.String()
		assert.Contains(t, logOutput, "Method: GET")
		assert.Contains(t, logOutput, "Path: /api/users")
	})

	t.Run("should validate JWT token format", func(t *testing.T) {
		authMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					http.Error(w, "Authorization header required", http.StatusUnauthorized)
					return
				}

				if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
					http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
					return
				}

				next.ServeHTTP(w, r)
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Authorized"))
		})

		t.Run("missing auth header", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/protected", nil)
			rec := httptest.NewRecorder()

			authMiddleware(handler).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})

		t.Run("invalid auth format", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/protected", nil)
			req.Header.Set("Authorization", "Basic abc123")
			rec := httptest.NewRecorder()

			authMiddleware(handler).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusUnauthorized, rec.Code)
		})

		t.Run("valid bearer token", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/protected", nil)
			req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
			rec := httptest.NewRecorder()

			authMiddleware(handler).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	})

	t.Run("should implement rate limiting", func(t *testing.T) {
		requestCount := 0
		limit := 3

		rateLimitMiddleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				if requestCount > limit {
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
			})
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK"))
		})

		// First 3 requests should succeed
		for i := 0; i < 3; i++ {
			req := httptest.NewRequest("GET", "/api/resource", nil)
			rec := httptest.NewRecorder()
			rateLimitMiddleware(handler).ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
		}

		// 4th request should be rate limited
		req := httptest.NewRequest("GET", "/api/resource", nil)
		rec := httptest.NewRecorder()
		rateLimitMiddleware(handler).ServeHTTP(rec, req)
		assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	})
}

func TestHTTPRequestHandling(t *testing.T) {
	t.Run("should parse JSON request body", func(t *testing.T) {
		type CreateRequest struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Failed to read body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var req CreateRequest
			if err := json.Unmarshal(body, &req); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			response := APIResponse{
				Success: true,
				Data:    req,
			}
			json.NewEncoder(w).Encode(response)
		})

		body := `{"name": "John", "email": "john@example.com"}`
		req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response APIResponse
		json.NewDecoder(rec.Body).Decode(&response)
		assert.True(t, response.Success)
	})

	t.Run("should validate required fields", func(t *testing.T) {
		type CreateRequest struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()

			var req CreateRequest
			json.Unmarshal(body, &req)

			if req.Name == "" || req.Email == "" {
				response := APIResponse{
					Success: false,
					Error: &APIError{
						Code:    "VALIDATION_ERROR",
						Message: "Name and email are required",
					},
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(response)
				return
			}

			w.Write([]byte("OK"))
		})

		body := `{"name": "", "email": ""}`
		req := httptest.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should handle query parameters", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			page := r.URL.Query().Get("page")
			limit := r.URL.Query().Get("limit")

			response := map[string]string{
				"page":  page,
				"limit": limit,
			}
			json.NewEncoder(w).Encode(response)
		})

		req := httptest.NewRequest("GET", "/api/users?page=2&limit=10", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		var response map[string]string
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, "2", response["page"])
		assert.Equal(t, "10", response["limit"])
	})

	t.Run("should handle path parameters", func(t *testing.T) {
		// Simulate path parameter extraction
		extractPathParam := func(path, pattern string) string {
			// Simple extraction for /api/users/{id}
			if len(path) > len("/api/users/") {
				return path[len("/api/users/"):]
			}
			return ""
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := extractPathParam(r.URL.Path, "/api/users/{id}")
			response := map[string]string{"id": id}
			json.NewEncoder(w).Encode(response)
		})

		req := httptest.NewRequest("GET", "/api/users/123", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		var response map[string]string
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, "123", response["id"])
	})
}

func TestHealthCheckEndpoint(t *testing.T) {
	t.Run("should return healthy status", func(t *testing.T) {
		type HealthResponse struct {
			Status  string `json:"status"`
			Service string `json:"service"`
			Version string `json:"version"`
		}

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			response := HealthResponse{
				Status:  "healthy",
				Service: "test-service",
				Version: "1.0.0",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		})

		req := httptest.NewRequest("GET", "/health", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response HealthResponse
		json.NewDecoder(rec.Body).Decode(&response)

		assert.Equal(t, "healthy", response.Status)
		assert.Equal(t, "test-service", response.Service)
	})

	t.Run("should return liveness probe", func(t *testing.T) {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		req := httptest.NewRequest("GET", "/health/live", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("should return readiness probe", func(t *testing.T) {
		dbReady := true
		cacheReady := true

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !dbReady || !cacheReady {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("Not Ready"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ready"))
		})

		req := httptest.NewRequest("GET", "/health/ready", nil)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
