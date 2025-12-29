// Package infrastructure_test provides unit tests for infrastructure components.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package infrastructure_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test configuration file handling

func TestConfigurationFileHandling(t *testing.T) {
	t.Run("should read YAML config file", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		content := `app:
  name: test-app
  port: 8080
database:
  host: localhost
  port: 5432
`
		err := os.WriteFile(configPath, []byte(content), 0644)
		require.NoError(t, err)

		readContent, err := os.ReadFile(configPath)
		require.NoError(t, err)

		assert.Contains(t, string(readContent), "app:")
		assert.Contains(t, string(readContent), "name: test-app")
		assert.Contains(t, string(readContent), "database:")
	})

	t.Run("should handle missing config file gracefully", func(t *testing.T) {
		_, err := os.ReadFile("/nonexistent/config.yaml")
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("should validate config file permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		err := os.WriteFile(configPath, []byte("test: value"), 0600)
		require.NoError(t, err)

		info, err := os.Stat(configPath)
		require.NoError(t, err)

		// Check that the file is not world-readable (security)
		perm := info.Mode().Perm()
		assert.False(t, perm&0004 != 0, "config should not be world-readable")
	})
}

func TestEnvironmentVariableHandling(t *testing.T) {
	t.Run("should read environment variable", func(t *testing.T) {
		t.Setenv("TEST_APP_NAME", "test-value")

		value := os.Getenv("TEST_APP_NAME")
		assert.Equal(t, "test-value", value)
	})

	t.Run("should handle missing environment variable", func(t *testing.T) {
		value := os.Getenv("NONEXISTENT_VAR_12345")
		assert.Empty(t, value)
	})

	t.Run("should use default when env var missing", func(t *testing.T) {
		value := os.Getenv("NONEXISTENT_VAR_12345")
		if value == "" {
			value = "default-value"
		}
		assert.Equal(t, "default-value", value)
	})

	t.Run("should override default with env var", func(t *testing.T) {
		t.Setenv("TEST_OVERRIDE_VAR", "env-value")

		value := os.Getenv("TEST_OVERRIDE_VAR")
		if value == "" {
			value = "default-value"
		}
		assert.Equal(t, "env-value", value)
	})
}

func TestEnvFileHandling(t *testing.T) {
	t.Run("should create .env.example template", func(t *testing.T) {
		tmpDir := t.TempDir()
		envPath := filepath.Join(tmpDir, ".env.example")

		content := `# Application Settings
APP_NAME=my-service
APP_ENV=development
APP_PORT=8080

# Database Settings
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USER=postgres
DB_PASSWORD=

# TelemetryFlow Settings
TELEMETRYFLOW_API_KEY_ID=
TELEMETRYFLOW_API_KEY_SECRET=
TELEMETRYFLOW_ENDPOINT=api.telemetryflow.id:4317
`
		err := os.WriteFile(envPath, []byte(content), 0644)
		require.NoError(t, err)

		readContent, err := os.ReadFile(envPath)
		require.NoError(t, err)

		assert.Contains(t, string(readContent), "APP_NAME=my-service")
		assert.Contains(t, string(readContent), "TELEMETRYFLOW_ENDPOINT=")
	})

	t.Run("should parse .env file format", func(t *testing.T) {
		content := `KEY1=value1
KEY2=value2
# Comment line
KEY3=value with spaces

KEY4=`
		lines := parseEnvContent(content)

		assert.Equal(t, "value1", lines["KEY1"])
		assert.Equal(t, "value2", lines["KEY2"])
		assert.Equal(t, "value with spaces", lines["KEY3"])
		assert.Equal(t, "", lines["KEY4"])
		assert.NotContains(t, lines, "# Comment line")
	})

	t.Run("should handle quoted values", func(t *testing.T) {
		content := `KEY1="quoted value"
KEY2='single quoted'
KEY3=unquoted`
		lines := parseEnvContent(content)

		assert.Equal(t, "quoted value", lines["KEY1"])
		assert.Equal(t, "single quoted", lines["KEY2"])
		assert.Equal(t, "unquoted", lines["KEY3"])
	})
}

// parseEnvContent simulates parsing .env file content
func parseEnvContent(content string) map[string]string {
	result := make(map[string]string)
	lines := splitLines(content)

	for _, line := range lines {
		// Skip empty lines and comments
		if line == "" || line[0] == '#' {
			continue
		}

		// Find the first = sign
		idx := indexOf(line, '=')
		if idx == -1 {
			continue
		}

		key := line[:idx]
		value := line[idx+1:]

		// Remove quotes if present
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		result[key] = value
	}

	return result
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func indexOf(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}
