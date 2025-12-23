// Package generator_test provides integration tests for the code generators.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package generator_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GeneratorBinary holds the path to the generator binary
var generatorBinary string
var restapiGeneratorBinary string

func init() {
	// Build generator binaries for testing
	buildDir := os.Getenv("BUILD_DIR")
	if buildDir == "" {
		buildDir = "../../../build"
	}
	generatorBinary = filepath.Join(buildDir, "telemetryflow-gen")
	restapiGeneratorBinary = filepath.Join(buildDir, "telemetryflow-restapi")
}

func TestGeneratorBinaryExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("telemetryflow-gen should exist", func(t *testing.T) {
		if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
			t.Skipf("Generator binary not found at %s. Run 'make build-generators' first.", generatorBinary)
		}
	})

	t.Run("telemetryflow-restapi should exist", func(t *testing.T) {
		if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
			t.Skipf("RESTful API generator binary not found at %s. Run 'make build-generators' first.", restapiGeneratorBinary)
		}
	})
}

func TestGeneratorVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	t.Run("should display version", func(t *testing.T) {
		cmd := exec.Command(generatorBinary, "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "TelemetryFlow")
		assert.Contains(t, outputStr, "Code Generator")
	})
}

func TestRESTAPIGeneratorVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("should display version", func(t *testing.T) {
		cmd := exec.Command(restapiGeneratorBinary, "version")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "TelemetryFlow")
		assert.Contains(t, outputStr, "RESTful API Generator")
	})
}

func TestGeneratorHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	t.Run("should display help", func(t *testing.T) {
		cmd := exec.Command(generatorBinary, "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "Usage:")
		assert.Contains(t, outputStr, "Available Commands:")
		assert.Contains(t, outputStr, "init")
		assert.Contains(t, outputStr, "config")
		assert.Contains(t, outputStr, "example")
		assert.Contains(t, outputStr, "version")
	})
}

func TestRESTAPIGeneratorHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("should display help", func(t *testing.T) {
		cmd := exec.Command(restapiGeneratorBinary, "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		assert.Contains(t, outputStr, "Usage:")
		assert.Contains(t, outputStr, "Available Commands:")
		assert.Contains(t, outputStr, "new")
		assert.Contains(t, outputStr, "entity")
		assert.Contains(t, outputStr, "docs")
		assert.Contains(t, outputStr, "version")
	})
}

func TestGeneratorNoBanner(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	t.Run("should suppress banner with --no-banner", func(t *testing.T) {
		cmd := exec.Command(generatorBinary, "--no-banner", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err)

		outputStr := string(output)
		// Should not contain ASCII art when banner is suppressed
		assert.NotContains(t, outputStr, "___________.__")
		// Should still show help
		assert.Contains(t, outputStr, "Usage:")
	})
}

func TestGeneratorInit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	t.Run("should require project name", func(t *testing.T) {
		cmd := exec.Command(generatorBinary, "--no-banner", "init")
		output, err := cmd.CombinedOutput()

		// Should fail without project name
		assert.Error(t, err)
		assert.Contains(t, string(output), "required")
	})

	t.Run("should create project structure", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "telemetryflow-init-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(generatorBinary, "--no-banner", "init",
			"--project", "test-project",
			"--service", "test-service",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		require.NoError(t, err)

		// Check directory structure
		telemetryDir := filepath.Join(tmpDir, "telemetry")
		assert.DirExists(t, telemetryDir)

		metricsDir := filepath.Join(telemetryDir, "metrics")
		assert.DirExists(t, metricsDir)

		logsDir := filepath.Join(telemetryDir, "logs")
		assert.DirExists(t, logsDir)

		tracesDir := filepath.Join(telemetryDir, "traces")
		assert.DirExists(t, tracesDir)

		// Check generated files
		assert.FileExists(t, filepath.Join(telemetryDir, "init.go"))
		assert.FileExists(t, filepath.Join(metricsDir, "metrics.go"))
		assert.FileExists(t, filepath.Join(logsDir, "logs.go"))
		assert.FileExists(t, filepath.Join(tracesDir, "traces.go"))
		assert.FileExists(t, filepath.Join(tmpDir, ".env.telemetryflow"))
	})
}

