package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"{{.ModulePath}}/telemetry"
	"{{.ModulePath}}/telemetry/logs"
	"{{.ModulePath}}/telemetry/metrics"
	"{{.ModulePath}}/telemetry/traces"
)

// Example gRPC server with TelemetryFlow instrumentation

func main() {
	// Initialize TelemetryFlow SDK
	if err := telemetry.Init(); err != nil {
		log.Fatal(err)
	}
	defer telemetry.Shutdown()

	// Create gRPC server with telemetry interceptors
	server := grpc.NewServer(
		grpc.UnaryInterceptor(unaryTelemetryInterceptor),
		grpc.StreamInterceptor(streamTelemetryInterceptor),
	)

	// Register your services here
	// pb.RegisterYourServiceServer(server, &yourServiceImpl{})

	listener, err := net.Listen("tcp", ":{{.Port}}")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	logs.Info("gRPC server starting", map[string]interface{}{
		"port": "{{.Port}}",
	})

	log.Printf("gRPC server listening on :{{.Port}}")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// unaryTelemetryInterceptor adds telemetry to unary RPCs
func unaryTelemetryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Extract metadata for tracing
	md, _ := metadata.FromIncomingContext(ctx)

	// Start a trace span
	spanID, err := traces.StartSpan(ctx, fmt.Sprintf("grpc.unary.%s", info.FullMethod), map[string]interface{}{
		"grpc.method":      info.FullMethod,
		"grpc.request_id":  getRequestID(md),
		"grpc.method_type": "unary",
	})
	if err == nil {
		defer func() {
			traces.EndSpan(ctx, spanID, nil)
		}()
	}

	// Call the handler
	resp, err := handler(ctx, req)

	// Record metrics and logs
	duration := time.Since(start).Seconds()
	grpcCode := status.Code(err)

	metrics.RecordHistogram("grpc.server.duration", duration, "s", map[string]interface{}{
		"method": info.FullMethod,
		"code":   grpcCode.String(),
	})

	metrics.IncrementCounter("grpc.server.requests.total", 1, map[string]interface{}{
		"method": info.FullMethod,
		"code":   grpcCode.String(),
	})

	if err != nil {
		logs.Error("gRPC request failed", map[string]interface{}{
			"method":     info.FullMethod,
			"duration_s": duration,
			"code":       grpcCode.String(),
			"error":      err.Error(),
		})
	} else {
		logs.Info("gRPC request completed", map[string]interface{}{
			"method":     info.FullMethod,
			"duration_s": duration,
			"code":       grpcCode.String(),
		})
	}

	return resp, err
}

// streamTelemetryInterceptor adds telemetry to streaming RPCs
func streamTelemetryInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	start := time.Now()
	ctx := ss.Context()

	// Extract metadata
	md, _ := metadata.FromIncomingContext(ctx)

	// Determine stream type
	streamType := "bidirectional"
	if info.IsClientStream && !info.IsServerStream {
		streamType = "client_stream"
	} else if !info.IsClientStream && info.IsServerStream {
		streamType = "server_stream"
	}

	// Start a trace span
	spanID, err := traces.StartSpan(ctx, fmt.Sprintf("grpc.stream.%s", info.FullMethod), map[string]interface{}{
		"grpc.method":      info.FullMethod,
		"grpc.stream_type": streamType,
		"grpc.request_id":  getRequestID(md),
	})
	if err == nil {
		defer func() {
			traces.EndSpan(ctx, spanID, nil)
		}()
	}

	// Wrap the stream to count messages
	wrapped := &telemetryServerStream{
		ServerStream: ss,
		method:       info.FullMethod,
		recvCount:    0,
		sendCount:    0,
	}

	// Call the handler
	handlerErr := handler(srv, wrapped)

	// Record metrics
	duration := time.Since(start).Seconds()
	grpcCode := status.Code(handlerErr)

	metrics.RecordHistogram("grpc.server.stream.duration", duration, "s", map[string]interface{}{
		"method":      info.FullMethod,
		"stream_type": streamType,
		"code":        grpcCode.String(),
	})

	metrics.IncrementCounter("grpc.server.streams.total", 1, map[string]interface{}{
		"method":      info.FullMethod,
		"stream_type": streamType,
		"code":        grpcCode.String(),
	})

	metrics.RecordGauge("grpc.server.stream.messages.received", float64(wrapped.recvCount), map[string]interface{}{
		"method": info.FullMethod,
	})

	metrics.RecordGauge("grpc.server.stream.messages.sent", float64(wrapped.sendCount), map[string]interface{}{
		"method": info.FullMethod,
	})

	if handlerErr != nil {
		logs.Error("gRPC stream failed", map[string]interface{}{
			"method":        info.FullMethod,
			"stream_type":   streamType,
			"duration_s":    duration,
			"code":          grpcCode.String(),
			"msgs_received": wrapped.recvCount,
			"msgs_sent":     wrapped.sendCount,
			"error":         handlerErr.Error(),
		})
	} else {
		logs.Info("gRPC stream completed", map[string]interface{}{
			"method":        info.FullMethod,
			"stream_type":   streamType,
			"duration_s":    duration,
			"code":          grpcCode.String(),
			"msgs_received": wrapped.recvCount,
			"msgs_sent":     wrapped.sendCount,
		})
	}

	return handlerErr
}

// telemetryServerStream wraps grpc.ServerStream to count messages
type telemetryServerStream struct {
	grpc.ServerStream
	method    string
	recvCount int
	sendCount int
}

func (s *telemetryServerStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil {
		s.recvCount++
		metrics.IncrementCounter("grpc.server.stream.messages.received.total", 1, map[string]interface{}{
			"method": s.method,
		})
	}
	return err
}

func (s *telemetryServerStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil {
		s.sendCount++
		metrics.IncrementCounter("grpc.server.stream.messages.sent.total", 1, map[string]interface{}{
			"method": s.method,
		})
	}
	return err
}

// getRequestID extracts request ID from metadata
func getRequestID(md metadata.MD) string {
	if values := md.Get("x-request-id"); len(values) > 0 {
		return values[0]
	}
	return ""
}

// Example service implementation (uncomment and modify for your use case)
/*
type yourServiceImpl struct {
	pb.UnimplementedYourServiceServer
}

func (s *yourServiceImpl) YourMethod(ctx context.Context, req *pb.YourRequest) (*pb.YourResponse, error) {
	// Your implementation here
	return &pb.YourResponse{}, nil
}
*/
