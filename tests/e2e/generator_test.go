// Package e2e_test provides end-to-end tests for TelemetryFlow Go SDK generators.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test binary paths
var (
	generatorBinary        string
	restapiGeneratorBinary string
)

func init() {
	buildDir := os.Getenv("BUILD_DIR")
	if buildDir == "" {
		buildDir = "../../build"
	}
	generatorBinary = filepath.Join(buildDir, "telemetryflow-gen")
	restapiGeneratorBinary = filepath.Join(buildDir, "telemetryflow-restapi")
}

func skipIfBinaryNotFound(t *testing.T, binary string) {
	if _, err := os.Stat(binary); os.IsNotExist(err) {
		t.Skipf("Binary not found at %s. Run 'make build-generators' first.", binary)
	}
}

func TestE2EGeneratorFullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	skipIfBinaryNotFound(t, generatorBinary)

	t.Run("complete SDK integration workflow", func(t *testing.T) {
		// Create temporary project directory
		tmpDir, err := os.MkdirTemp("", "e2e-telemetryflow-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		// Step 1: Initialize TelemetryFlow integration
		cmd := exec.Command(generatorBinary, "--no-banner", "init",
			"--project", "e2e-test-project",
			"--service", "e2e-test-service",
			"--key-id", "test-key-id",
			"--key-secret", "test-key-secret",
			"--endpoint", "localhost:4317",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("Init command output: %s", string(output))
		}
		require.NoError(t, err, "Init command should succeed")

		// Verify telemetry directory structure
		telemetryDir := filepath.Join(tmpDir, "telemetry")
		assert.DirExists(t, telemetryDir)
		assert.DirExists(t, filepath.Join(telemetryDir, "metrics"))
		assert.DirExists(t, filepath.Join(telemetryDir, "logs"))
		assert.DirExists(t, filepath.Join(telemetryDir, "traces"))

		// Verify generated files
		assert.FileExists(t, filepath.Join(telemetryDir, "init.go"))
		assert.FileExists(t, filepath.Join(telemetryDir, "metrics", "metrics.go"))
		assert.FileExists(t, filepath.Join(telemetryDir, "logs", "logs.go"))
		assert.FileExists(t, filepath.Join(telemetryDir, "traces", "traces.go"))
		assert.FileExists(t, filepath.Join(telemetryDir, "README.md"))
		assert.FileExists(t, filepath.Join(tmpDir, ".env.telemetryflow"))

		// Verify config content
		configContent, err := os.ReadFile(filepath.Join(tmpDir, ".env.telemetryflow"))
		require.NoError(t, err)
		assert.Contains(t, string(configContent), "e2e-test-service")
		assert.Contains(t, string(configContent), "test-key-id")
		assert.Contains(t, string(configContent), "localhost:4317")

		// Step 2: Generate examples
		for _, exampleType := range []string{"basic", "http-server"} {
			cmd = exec.Command(generatorBinary, "--no-banner", "example", exampleType,
				"--output", tmpDir,
			)
			output, err = cmd.CombinedOutput()
			if err != nil {
				t.Logf("Example %s command output: %s", exampleType, string(output))
			}
			require.NoError(t, err, "Example %s command should succeed", exampleType)
		}

		// Verify example files
		assert.FileExists(t, filepath.Join(tmpDir, "example_basic.go"))
		assert.FileExists(t, filepath.Join(tmpDir, "example_http_server.go"))
	})
}

func TestE2ERESTAPIGeneratorFullWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	skipIfBinaryNotFound(t, restapiGeneratorBinary)

	t.Run("complete RESTful API project workflow", func(t *testing.T) {
		// Create temporary project directory
		tmpDir, err := os.MkdirTemp("", "e2e-restapi-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		projectName := "e2e-api"
		modulePath := "github.com/e2e/api"

		// Step 1: Create new RESTful API project
		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", projectName,
			"--module", modulePath,
			"--service", "e2e-api-service",
			"--env", "development",
			"--db-driver", "postgres",
			"--db-host", "localhost",
			"--db-port", "5432",
			"--db-name", "e2e_api",
			"--db-user", "postgres",
			"--port", "8080",
			"--telemetry",
			"--swagger",
			"--cors",
			"--auth",
			"--rate-limit",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("New command output: %s", string(output))
		}
		require.NoError(t, err, "New command should succeed")

		projectDir := filepath.Join(tmpDir, projectName)

		// Verify DDD directory structure
		t.Run("verify DDD structure", func(t *testing.T) {
			// Domain layer
			assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "entity"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "repository"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "valueobject"))

			// Application layer (CQRS)
			assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "command"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "query"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "handler"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "dto"))

			// Infrastructure layer
			assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "middleware"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "handler"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "persistence"))
			assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "config"))
		})

		// Verify project files
		t.Run("verify project files", func(t *testing.T) {
			assert.FileExists(t, filepath.Join(projectDir, "go.mod"))
			assert.FileExists(t, filepath.Join(projectDir, "cmd", "api", "main.go"))
			assert.FileExists(t, filepath.Join(projectDir, "Makefile"))
			assert.FileExists(t, filepath.Join(projectDir, "Dockerfile"))
			assert.FileExists(t, filepath.Join(projectDir, "docker-compose.yml"))
			assert.FileExists(t, filepath.Join(projectDir, ".env.example"))
			assert.FileExists(t, filepath.Join(projectDir, ".gitignore"))
			assert.FileExists(t, filepath.Join(projectDir, "README.md"))
		})

		// Verify documentation files
		t.Run("verify documentation", func(t *testing.T) {
			// API documentation
			assert.FileExists(t, filepath.Join(projectDir, "docs", "api", "openapi.yaml"))
			assert.FileExists(t, filepath.Join(projectDir, "docs", "api", "swagger.json"))

			// Diagrams
			assert.FileExists(t, filepath.Join(projectDir, "docs", "diagrams", "ERD.md"))
			assert.FileExists(t, filepath.Join(projectDir, "docs", "diagrams", "DFD.md"))

			// Postman
			assert.FileExists(t, filepath.Join(projectDir, "docs", "postman", "collection.json"))
			assert.FileExists(t, filepath.Join(projectDir, "docs", "postman", "environment.json"))
		})

		// Verify infrastructure files
		t.Run("verify infrastructure", func(t *testing.T) {
			// HTTP server
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "server.go"))
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "router.go"))

			// Middleware
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "middleware", "logger.go"))
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "middleware", "auth.go"))
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "middleware", "cors.go"))
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http", "middleware", "ratelimit.go"))

			// Database
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "persistence", "database.go"))

			// Config
			assert.FileExists(t, filepath.Join(projectDir, "internal", "infrastructure", "config", "config.go"))
			assert.FileExists(t, filepath.Join(projectDir, "configs", "config.yaml"))
		})

		// Verify telemetry integration
		t.Run("verify telemetry", func(t *testing.T) {
			assert.DirExists(t, filepath.Join(projectDir, "telemetry"))
			assert.FileExists(t, filepath.Join(projectDir, "telemetry", "init.go"))
			assert.FileExists(t, filepath.Join(projectDir, "telemetry", "metrics", "metrics.go"))
			assert.FileExists(t, filepath.Join(projectDir, "telemetry", "logs", "logs.go"))
			assert.FileExists(t, filepath.Join(projectDir, "telemetry", "traces", "traces.go"))
		})

		// Verify migrations
		t.Run("verify migrations", func(t *testing.T) {
			assert.DirExists(t, filepath.Join(projectDir, "migrations"))
			assert.FileExists(t, filepath.Join(projectDir, "migrations", "000001_init.up.sql"))
			assert.FileExists(t, filepath.Join(projectDir, "migrations", "000001_init.down.sql"))
		})

		// Verify scripts
		t.Run("verify scripts", func(t *testing.T) {
			assert.DirExists(t, filepath.Join(projectDir, "scripts"))
			assert.FileExists(t, filepath.Join(projectDir, "scripts", "run.sh"))
			assert.FileExists(t, filepath.Join(projectDir, "scripts", "test.sh"))
		})

		// Verify PKG utilities
		t.Run("verify pkg utilities", func(t *testing.T) {
			assert.FileExists(t, filepath.Join(projectDir, "pkg", "logger", "logger.go"))
			assert.FileExists(t, filepath.Join(projectDir, "pkg", "validator", "validator.go"))
			assert.FileExists(t, filepath.Join(projectDir, "pkg", "response", "response.go"))
		})

		// Verify go.mod content
		t.Run("verify go.mod content", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "module "+modulePath)
		})

		// Verify main.go content
		t.Run("verify main.go content", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "cmd", "api", "main.go"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "package main")
			assert.Contains(t, string(content), "func main()")
		})

		// Verify OpenAPI spec content
		t.Run("verify OpenAPI content", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "docs", "api", "openapi.yaml"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "openapi:")
			assert.Contains(t, string(content), "info:")
			assert.Contains(t, string(content), projectName)
		})

		// Verify Postman collection content
		t.Run("verify Postman collection", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "docs", "postman", "collection.json"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "\"name\":")
			assert.Contains(t, string(content), "\"item\":")
		})

		// Verify ERD content (Mermaid)
		t.Run("verify ERD content", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "docs", "diagrams", "ERD.md"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "```mermaid")
			assert.Contains(t, string(content), "erDiagram")
		})

		// Verify DFD content (Mermaid)
		t.Run("verify DFD content", func(t *testing.T) {
			content, err := os.ReadFile(filepath.Join(projectDir, "docs", "diagrams", "DFD.md"))
			require.NoError(t, err)
			assert.Contains(t, string(content), "```mermaid")
		})
	})
}

func TestE2EGeneratorWithDisabledFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	skipIfBinaryNotFound(t, restapiGeneratorBinary)

	t.Run("create project with disabled features", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "e2e-disabled-features-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", "minimal-api",
			"--module", "github.com/test/minimal",
			"--telemetry=false",
			"--swagger=false",
			"--cors=false",
			"--auth=false",
			"--rate-limit=false",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Logf("New command output: %s", string(output))
		}
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "minimal-api")

		// Core structure should still exist
		assert.DirExists(t, filepath.Join(projectDir, "cmd", "api"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "domain"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "application"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure"))

		// Basic files should exist
		assert.FileExists(t, filepath.Join(projectDir, "go.mod"))
		assert.FileExists(t, filepath.Join(projectDir, "cmd", "api", "main.go"))
	})
}

func TestE2EMultipleProjects(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	skipIfBinaryNotFound(t, restapiGeneratorBinary)

	t.Run("create multiple independent projects", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "e2e-multi-project-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)

		projects := []struct {
			name   string
			module string
		}{
			{"project-a", "github.com/test/project-a"},
			{"project-b", "github.com/test/project-b"},
			{"project-c", "github.com/test/project-c"},
		}

		for _, p := range projects {
			cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
				"--name", p.name,
				"--module", p.module,
				"--output", tmpDir,
			)
			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Logf("Creating %s output: %s", p.name, string(output))
			}
			require.NoError(t, err, "Should create %s", p.name)

			// Verify each project is independent
			projectDir := filepath.Join(tmpDir, p.name)
			assert.DirExists(t, projectDir)
			assert.FileExists(t, filepath.Join(projectDir, "go.mod"))

			content, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
			require.NoError(t, err)
			assert.Contains(t, string(content), p.module)
		}
	})
}

func TestE2EBannerDisplay(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	skipIfBinaryNotFound(t, generatorBinary)
	skipIfBinaryNotFound(t, restapiGeneratorBinary)

	t.Run("telemetryflow-gen displays banner", func(t *testing.T) {
		cmd := exec.Command(generatorBinary, "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "TelemetryFlow")
		assert.Contains(t, outputStr, "Code Generator")
		assert.Contains(t, outputStr, "DevOpsCorner Indonesia")
		assert.Contains(t, outputStr, "Apache-2.0")
	})

	t.Run("telemetryflow-restapi displays banner", func(t *testing.T) {
		cmd := exec.Command(restapiGeneratorBinary, "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "TelemetryFlow")
		assert.Contains(t, outputStr, "RESTful API Generator")
		assert.Contains(t, outputStr, "DDD + CQRS")
		assert.Contains(t, outputStr, "DevOpsCorner Indonesia")
	})
}
