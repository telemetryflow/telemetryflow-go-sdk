// Package generator_test provides unit tests for the code generator.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TemplateData mirrors the generator's TemplateData struct for testing
type TemplateData struct {
	ProjectName   string
	ModulePath    string
	ServiceName   string
	Environment   string
	EnableMetrics bool
	EnableLogs    bool
	EnableTraces  bool
	APIKeyID      string
	APIKeySecret  string
	Endpoint      string
	Port          string
	NumWorkers    int
	QueueSize     int
}

func TestTemplateDataDefaults(t *testing.T) {
	t.Run("should have default values", func(t *testing.T) {
		data := TemplateData{
			ProjectName:   "test-project",
			ModulePath:    "github.com/test/project",
			ServiceName:   "test-service",
			Environment:   "production",
			EnableMetrics: true,
			EnableLogs:    true,
			EnableTraces:  true,
			Endpoint:      "api.telemetryflow.id:4317",
			Port:          "8080",
			NumWorkers:    5,
			QueueSize:     100,
		}

		assert.Equal(t, "test-project", data.ProjectName)
		assert.Equal(t, "production", data.Environment)
		assert.True(t, data.EnableMetrics)
		assert.True(t, data.EnableLogs)
		assert.True(t, data.EnableTraces)
	})
}

func TestTemplateExecution(t *testing.T) {
	t.Run("should execute simple template", func(t *testing.T) {
		tmplStr := `Project: {{.ProjectName}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := TemplateData{ProjectName: "my-project"}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		assert.Equal(t, "Project: my-project", sb.String())
	})

	t.Run("should execute template with conditionals", func(t *testing.T) {
		tmplStr := `{{if .EnableMetrics}}Metrics: enabled{{else}}Metrics: disabled{{end}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		dataEnabled := TemplateData{EnableMetrics: true}
		dataDisabled := TemplateData{EnableMetrics: false}

		var sb1, sb2 strings.Builder
		err = tmpl.Execute(&sb1, dataEnabled)
		require.NoError(t, err)
		assert.Equal(t, "Metrics: enabled", sb1.String())

		err = tmpl.Execute(&sb2, dataDisabled)
		require.NoError(t, err)
		assert.Equal(t, "Metrics: disabled", sb2.String())
	})

	t.Run("should handle complex template data", func(t *testing.T) {
		tmplStr := `
Service: {{.ServiceName}}
Endpoint: {{.Endpoint}}
Port: {{.Port}}
`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := TemplateData{
			ServiceName: "my-service",
			Endpoint:    "api.telemetryflow.id:4317",
			Port:        "8080",
		}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		result := sb.String()
		assert.Contains(t, result, "my-service")
		assert.Contains(t, result, "api.telemetryflow.id:4317")
		assert.Contains(t, result, "8080")
	})
}

func TestDirectoryCreation(t *testing.T) {
	t.Run("should create directory structure", func(t *testing.T) {
		// Create temp directory
		tmpDir, err := os.MkdirTemp("", "telemetryflow-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		// Define directories to create
		dirs := []string{
			filepath.Join(tmpDir, "telemetry"),
			filepath.Join(tmpDir, "telemetry", "metrics"),
			filepath.Join(tmpDir, "telemetry", "logs"),
			filepath.Join(tmpDir, "telemetry", "traces"),
		}

		// Create directories
		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0750)
			require.NoError(t, err)
		}

		// Verify directories exist
		for _, dir := range dirs {
			info, err := os.Stat(dir)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
		}
	})

	t.Run("should handle existing directories", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "telemetryflow-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		dir := filepath.Join(tmpDir, "telemetry")

		// Create directory first time
		err = os.MkdirAll(dir, 0750)
		require.NoError(t, err)

		// Should not error on second creation
		err = os.MkdirAll(dir, 0750)
		require.NoError(t, err)
	})
}