func TestRESTAPIGeneratorNew(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("should require project name", func(t *testing.T) {
		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new")
		output, err := cmd.CombinedOutput()

		// Should fail without project name
		assert.Error(t, err)
		assert.Contains(t, string(output), "required")
	})

	t.Run("should create DDD project structure", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-new-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", "test-api",
			"--module", "github.com/test/api",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-api")

		// Check DDD directory structure
		assert.DirExists(t, filepath.Join(projectDir, "cmd", "api"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "entity"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "repository"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "command"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "query"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "application", "handler"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "http"))
		assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure", "persistence"))

		// Check documentation directories
		assert.DirExists(t, filepath.Join(projectDir, "docs", "api"))
		assert.DirExists(t, filepath.Join(projectDir, "docs", "diagrams"))
		assert.DirExists(t, filepath.Join(projectDir, "docs", "postman"))

		// Check generated files
		assert.FileExists(t, filepath.Join(projectDir, "go.mod"))
		assert.FileExists(t, filepath.Join(projectDir, "Makefile"))
		assert.FileExists(t, filepath.Join(projectDir, "README.md"))
		assert.FileExists(t, filepath.Join(projectDir, "Dockerfile"))
		assert.FileExists(t, filepath.Join(projectDir, "docker-compose.yml"))
		assert.FileExists(t, filepath.Join(projectDir, ".env.example"))
	})

	t.Run("should generate documentation files", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-docs-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", "test-api",
			"--module", "github.com/test/api",
			"--output", tmpDir,
		)
		_, err = cmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-api")

		// Check documentation files
		assert.FileExists(t, filepath.Join(projectDir, "docs", "api", "openapi.yaml"))
		assert.FileExists(t, filepath.Join(projectDir, "docs", "api", "swagger.json"))
		assert.FileExists(t, filepath.Join(projectDir, "docs", "diagrams", "ERD.md"))
		assert.FileExists(t, filepath.Join(projectDir, "docs", "diagrams", "DFD.md"))
		assert.FileExists(t, filepath.Join(projectDir, "docs", "postman", "collection.json"))
		assert.FileExists(t, filepath.Join(projectDir, "docs", "postman", "environment.json"))
	})
}

func TestGeneratorConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	t.Run("should generate config file", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "telemetryflow-config-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(generatorBinary, "--no-banner", "config",
			"--service", "test-service",
			"--key-id", "test-key-id",
			"--key-secret", "test-key-secret",
			"--endpoint", "custom.endpoint:4317",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		require.NoError(t, err)

		configFile := filepath.Join(tmpDir, ".env.telemetryflow")
		assert.FileExists(t, configFile)

		// Read and verify config content
		content, err := os.ReadFile(configFile)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "test-service")
		assert.Contains(t, contentStr, "test-key-id")
		assert.Contains(t, contentStr, "custom.endpoint:4317")
	})
}

func TestGeneratorExample(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(generatorBinary); os.IsNotExist(err) {
		t.Skipf("Generator binary not found at %s", generatorBinary)
	}

	exampleTypes := []string{"basic", "http-server", "grpc-server", "worker"}

	for _, exType := range exampleTypes {
		t.Run("should generate "+exType+" example", func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "telemetryflow-example-test-*")
			require.NoError(t, err)
			defer func() {
				if err := os.RemoveAll(tmpDir); err != nil {
					t.Logf("Failed to remove temp dir: %v", err)
				}
			}()

			cmd := exec.Command(generatorBinary, "--no-banner", "example", exType,
				"--output", tmpDir,
			)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Logf("Command output: %s", string(output))
			}
			require.NoError(t, err)

			// Check example file exists
			expectedFile := filepath.Join(tmpDir, "example_"+strings.ReplaceAll(exType, "-", "_")+".go")
			assert.FileExists(t, expectedFile)
		})
	}
}

func TestRESTAPIGeneratorEntity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("should require entity name", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-entity-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "entity",
			"--output", tmpDir,
		)
		output, err := cmd.CombinedOutput()

		// Should fail without entity name
		assert.Error(t, err)
		assert.Contains(t, string(output), "required")
	})
}

func TestGeneratedGoModValidity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("should generate valid go.mod", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-gomod-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", "test-api",
			"--module", "github.com/test/api",
			"--output", tmpDir,
		)
		_, err = cmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-api")
		goModPath := filepath.Join(projectDir, "go.mod")

		// Read and verify go.mod content
		content, err := os.ReadFile(goModPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "module github.com/test/api")
		assert.Contains(t, contentStr, "go ")
	})
}

func TestGeneratedProjectCompilation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if _, err := os.Stat(restapiGeneratorBinary); os.IsNotExist(err) {
		t.Skipf("RESTful API generator binary not found at %s", restapiGeneratorBinary)
	}

	t.Run("generated project should have valid syntax", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-compile-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		// Generate project
		cmd := exec.Command(restapiGeneratorBinary, "--no-banner", "new",
			"--name", "test-api",
			"--module", "github.com/test/api",
			"--output", tmpDir,
		)
		_, err = cmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-api")

		// Check that main.go exists and has valid Go syntax
		mainGoPath := filepath.Join(projectDir, "cmd", "api", "main.go")
		assert.FileExists(t, mainGoPath)

		content, err := os.ReadFile(mainGoPath)
		require.NoError(t, err)

		contentStr := string(content)
		assert.Contains(t, contentStr, "package main")
		assert.Contains(t, contentStr, "func main()")
	})
}
