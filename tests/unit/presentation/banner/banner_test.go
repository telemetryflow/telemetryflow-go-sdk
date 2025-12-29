// Package banner_test provides unit tests for the banner package.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
package banner_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/banner"
)

func TestDefaultConfig(t *testing.T) {
	t.Run("should return valid default config", func(t *testing.T) {
		cfg := banner.DefaultConfig()

		assert.Equal(t, "TelemetryFlow Go SDK", cfg.ProductName)
		assert.Equal(t, "1.1.1", cfg.Version)
		assert.Equal(t, "Community Enterprise Observability Platform (CEOP)", cfg.Motto)
		assert.Equal(t, "TelemetryFlow", cfg.Vendor)
		assert.Equal(t, "https://telemetryflow.id", cfg.VendorURL)
		assert.Equal(t, "DevOpsCorner Indonesia", cfg.Developer)
		assert.Equal(t, "Apache-2.0", cfg.License)
		assert.Equal(t, "https://docs.telemetryflow.id", cfg.SupportURL)
	})

	t.Run("should have copyright info", func(t *testing.T) {
		cfg := banner.DefaultConfig()

		assert.Contains(t, cfg.Copyright, "DevOpsCorner Indonesia")
		assert.Contains(t, cfg.Copyright, "2024-2026")
	})
}

func TestGeneratorConfig(t *testing.T) {
	t.Run("should return generator-specific config", func(t *testing.T) {
		cfg := banner.GeneratorConfig()

		assert.Equal(t, "TelemetryFlow Code Generator", cfg.ProductName)
		assert.Equal(t, "Community Enterprise Observability Platform (CEOP)", cfg.Motto)
	})

	t.Run("should inherit default values", func(t *testing.T) {
		cfg := banner.GeneratorConfig()

		assert.Equal(t, "TelemetryFlow", cfg.Vendor)
		assert.Equal(t, "DevOpsCorner Indonesia", cfg.Developer)
		assert.Equal(t, "Apache-2.0", cfg.License)
	})
}

func TestRESTfulAPIGeneratorConfig(t *testing.T) {
	t.Run("should return RESTful API generator config", func(t *testing.T) {
		cfg := banner.RESTfulAPIGeneratorConfig()

		assert.Equal(t, "TelemetryFlow RESTful API Generator", cfg.ProductName)
		assert.Equal(t, "DDD + CQRS Pattern RESTful API Generator", cfg.Motto)
	})

	t.Run("should inherit default values", func(t *testing.T) {
		cfg := banner.RESTfulAPIGeneratorConfig()

		assert.Equal(t, "TelemetryFlow", cfg.Vendor)
		assert.Equal(t, "DevOpsCorner Indonesia", cfg.Developer)
		assert.Equal(t, "Apache-2.0", cfg.License)
	})
}

func TestGenerate(t *testing.T) {
	cfg := banner.DefaultConfig()

	t.Run("should return ASCII art banner", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.NotEmpty(t, b)
		// Banner should contain ASCII art characters
		assert.Contains(t, b, "___")
		assert.Contains(t, b, "/")
		assert.Contains(t, b, "\\")
	})

	t.Run("should contain product name", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.ProductName)
	})

	t.Run("should contain version", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.Version)
	})

	t.Run("should contain motto", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.Motto)
	})

	t.Run("should contain vendor info", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.Vendor)
		assert.Contains(t, b, cfg.VendorURL)
	})

	t.Run("should contain developer info", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.Developer)
	})

	t.Run("should contain license info", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.License)
	})

	t.Run("should contain support URL", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.SupportURL)
	})

	t.Run("should contain copyright", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, cfg.Copyright)
	})

	t.Run("should be multi-line", func(t *testing.T) {
		b := banner.Generate(cfg)
		lines := strings.Split(b, "\n")

		// Full banner should have many lines (ASCII art)
		assert.Greater(t, len(lines), 20)
	})

	t.Run("should contain separator lines", func(t *testing.T) {
		b := banner.Generate(cfg)

		assert.Contains(t, b, strings.Repeat("=", 78))
		assert.Contains(t, b, strings.Repeat("-", 78))
	})
}

