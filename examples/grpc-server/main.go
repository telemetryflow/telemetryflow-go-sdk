// Package main demonstrates TelemetryFlow SDK integration with a gRPC server.
//
// This example shows:
// - gRPC unary and streaming interceptors for tracing
// - RPC metrics (latency, errors, call counts)
// - Service-level health monitoring
// - Graceful shutdown
//
// Note: This example simulates gRPC patterns without the full gRPC dependency
// to keep the example self-contained. In a real application, you would use
// google.golang.org/grpc with proper interceptors.
//
// Run with:
//
//	export TELEMETRYFLOW_API_KEY_ID=tfk_your_key
//	export TELEMETRYFLOW_API_KEY_SECRET=tfs_your_secret
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Simulated gRPC-like types for demonstration
type (
	// UnaryServerInfo represents information about a unary RPC
	UnaryServerInfo struct {
		FullMethod string
		Service    string
		Method     string
	}

	// StreamServerInfo represents information about a streaming RPC
	StreamServerInfo struct {
		FullMethod     string
		Service        string
		Method         string
		IsClientStream bool
		IsServerStream bool
	}
)

// UserService simulates a gRPC service
type UserService struct {
	client *telemetryflow.Client
}

// OrderService simulates a gRPC service
type OrderService struct {
	client *telemetryflow.Client
}

func main() {
	// Initialize TelemetryFlow client
	var err error
	client, err = telemetryflow.NewBuilder().
		WithAPIKeyFromEnv().
		WithEndpointFromEnv().
		WithService("grpc-server-example", "1.0.0").
		WithEnvironmentFromEnv().
		WithGRPC().
		WithCustomAttribute("example", "grpc-server").
		WithCustomAttribute("grpc.port", "50051").
		Build()

	if err != nil {
		log.Fatalf("Failed to create TelemetryFlow client: %v", err)
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		log.Fatalf("Failed to initialize TelemetryFlow: %v", err)
	}

	if err := client.LogInfo(ctx, "gRPC server starting", map[string]interface{}{
		"port":    50051,
		"version": "1.0.0",
	}); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	// Create services
	userService := &UserService{client: client}
	orderService := &OrderService{client: client}

	// Simulate gRPC server running
	log.Println("gRPC server listening on :50051")

	// Start request simulator
	quit := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		simulateGRPCRequests(ctx, userService, orderService, quit)
	}()

	// Start health reporter
	wg.Add(1)
	go func() {
		defer wg.Done()
		reportGRPCHealth(ctx, quit)
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutdown signal received...")
	if err := client.LogInfo(ctx, "gRPC server shutdown initiated", nil); err != nil {
		log.Printf("Failed to log info: %v", err)
	}

	close(quit)
	wg.Wait()

	// Flush and shutdown telemetry
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Flush(shutdownCtx); err != nil {
		log.Printf("Failed to flush: %v", err)
	}
	if err := client.Shutdown(shutdownCtx); err != nil {
		log.Printf("Failed to shutdown: %v", err)
	}

	log.Println("gRPC server stopped")
}