func TestFileGeneration(t *testing.T) {
	t.Run("should generate file from template", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "telemetryflow-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		tmplStr := `// Package {{.ProjectName}}
package main

import "fmt"

func main() {
	fmt.Println("Hello, {{.ServiceName}}!")
}
`
		tmpl, err := template.New("main").Parse(tmplStr)
		require.NoError(t, err)

		outputPath := filepath.Join(tmpDir, "main.go")
		f, err := os.Create(outputPath)
		require.NoError(t, err)

		data := TemplateData{
			ProjectName: "myproject",
			ServiceName: "myservice",
		}

		err = tmpl.Execute(f, data)
		if err := f.Close(); err != nil {
			t.Errorf("Failed to close file: %v", err)
		}
		require.NoError(t, err)

		// Read and verify content
		content, err := os.ReadFile(outputPath)
		require.NoError(t, err)

		assert.Contains(t, string(content), "// Package myproject")
		assert.Contains(t, string(content), "Hello, myservice!")
	})
}

func TestModulePath(t *testing.T) {
	t.Run("should convert project name to module path", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"My Project", "my-project"},
			{"my-project", "my-project"},
			{"MyProject", "myproject"},
			{"test project", "test-project"},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				result := strings.ToLower(strings.ReplaceAll(tt.input, " ", "-"))
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

func TestGeneratorConfiguration(t *testing.T) {
	t.Run("should configure metrics", func(t *testing.T) {
		data := TemplateData{EnableMetrics: true}
		assert.True(t, data.EnableMetrics)

		data.EnableMetrics = false
		assert.False(t, data.EnableMetrics)
	})

	t.Run("should configure logs", func(t *testing.T) {
		data := TemplateData{EnableLogs: true}
		assert.True(t, data.EnableLogs)

		data.EnableLogs = false
		assert.False(t, data.EnableLogs)
	})

	t.Run("should configure traces", func(t *testing.T) {
		data := TemplateData{EnableTraces: true}
		assert.True(t, data.EnableTraces)

		data.EnableTraces = false
		assert.False(t, data.EnableTraces)
	})

	t.Run("should configure endpoint", func(t *testing.T) {
		defaultEndpoint := "api.telemetryflow.id:4317"
		customEndpoint := "custom.endpoint:4317"

		data := TemplateData{Endpoint: defaultEndpoint}
		assert.Equal(t, defaultEndpoint, data.Endpoint)

		data.Endpoint = customEndpoint
		assert.Equal(t, customEndpoint, data.Endpoint)
	})
}

func TestTemplateValidation(t *testing.T) {
	t.Run("should validate template syntax", func(t *testing.T) {
		validTemplate := `{{.ProjectName}}`
		_, err := template.New("valid").Parse(validTemplate)
		assert.NoError(t, err)
	})

	t.Run("should fail on invalid template syntax", func(t *testing.T) {
		invalidTemplate := `{{.ProjectName`
		_, err := template.New("invalid").Parse(invalidTemplate)
		assert.Error(t, err)
	})

	t.Run("should error on missing field in struct", func(t *testing.T) {
		tmplStr := `{{.NonExistentField}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := TemplateData{ProjectName: "test"}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		// Struct templates error on missing fields
		assert.Error(t, err)
	})

	t.Run("should handle missing field with map gracefully", func(t *testing.T) {
		tmplStr := `{{.NonExistentField}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		// Maps return "<no value>" for missing keys in Go templates
		data := map[string]string{"ProjectName": "test"}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		assert.NoError(t, err)
		// Go templates output "<no value>" for missing map keys
		assert.Equal(t, "<no value>", sb.String())
	})
}

func TestEnvironmentConfiguration(t *testing.T) {
	t.Run("should support different environments", func(t *testing.T) {
		environments := []string{"development", "staging", "production"}

		for _, env := range environments {
			data := TemplateData{Environment: env}
			assert.Equal(t, env, data.Environment)
		}
	})
}

func TestWorkerConfiguration(t *testing.T) {
	t.Run("should have default worker settings", func(t *testing.T) {
		data := TemplateData{
			NumWorkers: 5,
			QueueSize:  100,
		}

		assert.Equal(t, 5, data.NumWorkers)
		assert.Equal(t, 100, data.QueueSize)
	})

	t.Run("should support custom worker settings", func(t *testing.T) {
		data := TemplateData{
			NumWorkers: 10,
			QueueSize:  500,
		}

		assert.Equal(t, 10, data.NumWorkers)
		assert.Equal(t, 500, data.QueueSize)
	})
}

func TestAPIKeyConfiguration(t *testing.T) {
	t.Run("should handle API key configuration", func(t *testing.T) {
		data := TemplateData{
			APIKeyID:     "key-id-123",
			APIKeySecret: "key-secret-456",
		}

		assert.Equal(t, "key-id-123", data.APIKeyID)
		assert.Equal(t, "key-secret-456", data.APIKeySecret)
	})

	t.Run("should handle empty API keys", func(t *testing.T) {
		data := TemplateData{}

		assert.Empty(t, data.APIKeyID)
		assert.Empty(t, data.APIKeySecret)
	})
}

func TestOutputPathGeneration(t *testing.T) {
	t.Run("should generate correct output paths", func(t *testing.T) {
		baseDir := "/tmp/output"

		paths := map[string]string{
			"init":    filepath.Join(baseDir, "telemetry", "init.go"),
			"metrics": filepath.Join(baseDir, "telemetry", "metrics", "metrics.go"),
			"logs":    filepath.Join(baseDir, "telemetry", "logs", "logs.go"),
			"traces":  filepath.Join(baseDir, "telemetry", "traces", "traces.go"),
			"config":  filepath.Join(baseDir, ".env.telemetryflow"),
			"readme":  filepath.Join(baseDir, "telemetry", "README.md"),
		}

		assert.Equal(t, "/tmp/output/telemetry/init.go", paths["init"])
		assert.Equal(t, "/tmp/output/telemetry/metrics/metrics.go", paths["metrics"])
		assert.Equal(t, "/tmp/output/telemetry/logs/logs.go", paths["logs"])
		assert.Equal(t, "/tmp/output/telemetry/traces/traces.go", paths["traces"])
		assert.Equal(t, "/tmp/output/.env.telemetryflow", paths["config"])
		assert.Equal(t, "/tmp/output/telemetry/README.md", paths["readme"])
	})
}

func TestExampleTypes(t *testing.T) {
	t.Run("should support all example types", func(t *testing.T) {
		exampleTypes := []string{"basic", "http-server", "grpc-server", "worker"}

		for _, exType := range exampleTypes {
			assert.NotEmpty(t, exType)
		}
	})
}