func TestGenerateCompact(t *testing.T) {
	cfg := banner.DefaultConfig()

	t.Run("should return compact banner", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		assert.NotEmpty(t, b)
	})

	t.Run("should contain product name and version", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		assert.Contains(t, b, cfg.ProductName)
		assert.Contains(t, b, cfg.Version)
	})

	t.Run("should contain motto", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		assert.Contains(t, b, cfg.Motto)
	})

	t.Run("should contain copyright", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		assert.Contains(t, b, cfg.Copyright)
	})

	t.Run("should be shorter than full banner", func(t *testing.T) {
		full := banner.Generate(cfg)
		compact := banner.GenerateCompact(cfg)

		assert.Less(t, len(compact), len(full))
	})

	t.Run("should not contain ASCII art", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		// Compact banner should not have the complex ASCII art patterns
		assert.NotContains(t, b, "___________.__")
		assert.NotContains(t, b, "\\_   _____/")
	})
}

func TestGenerateMinimal(t *testing.T) {
	cfg := banner.DefaultConfig()

	t.Run("should return minimal banner", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)

		assert.NotEmpty(t, b)
	})

	t.Run("should be single line", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)
		// Trim the newline and check
		trimmed := strings.TrimSuffix(b, "\n")

		assert.NotContains(t, trimmed, "\n")
	})

	t.Run("should contain product name", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)

		assert.Contains(t, b, cfg.ProductName)
	})

	t.Run("should contain version", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)

		assert.Contains(t, b, cfg.Version)
	})

	t.Run("should contain motto", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)

		assert.Contains(t, b, cfg.Motto)
	})

	t.Run("should be shorter than compact banner", func(t *testing.T) {
		compact := banner.GenerateCompact(cfg)
		minimal := banner.GenerateMinimal(cfg)

		assert.Less(t, len(minimal), len(compact))
	})
}

func TestConfigCustomization(t *testing.T) {
	t.Run("should use custom values", func(t *testing.T) {
		cfg := banner.Config{
			ProductName: "Custom Product",
			Version:     "2.0.0",
			Motto:       "Custom Motto",
			GitCommit:   "abc123",
			BuildTime:   "2024-01-01T00:00:00Z",
			GoVersion:   "go1.21.0",
			Platform:    "linux/amd64",
			Vendor:      "Custom Vendor",
			VendorURL:   "https://custom.vendor",
			Developer:   "Custom Developer",
			License:     "MIT",
			SupportURL:  "https://custom.support",
			Copyright:   "Custom Copyright",
		}

		b := banner.Generate(cfg)

		assert.Contains(t, b, "Custom Product")
		assert.Contains(t, b, "2.0.0")
		assert.Contains(t, b, "Custom Motto")
		assert.Contains(t, b, "abc123")
		assert.Contains(t, b, "2024-01-01T00:00:00Z")
		assert.Contains(t, b, "go1.21.0")
		assert.Contains(t, b, "linux/amd64")
		assert.Contains(t, b, "Custom Vendor")
		assert.Contains(t, b, "https://custom.vendor")
		assert.Contains(t, b, "Custom Developer")
		assert.Contains(t, b, "MIT")
		assert.Contains(t, b, "https://custom.support")
		assert.Contains(t, b, "Custom Copyright")
	})
}

func TestBannerFormat(t *testing.T) {
	cfg := banner.DefaultConfig()

	t.Run("Generate should have correct format", func(t *testing.T) {
		b := banner.Generate(cfg)

		// Should start with newline
		assert.True(t, strings.HasPrefix(b, "\n"))
		// Should end with newline
		assert.True(t, strings.HasSuffix(b, "\n"))
	})

	t.Run("GenerateCompact should have correct format", func(t *testing.T) {
		b := banner.GenerateCompact(cfg)

		// Should start with newline
		assert.True(t, strings.HasPrefix(b, "\n"))
		// Should end with newline
		assert.True(t, strings.HasSuffix(b, "\n"))
	})

	t.Run("GenerateMinimal should end with newline", func(t *testing.T) {
		b := banner.GenerateMinimal(cfg)

		assert.True(t, strings.HasSuffix(b, "\n"))
	})
}