// UnaryServerInterceptor creates telemetry for unary RPCs
func UnaryServerInterceptor(ctx context.Context, info *UnaryServerInfo, handler func(context.Context) (interface{}, error)) (interface{}, error) {
	start := time.Now()

	// Start span for RPC
	spanID, err := client.StartSpan(ctx, info.FullMethod, "server", map[string]interface{}{
		"rpc.system":           "grpc",
		"rpc.service":          info.Service,
		"rpc.method":           info.Method,
		"rpc.grpc.status_code": 0, // Will be updated
	})
	if err != nil {
		log.Printf("Failed to start span: %v", err)
	}

	// Record RPC started
	if err := client.IncrementCounter(ctx, "grpc.server.started", 1, map[string]interface{}{
		"grpc.service": info.Service,
		"grpc.method":  info.Method,
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	// Call handler
	resp, handlerErr := handler(ctx)

	duration := time.Since(start)

	// Determine status code
	statusCode := "OK"
	if handlerErr != nil {
		statusCode = "INTERNAL"
		if handlerErr.Error() == "not found" {
			statusCode = "NOT_FOUND"
		} else if handlerErr.Error() == "invalid argument" {
			statusCode = "INVALID_ARGUMENT"
		}
	}

	// Record metrics
	if err := client.RecordHistogram(ctx, "grpc.server.duration", duration.Seconds(), "s", map[string]interface{}{
		"grpc.service":     info.Service,
		"grpc.method":      info.Method,
		"grpc.status_code": statusCode,
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	if err := client.IncrementCounter(ctx, "grpc.server.handled", 1, map[string]interface{}{
		"grpc.service":     info.Service,
		"grpc.method":      info.Method,
		"grpc.status_code": statusCode,
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	// Log errors
	if handlerErr != nil {
		if err := client.LogError(ctx, "gRPC call failed", map[string]interface{}{
			"grpc.service":     info.Service,
			"grpc.method":      info.Method,
			"grpc.status_code": statusCode,
			"error":            handlerErr.Error(),
			"duration_ms":      duration.Milliseconds(),
		}); err != nil {
			log.Printf("Failed to log error: %v", err)
		}
	}

	// End span
	if spanID != "" {
		if err := client.EndSpan(ctx, spanID, handlerErr); err != nil {
			log.Printf("Failed to end span: %v", err)
		}
	}

	return resp, handlerErr
}

// StreamServerInterceptor creates telemetry for streaming RPCs
func StreamServerInterceptor(ctx context.Context, info *StreamServerInfo, handler func(context.Context) error) error {
	start := time.Now()

	streamType := "bidirectional"
	if info.IsClientStream && !info.IsServerStream {
		streamType = "client_streaming"
	} else if !info.IsClientStream && info.IsServerStream {
		streamType = "server_streaming"
	}

	// Start span for stream
	spanID, _ := client.StartSpan(ctx, info.FullMethod, "server", map[string]interface{}{
		"rpc.system":       "grpc",
		"rpc.service":      info.Service,
		"rpc.method":       info.Method,
		"grpc.stream_type": streamType,
	})

	if err := client.IncrementCounter(ctx, "grpc.server.stream.started", 1, map[string]interface{}{
		"grpc.service":     info.Service,
		"grpc.method":      info.Method,
		"grpc.stream_type": streamType,
	}); err != nil {
		log.Printf("Failed to increment counter: %v", err)
	}

	// Call handler
	handlerErr := handler(ctx)

	duration := time.Since(start)

	if err := client.RecordHistogram(ctx, "grpc.server.stream.duration", duration.Seconds(), "s", map[string]interface{}{
		"grpc.service":     info.Service,
		"grpc.method":      info.Method,
		"grpc.stream_type": streamType,
	}); err != nil {
		log.Printf("Failed to record histogram: %v", err)
	}

	if spanID != "" {
		if err := client.EndSpan(ctx, spanID, handlerErr); err != nil {
			log.Printf("Failed to end span: %v", err)
		}
	}

	return handlerErr
}

// UserService methods

func (s *UserService) GetUser(ctx context.Context, userID string) (map[string]interface{}, error) {
	info := &UnaryServerInfo{
		FullMethod: "/user.UserService/GetUser",
		Service:    "user.UserService",
		Method:     "GetUser",
	}

	result, err := UnaryServerInterceptor(ctx, info, func(ctx context.Context) (interface{}, error) {
		// Start database span
		dbSpanID, _ := client.StartSpan(ctx, "db.query.user", "client", map[string]interface{}{
			"db.system":    "postgresql",
			"db.operation": "SELECT",
			"db.table":     "users",
		})

		// Simulate database query
		time.Sleep(time.Duration(20+rand.Intn(30)) * time.Millisecond)

		if err := client.EndSpan(ctx, dbSpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Simulate not found
		if rand.Float32() < 0.1 {
			return nil, fmt.Errorf("not found")
		}

		return map[string]interface{}{
			"id":    userID,
			"name":  "John Doe",
			"email": "john@example.com",
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

func (s *UserService) CreateUser(ctx context.Context, user map[string]interface{}) (map[string]interface{}, error) {
	info := &UnaryServerInfo{
		FullMethod: "/user.UserService/CreateUser",
		Service:    "user.UserService",
		Method:     "CreateUser",
	}

	result, err := UnaryServerInterceptor(ctx, info, func(ctx context.Context) (interface{}, error) {
		// Validate
		if user["email"] == nil {
			return nil, fmt.Errorf("invalid argument")
		}

		// Start database span
		dbSpanID, _ := client.StartSpan(ctx, "db.insert.user", "client", map[string]interface{}{
			"db.system":    "postgresql",
			"db.operation": "INSERT",
			"db.table":     "users",
		})

		time.Sleep(time.Duration(30+rand.Intn(50)) * time.Millisecond)
		if err := client.EndSpan(ctx, dbSpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Record business metric
		if err := client.IncrementCounter(ctx, "users.created", 1, map[string]interface{}{
			"source": "grpc",
		}); err != nil {
			log.Printf("Failed to increment counter: %v", err)
		}

		return map[string]interface{}{
			"id":      fmt.Sprintf("user_%d", rand.Intn(10000)),
			"name":    user["name"],
			"email":   user["email"],
			"created": time.Now().Format(time.RFC3339),
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

func (s *UserService) ListUsers(ctx context.Context, pageSize int) error {
	info := &StreamServerInfo{
		FullMethod:     "/user.UserService/ListUsers",
		Service:        "user.UserService",
		Method:         "ListUsers",
		IsServerStream: true,
	}

	return StreamServerInterceptor(ctx, info, func(ctx context.Context) error {
		// Simulate streaming response
		for i := 0; i < pageSize; i++ {
			// Simulate sending each user
			time.Sleep(time.Duration(10+rand.Intn(20)) * time.Millisecond)

			if err := client.IncrementCounter(ctx, "grpc.server.stream.msg_sent", 1, map[string]interface{}{
				"grpc.service": "user.UserService",
				"grpc.method":  "ListUsers",
			}); err != nil {
				log.Printf("Failed to increment counter: %v", err)
			}
		}
		return nil
	})
}

// OrderService methods

func (s *OrderService) GetOrder(ctx context.Context, orderID string) (map[string]interface{}, error) {
	info := &UnaryServerInfo{
		FullMethod: "/order.OrderService/GetOrder",
		Service:    "order.OrderService",
		Method:     "GetOrder",
	}

	result, err := UnaryServerInterceptor(ctx, info, func(ctx context.Context) (interface{}, error) {
		// Query database
		dbSpanID, _ := client.StartSpan(ctx, "db.query.order", "client", map[string]interface{}{
			"db.system":    "postgresql",
			"db.operation": "SELECT",
			"db.table":     "orders",
		})

		time.Sleep(time.Duration(25+rand.Intn(35)) * time.Millisecond)
		if err := client.EndSpan(ctx, dbSpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		return map[string]interface{}{
			"id":     orderID,
			"total":  99.99,
			"status": "pending",
			"items":  3,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order map[string]interface{}) (map[string]interface{}, error) {
	info := &UnaryServerInfo{
		FullMethod: "/order.OrderService/CreateOrder",
		Service:    "order.OrderService",
		Method:     "CreateOrder",
	}

	result, err := UnaryServerInterceptor(ctx, info, func(ctx context.Context) (interface{}, error) {
		// Validate order
		spanID, _ := client.StartSpan(ctx, "order.validate", "internal", nil)
		time.Sleep(10 * time.Millisecond)
		if err := client.EndSpan(ctx, spanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Check inventory (external call)
		invSpanID, _ := client.StartSpan(ctx, "inventory.check", "client", map[string]interface{}{
			"rpc.system":  "grpc",
			"rpc.service": "inventory.InventoryService",
		})
		time.Sleep(time.Duration(30+rand.Intn(40)) * time.Millisecond)
		if err := client.EndSpan(ctx, invSpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Process payment (external call)
		paySpanID, _ := client.StartSpan(ctx, "payment.process", "client", map[string]interface{}{
			"payment.provider": "stripe",
		})
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

		// Simulate occasional payment failures
		if rand.Float32() < 0.05 {
			err := fmt.Errorf("payment declined")
			if err := client.EndSpan(ctx, paySpanID, err); err != nil {
				log.Printf("Failed to end span: %v", err)
			}
			return nil, err
		}
		if err := client.EndSpan(ctx, paySpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Save order
		dbSpanID, _ := client.StartSpan(ctx, "db.insert.order", "client", map[string]interface{}{
			"db.system":    "postgresql",
			"db.operation": "INSERT",
			"db.table":     "orders",
		})
		time.Sleep(time.Duration(20+rand.Intn(30)) * time.Millisecond)
		if err := client.EndSpan(ctx, dbSpanID, nil); err != nil {
			log.Printf("Failed to end span: %v", err)
		}

		// Record business metrics
		total := order["total"].(float64)
		if err := client.IncrementCounter(ctx, "orders.created", 1, map[string]interface{}{
			"source": "grpc",
		}); err != nil {
			log.Printf("Failed to increment counter: %v", err)
		}
		if err := client.RecordHistogram(ctx, "order.value", total, "usd", nil); err != nil {
			log.Printf("Failed to record histogram: %v", err)
		}

		return map[string]interface{}{
			"id":      fmt.Sprintf("ord_%d", rand.Intn(100000)),
			"total":   total,
			"status":  "confirmed",
			"created": time.Now().Format(time.RFC3339),
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return result.(map[string]interface{}), nil
}

// simulateGRPCRequests simulates incoming gRPC requests
func simulateGRPCRequests(ctx context.Context, userSvc *UserService, orderSvc *OrderService, quit chan struct{}) {
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			// Randomly call different services
			switch rand.Intn(5) {
			case 0:
				if _, err := userSvc.GetUser(ctx, fmt.Sprintf("user_%d", rand.Intn(1000))); err != nil {
					log.Printf("GetUser failed: %v", err)
				}
			case 1:
				if _, err := userSvc.CreateUser(ctx, map[string]interface{}{
					"name":  "New User",
					"email": fmt.Sprintf("user%d@example.com", rand.Intn(10000)),
				}); err != nil {
					log.Printf("CreateUser failed: %v", err)
				}
			case 2:
				if err := userSvc.ListUsers(ctx, 5+rand.Intn(10)); err != nil {
					log.Printf("ListUsers failed: %v", err)
				}
			case 3:
				if _, err := orderSvc.GetOrder(ctx, fmt.Sprintf("ord_%d", rand.Intn(10000))); err != nil {
					log.Printf("GetOrder failed: %v", err)
				}
			case 4:
				if _, err := orderSvc.CreateOrder(ctx, map[string]interface{}{
					"total": float64(rand.Intn(500) + 50),
					"items": rand.Intn(5) + 1,
				}); err != nil {
					log.Printf("CreateOrder failed: %v", err)
				}
			}
		}
	}
}

// reportGRPCHealth periodically reports gRPC server health
func reportGRPCHealth(ctx context.Context, quit chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			// Record server health metrics
			if err := client.RecordGauge(ctx, "grpc.server.connections", float64(rand.Intn(50)+10), map[string]interface{}{
				"state": "active",
			}); err != nil {
				log.Printf("Failed to record gauge: %v", err)
			}

			// Simulate memory usage
			memUsage := float64(rand.Intn(100) + 200)
			if err := client.RecordGauge(ctx, "grpc.server.memory", memUsage, map[string]interface{}{
				"unit": "MB",
			}); err != nil {
				log.Printf("Failed to record gauge: %v", err)
			}

			// Log health status
			if err := client.LogInfo(ctx, "gRPC server health check", map[string]interface{}{
				"status":      "healthy",
				"memory_mb":   memUsage,
				"connections": rand.Intn(50) + 10,
			}); err != nil {
				log.Printf("Failed to log info: %v", err)
			}
		}
	}
}
