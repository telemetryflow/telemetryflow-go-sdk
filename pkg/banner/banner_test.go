// Package banner provides unit tests for the banner package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package banner

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("should return default configuration", func(t *testing.T) {
		cfg := DefaultConfig()

		assert.Equal(t, "TelemetryFlow Go SDK", cfg.ProductName)
		assert.Equal(t, "1.1.0", cfg.Version)
		assert.Equal(t, "Community Enterprise Observability Platform (CEOP)", cfg.Motto)
		assert.Equal(t, "TelemetryFlow", cfg.Vendor)
		assert.Equal(t, "https://telemetryflow.id", cfg.VendorURL)
		assert.Equal(t, "DevOpsCorner Indonesia", cfg.Developer)
		assert.Equal(t, "Apache-2.0", cfg.License)
		assert.Equal(t, "https://docs.telemetryflow.id", cfg.SupportURL)
		assert.Contains(t, cfg.Copyright, "DevOpsCorner Indonesia")
	})

	t.Run("should have unknown as default for build info", func(t *testing.T) {
		cfg := DefaultConfig()

		assert.Equal(t, "unknown", cfg.GitCommit)
		assert.Equal(t, "unknown", cfg.BuildTime)
		assert.Equal(t, "unknown", cfg.GoVersion)
		assert.Equal(t, "unknown", cfg.Platform)
	})
}

func TestGeneratorConfig(t *testing.T) {
	t.Run("should return generator configuration", func(t *testing.T) {
		cfg := GeneratorConfig()

		assert.Equal(t, "TelemetryFlow Code Generator", cfg.ProductName)
		assert.Equal(t, "1.1.0", cfg.Version)
		assert.Equal(t, "TelemetryFlow", cfg.Vendor)
	})
}

func TestRESTfulAPIGeneratorConfig(t *testing.T) {
	t.Run("should return RESTful API generator configuration", func(t *testing.T) {
		cfg := RESTfulAPIGeneratorConfig()

		assert.Equal(t, "TelemetryFlow RESTful API Generator", cfg.ProductName)
		assert.Equal(t, "DDD + CQRS Pattern RESTful API Generator", cfg.Motto)
		assert.Equal(t, "1.1.0", cfg.Version)
	})
}

func TestGenerate(t *testing.T) {
	t.Run("should generate full banner with ASCII art", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := Generate(cfg)

		assert.Contains(t, banner, "TelemetryFlow Go SDK")
		assert.Contains(t, banner, "1.1.0")
		assert.Contains(t, banner, "Community Enterprise Observability Platform")
		assert.Contains(t, banner, "DevOpsCorner Indonesia")
		assert.Contains(t, banner, "TelemetryFlow")
		// Check for ASCII art presence
		assert.Contains(t, banner, "___________")
		assert.Contains(t, banner, "\\_   _____/")
	})

	t.Run("should include all config values", func(t *testing.T) {
		cfg := Config{
			ProductName: "Test Product",
			Version:     "2.0.0",
			Motto:       "Test Motto",
			GitCommit:   "abc123",
			BuildTime:   "2024-01-01",
			GoVersion:   "1.24",
			Platform:    "linux/amd64",
			Vendor:      "Test Vendor",
			VendorURL:   "https://test.com",
			Developer:   "Test Dev",
			License:     "MIT",
			SupportURL:  "https://support.test.com",
			Copyright:   "Copyright Test",
		}

		banner := Generate(cfg)

		assert.Contains(t, banner, "Test Product")
		assert.Contains(t, banner, "2.0.0")
		assert.Contains(t, banner, "Test Motto")
		assert.Contains(t, banner, "abc123")
		assert.Contains(t, banner, "2024-01-01")
		assert.Contains(t, banner, "1.24")
		assert.Contains(t, banner, "linux/amd64")
		assert.Contains(t, banner, "Test Vendor")
		assert.Contains(t, banner, "https://test.com")
		assert.Contains(t, banner, "Test Dev")
		assert.Contains(t, banner, "MIT")
		assert.Contains(t, banner, "https://support.test.com")
		assert.Contains(t, banner, "Copyright Test")
	})
}

func TestGenerateCompact(t *testing.T) {
	t.Run("should generate compact banner", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := GenerateCompact(cfg)

		assert.Contains(t, banner, "TelemetryFlow Go SDK")
		assert.Contains(t, banner, "1.1.0")
		assert.Contains(t, banner, "DevOpsCorner Indonesia")
		// Should not contain ASCII art
		assert.NotContains(t, banner, "\\_   _____/")
	})

	t.Run("should be shorter than full banner", func(t *testing.T) {
		cfg := DefaultConfig()

		fullBanner := Generate(cfg)
		compactBanner := GenerateCompact(cfg)

		assert.Less(t, len(compactBanner), len(fullBanner))
	})
}

func TestGenerateMinimal(t *testing.T) {
	t.Run("should generate minimal one-line banner", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := GenerateMinimal(cfg)

		assert.Contains(t, banner, "TelemetryFlow Go SDK")
		assert.Contains(t, banner, "v1.1.0")
		assert.Contains(t, banner, "Community Enterprise Observability Platform")
	})

	t.Run("should be single line", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := GenerateMinimal(cfg)

		lines := strings.Split(strings.TrimSpace(banner), "\n")
		assert.Equal(t, 1, len(lines))
	})

	t.Run("should be shorter than compact banner", func(t *testing.T) {
		cfg := DefaultConfig()

		compactBanner := GenerateCompact(cfg)
		minimalBanner := GenerateMinimal(cfg)

		assert.Less(t, len(minimalBanner), len(compactBanner))
	})
}

func TestConfig_CustomValues(t *testing.T) {
	t.Run("should use custom configuration values", func(t *testing.T) {
		cfg := Config{
			ProductName: "Custom SDK",
			Version:     "3.0.0",
			Motto:       "Custom Motto",
		}

		banner := GenerateMinimal(cfg)

		assert.Contains(t, banner, "Custom SDK")
		assert.Contains(t, banner, "v3.0.0")
		assert.Contains(t, banner, "Custom Motto")
	})
}

func TestBannerSeparators(t *testing.T) {
	t.Run("should contain separators in full banner", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := Generate(cfg)

		// Check for separator lines (78 = characters)
		assert.Contains(t, banner, strings.Repeat("=", 78))
		assert.Contains(t, banner, strings.Repeat("-", 78))
	})

	t.Run("should contain separators in compact banner", func(t *testing.T) {
		cfg := DefaultConfig()

		banner := GenerateCompact(cfg)

		assert.Contains(t, banner, strings.Repeat("=", 78))
	})
}
