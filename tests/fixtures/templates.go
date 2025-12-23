// =============================================================================
// Test Fixtures - Templates
// =============================================================================
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
//
// =============================================================================

package fixtures

// TemplateData represents template data for generator tests
type TemplateData struct {
	ProjectName    string
	ModulePath     string
	ServiceName    string
	ServiceVersion string
	Environment    string
	APIKeyID       string
	APIKeySecret   string
	Endpoint       string
	Port           string
	EnableMetrics  bool
	EnableLogs     bool
	EnableTraces   bool
	NumWorkers     int
	QueueSize      int
	CustomAttrs    map[string]string
}

// SampleTemplateData contains sample template data for tests
var SampleTemplateData = []TemplateData{
	{
		ProjectName:    "my-api",
		ModulePath:     "github.com/example/my-api",
		ServiceName:    "my-api-service",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		APIKeyID:       "tfk_test_key",
		APIKeySecret:   "tfs_test_secret",
		Endpoint:       "localhost:4317",
		Port:           "8080",
		EnableMetrics:  true,
		EnableLogs:     true,
		EnableTraces:   true,
		NumWorkers:     4,
		QueueSize:      100,
		CustomAttrs: map[string]string{
			"team": "backend",
		},
	},
	{
		ProjectName:    "order-service",
		ModulePath:     "github.com/myorg/order-service",
		ServiceName:    "order-api",
		ServiceVersion: "2.1.0",
		Environment:    "production",
		APIKeyID:       "tfk_prod_key",
		APIKeySecret:   "tfs_prod_secret",
		Endpoint:       "api.telemetryflow.id:4317",
		Port:           "3000",
		EnableMetrics:  true,
		EnableLogs:     true,
		EnableTraces:   true,
		NumWorkers:     8,
		QueueSize:      500,
		CustomAttrs: map[string]string{
			"team":       "orders",
			"region":     "us-east-1",
			"datacenter": "dc1",
		},
	},
}

// SampleInitTemplate is a sample init.go template for testing
var SampleInitTemplate = `package telemetry

import (
    "context"
    "log"

    "github.com/telemetryflow/telemetryflow-go-sdk/pkg/telemetryflow"
)

var client *telemetryflow.Client

// Init initializes TelemetryFlow SDK
func Init() error {
    var err error
    client, err = telemetryflow.NewBuilder().
        WithAPIKey("{{.APIKeyID}}", "{{.APIKeySecret}}").
        WithEndpoint("{{.Endpoint}}").
        WithService("{{.ServiceName}}", "{{.ServiceVersion}}").
        WithEnvironment("{{.Environment}}").
        WithSignals({{.EnableMetrics}}, {{.EnableLogs}}, {{.EnableTraces}}).
        Build()
    if err != nil {
        return err
    }

    ctx := context.Background()
    return client.Initialize(ctx)
}

// Shutdown gracefully shuts down the SDK
func Shutdown() {
    if client != nil {
        ctx := context.Background()
        client.Shutdown(ctx)
    }
}

// Client returns the TelemetryFlow client
func Client() *telemetryflow.Client {
    return client
}
`

// SampleMetricsTemplate is a sample metrics.go template for testing
var SampleMetricsTemplate = `package metrics

import (
    "context"

    "../"
)

// RecordCounter records a counter metric
func RecordCounter(name string, value int64, attrs map[string]string) {
    client := telemetry.Client()
    if client != nil {
        ctx := context.Background()
        attrMap := make(map[string]interface{})
        for k, v := range attrs {
            attrMap[k] = v
        }
        client.IncrementCounter(ctx, name, value, attrMap)
    }
}

// RecordGauge records a gauge metric
func RecordGauge(name string, value float64, attrs map[string]string) {
    client := telemetry.Client()
    if client != nil {
        ctx := context.Background()
        attrMap := make(map[string]interface{})
        for k, v := range attrs {
            attrMap[k] = v
        }
        client.RecordGauge(ctx, name, value, attrMap)
    }
}

// RecordHistogram records a histogram metric
func RecordHistogram(name string, value float64, unit string, attrs map[string]string) {
    client := telemetry.Client()
    if client != nil {
        ctx := context.Background()
        attrMap := make(map[string]interface{})
        for k, v := range attrs {
            attrMap[k] = v
        }
        client.RecordHistogram(ctx, name, value, unit, attrMap)
    }
}
`

// EntityTemplateData represents entity template data for RESTful API generator
type EntityTemplateData struct {
	EntityName       string
	EntityNameLower  string
	EntityNamePlural string
	ModulePath       string
	Fields           []EntityField
}

// EntityField represents an entity field
type EntityField struct {
	Name     string
	Type     string
	JSONName string
	DBColumn string
	GoType   string
	Nullable bool
}

// SampleEntityData contains sample entity data for tests
var SampleEntityData = []EntityTemplateData{
	{
		EntityName:       "User",
		EntityNameLower:  "user",
		EntityNamePlural: "users",
		ModulePath:       "github.com/example/my-api",
		Fields: []EntityField{
			{Name: "Name", Type: "string", JSONName: "name", DBColumn: "name", GoType: "string", Nullable: false},
			{Name: "Email", Type: "string", JSONName: "email", DBColumn: "email", GoType: "string", Nullable: false},
			{Name: "Age", Type: "int", JSONName: "age", DBColumn: "age", GoType: "int", Nullable: true},
			{Name: "Active", Type: "bool", JSONName: "active", DBColumn: "active", GoType: "bool", Nullable: false},
		},
	},
	{
		EntityName:       "Product",
		EntityNameLower:  "product",
		EntityNamePlural: "products",
		ModulePath:       "github.com/example/my-api",
		Fields: []EntityField{
			{Name: "Name", Type: "string", JSONName: "name", DBColumn: "name", GoType: "string", Nullable: false},
			{Name: "Description", Type: "text", JSONName: "description", DBColumn: "description", GoType: "string", Nullable: true},
			{Name: "Price", Type: "decimal", JSONName: "price", DBColumn: "price", GoType: "float64", Nullable: false},
			{Name: "Stock", Type: "int", JSONName: "stock", DBColumn: "stock", GoType: "int", Nullable: false},
		},
	},
	{
		EntityName:       "Order",
		EntityNameLower:  "order",
		EntityNamePlural: "orders",
		ModulePath:       "github.com/example/my-api",
		Fields: []EntityField{
			{Name: "CustomerID", Type: "uuid", JSONName: "customer_id", DBColumn: "customer_id", GoType: "uuid.UUID", Nullable: false},
			{Name: "Total", Type: "decimal", JSONName: "total", DBColumn: "total", GoType: "float64", Nullable: false},
			{Name: "Status", Type: "string", JSONName: "status", DBColumn: "status", GoType: "string", Nullable: false},
			{Name: "CreatedAt", Type: "timestamp", JSONName: "created_at", DBColumn: "created_at", GoType: "time.Time", Nullable: false},
		},
	},
}
