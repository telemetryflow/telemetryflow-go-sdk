package telemetry

import (
	"context"
	"log"
	"time"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Init initializes the TelemetryFlow SDK
func Init() error {
	var err error
	client, err = telemetryflow.NewBuilder().
		WithAutoConfiguration().
		WithSignals({{.EnableMetrics}}, {{.EnableLogs}}, {{.EnableTraces}}).
		Build()
	if err != nil {
		return err
	}

	ctx := context.Background()
	if err := client.Initialize(ctx); err != nil {
		return err
	}

	log.Println("TelemetryFlow SDK initialized successfully")
	return nil
}

// Shutdown gracefully shuts down the SDK
func Shutdown() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down TelemetryFlow: %v", err)
		}
	}
}

// Client returns the global TelemetryFlow client
func Client() *telemetryflow.Client {
	return client
}
