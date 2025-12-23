// Package banner provides ASCII art banner for TelemetryFlow Go SDK.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package banner

import (
	"fmt"
	"strings"
)

// Config holds banner configuration
type Config struct {
	ProductName string
	Version     string
	Motto       string
	GitCommit   string
	BuildTime   string
	GoVersion   string
	Platform    string
	Vendor      string
	VendorURL   string
	Developer   string
	License     string
	SupportURL  string
	Copyright   string
}

// DefaultConfig returns default configuration for SDK
func DefaultConfig() Config {
	return Config{
		ProductName: "TelemetryFlow Go SDK",
		Version:     "1.0.0",
		Motto:       "Community Enterprise Observability Platform (CEOP)",
		GitCommit:   "unknown",
		BuildTime:   "unknown",
		GoVersion:   "unknown",
		Platform:    "unknown",
		Vendor:      "TelemetryFlow",
		VendorURL:   "https://telemetryflow.id",
		Developer:   "DevOpsCorner Indonesia",
		License:     "Apache-2.0",
		SupportURL:  "https://docs.telemetryflow.id",
		Copyright:   "Copyright (c) 2024-2026 DevOpsCorner Indonesia",
	}
}

// GeneratorConfig returns configuration for code generator
func GeneratorConfig() Config {
	cfg := DefaultConfig()
	cfg.ProductName = "TelemetryFlow Code Generator"
	return cfg
}

// RESTfulAPIGeneratorConfig returns configuration for RESTful API generator
func RESTfulAPIGeneratorConfig() Config {
	cfg := DefaultConfig()
	cfg.ProductName = "TelemetryFlow RESTful API Generator"
	cfg.Motto = "DDD + CQRS Pattern RESTful API Generator"
	return cfg
}

// Generate creates the full banner string with ASCII art
func Generate(cfg Config) string {
	return fmt.Sprintf(`
    ___________    .__                        __
    \__    ___/___ |  |   ____   _____   _____/  |________ ___.__.
      |    |_/ __ \|  | _/ __ \ /     \_/ __ \   __\_  __ <   |  |
      |    |\  ___/|  |_\  ___/|  Y Y  \  ___/|  |  |  | \/\___  |
      |____| \___  >____/\___  >__|_|  /\___  >__|  |__|   / ____|
                 \/          \/      \/     \/             \/
                    ___________.__
                    \_   _____/|  |   ______  _  __
                     |    __)  |  |  /  _ \ \/ \/ /
                     |     \   |  |_(  <_> )     /
                     |___  /   |____/\____/ \/\_/
                         \/
              _________  ________   ____  __.
             /   _____/ \______ \ |    |/ _|
             \_____  \   |    |  \|      <
             /        \  |    '   \    |  \
            /_______  / /_______  /____|__ \
                    \/          \/        \/

  %s
    %s v%s
    %s
  %s
    Platform     %s
    Go Version   %s
    Commit       %s
    Built        %s
  %s
    Vendor       %s (%s)
    Developer    %s
    License      %s
    Support      %s
  %s
    %s
  %s

`, strings.Repeat("=", 78),
		cfg.ProductName, cfg.Version, cfg.Motto,
		strings.Repeat("=", 78),
		cfg.Platform, cfg.GoVersion, cfg.GitCommit, cfg.BuildTime,
		strings.Repeat("-", 78),
		cfg.Vendor, cfg.VendorURL, cfg.Developer, cfg.License, cfg.SupportURL,
		strings.Repeat("-", 78),
		cfg.Copyright,
		strings.Repeat("=", 78))
}

// GenerateCompact creates a compact banner without ASCII art
func GenerateCompact(cfg Config) string {
	return fmt.Sprintf(`
  %s
    %s v%s - %s
  %s
    %s
  %s

`, strings.Repeat("=", 78),
		cfg.ProductName, cfg.Version, cfg.Motto,
		strings.Repeat("=", 78),
		cfg.Copyright,
		strings.Repeat("=", 78))
}

// GenerateMinimal creates a minimal one-line banner
func GenerateMinimal(cfg Config) string {
	return fmt.Sprintf("%s v%s - %s\n", cfg.ProductName, cfg.Version, cfg.Motto)
}

// Print prints the full banner to stdout
func Print(cfg Config) {
	fmt.Print(Generate(cfg))
}

// PrintCompact prints the compact banner to stdout
func PrintCompact(cfg Config) {
	fmt.Print(GenerateCompact(cfg))
}

// PrintMinimal prints the minimal banner to stdout
func PrintMinimal(cfg Config) {
	fmt.Print(GenerateMinimal(cfg))
}
