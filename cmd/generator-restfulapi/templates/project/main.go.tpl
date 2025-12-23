// Package main is the entry point for {{.ProjectName}} API.
//
// {{.ServiceName}} - RESTful API with DDD + CQRS Pattern
// Copyright (c) 2024-2026 {{.ProjectName}}. All rights reserved.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.ModulePath}}/internal/infrastructure/config"
	"{{.ModulePath}}/internal/infrastructure/http"
	"{{.ModulePath}}/internal/infrastructure/persistence"
{{- if .EnableTelemetry}}
	"{{.ModulePath}}/telemetry"
{{- end}}
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

{{- if .EnableTelemetry}}
	// Initialize TelemetryFlow
	if err := telemetry.Init(); err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer telemetry.Shutdown()
{{- end}}

	// Initialize database
	db, err := persistence.NewDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create HTTP server
	server := http.NewServer(cfg, db)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Printf("{{.ServiceName}} v{{.ServiceVersion}} started on port %s", cfg.Server.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
