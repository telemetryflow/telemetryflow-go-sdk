package main

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/telemetryflow/telemetryflow-go-sdk/internal/version"
	"github.com/telemetryflow/telemetryflow-go-sdk/pkg/banner"
)

//go:embed templates/*.tpl
var templateFS embed.FS

var (
	projectName   string
	apiKeyID      string
	apiKeySecret  string
	endpoint      string
	serviceName   string
	enableMetrics bool
	enableLogs    bool
	enableTraces  bool
	outputDir     string
	templateDir   string
	noBanner      bool
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "telemetryflow-gen",
		Short: "TelemetryFlow SDK code generator",
		Long:  "Generate boilerplate code for integrating TelemetryFlow into your Go application",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !noBanner && cmd.Name() != "version" {
				cfg := banner.GeneratorConfig()
				cfg.Version = version.Version
				cfg.GitCommit = version.GitCommit
				cfg.BuildTime = version.BuildTime
				cfg.GoVersion = version.GoVersion()
				cfg.Platform = version.Platform()
				banner.PrintCompact(cfg)
			}
		},
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize a new TelemetryFlow integration",
		Long:  "Generate all necessary files to integrate TelemetryFlow into your project",
		Run:   runInit,
	}

	var exampleCmd = &cobra.Command{
		Use:   "example [type]",
		Short: "Generate example code",
		Long:  "Generate example code for specific use cases (basic, http-server, grpc-server, worker)",
		Args:  cobra.ExactArgs(1),
		Run:   runExample,
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Generate configuration file",
		Long:  "Generate a .env configuration file with TelemetryFlow settings",
		Run:   runConfig,
	}

	// Init command flags
	initCmd.Flags().StringVarP(&projectName, "project", "p", "", "Project name (required)")
	initCmd.Flags().StringVarP(&apiKeyID, "key-id", "k", "", "TelemetryFlow API Key ID")
	initCmd.Flags().StringVarP(&apiKeySecret, "key-secret", "s", "", "TelemetryFlow API Key Secret")
	initCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "api.telemetryflow.id:4317", "OTLP endpoint")
	initCmd.Flags().StringVarP(&serviceName, "service", "n", "", "Service name (defaults to project name)")
	initCmd.Flags().BoolVar(&enableMetrics, "metrics", true, "Enable metrics")
	initCmd.Flags().BoolVar(&enableLogs, "logs", true, "Enable logs")
	initCmd.Flags().BoolVar(&enableTraces, "traces", true, "Enable traces")
	initCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	_ = initCmd.MarkFlagRequired("project")

	// Config command flags
	configCmd.Flags().StringVarP(&apiKeyID, "key-id", "k", "", "TelemetryFlow API Key ID")
	configCmd.Flags().StringVarP(&apiKeySecret, "key-secret", "s", "", "TelemetryFlow API Key Secret")
	configCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "api.telemetryflow.id:4317", "OTLP endpoint")
	configCmd.Flags().StringVarP(&serviceName, "service", "n", "", "Service name")
	configCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")

	// Example command flags
	exampleCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := banner.GeneratorConfig()
			cfg.Version = version.Version
			cfg.GitCommit = version.GitCommit
			cfg.BuildTime = version.BuildTime
			cfg.GoVersion = version.GoVersion()
			cfg.Platform = version.Platform()
			banner.Print(cfg)
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&templateDir, "template-dir", "", "Custom template directory (uses embedded templates if not set)")
	rootCmd.PersistentFlags().BoolVar(&noBanner, "no-banner", false, "Disable banner output")

	rootCmd.AddCommand(initCmd, exampleCmd, configCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runInit(cmd *cobra.Command, args []string) {
	if serviceName == "" {
		serviceName = projectName
	}

	fmt.Printf("Initializing TelemetryFlow integration for project: %s\n", projectName)

	// Create directory structure
	dirs := []string{
		filepath.Join(outputDir, "telemetry"),
		filepath.Join(outputDir, "telemetry", "metrics"),
		filepath.Join(outputDir, "telemetry", "logs"),
		filepath.Join(outputDir, "telemetry", "traces"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Generate files
	generateInitFile()
	generateMetricsFile()
	generateLogsFile()
	generateTracesFile()
	generateConfigFile()
	generateReadme()

	fmt.Println("\n✅ TelemetryFlow integration initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review and update the generated .env file with your API credentials")
	fmt.Println("2. Import the telemetry package in your main.go:")
	fmt.Println("   import \"your-module/telemetry\"")
	fmt.Println("3. Initialize in your main function:")
	fmt.Println("   telemetry.Init()")
	fmt.Println("   defer telemetry.Shutdown()")
}

func runExample(cmd *cobra.Command, args []string) {
	exampleType := args[0]

	fmt.Printf("Generating %s example...\n", exampleType)

	switch exampleType {
	case "basic":
		generateBasicExample()
	case "http-server":
		generateHTTPServerExample()
	case "grpc-server":
		generateGRPCServerExample()
	case "worker":
		generateWorkerExample()
	default:
		fmt.Fprintf(os.Stderr, "Unknown example type: %s\n", exampleType)
		fmt.Println("Available types: basic, http-server, grpc-server, worker")
		os.Exit(1)
	}

	fmt.Println("✅ Example generated successfully!")
}

func runConfig(cmd *cobra.Command, args []string) {
	generateConfigFile()
	fmt.Println("✅ Configuration file generated: .env.telemetryflow")
}

// ===== TEMPLATE DATA STRUCTURES =====

// TemplateData holds all data passed to templates
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

// newTemplateData creates a new TemplateData with current configuration
func newTemplateData() TemplateData {
	return TemplateData{
		ProjectName:   projectName,
		ModulePath:    getModulePath(),
		ServiceName:   serviceName,
		Environment:   "production",
		EnableMetrics: enableMetrics,
		EnableLogs:    enableLogs,
		EnableTraces:  enableTraces,
		APIKeyID:      apiKeyID,
		APIKeySecret:  apiKeySecret,
		Endpoint:      endpoint,
		Port:          "8080",
		NumWorkers:    5,
		QueueSize:     100,
	}
}

// ===== TEMPLATE LOADING =====

// loadTemplate loads a template from file or embedded FS
func loadTemplate(name string) (*template.Template, error) {
	var content []byte
	var err error

	// Try loading from custom template directory first
	if templateDir != "" {
		filePath := filepath.Join(templateDir, name)
		content, err = os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", filePath, err)
		}
	} else {
		// Use embedded templates
		content, err = templateFS.ReadFile("templates/" + name)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded template %s: %w", name, err)
		}
	}

	tmpl, err := template.New(name).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return tmpl, nil
}

// executeTemplate loads and executes a template to a file
func executeTemplate(templateName string, data interface{}, outputPath string) error {
	tmpl, err := loadTemplate(templateName)
	if err != nil {
		return err
	}

	// Sanitize the output path to prevent path traversal (G304)
	cleanPath := filepath.Clean(outputPath)

	f, err := os.Create(cleanPath) // #nosec G304 - path is sanitized above
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", cleanPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", cleanPath, err)
		}
	}()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	fmt.Printf("Generated: %s\n", cleanPath)
	return nil
}

// ===== FILE GENERATORS =====

func generateInitFile() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "telemetry", "init.go")

	if err := executeTemplate("init.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating init.go: %v\n", err)
		os.Exit(1)
	}
}

func generateMetricsFile() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "telemetry", "metrics", "metrics.go")

	if err := executeTemplate("metrics.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating metrics.go: %v\n", err)
		os.Exit(1)
	}
}

func generateLogsFile() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "telemetry", "logs", "logs.go")

	if err := executeTemplate("logs.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating logs.go: %v\n", err)
		os.Exit(1)
	}
}

func generateTracesFile() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "telemetry", "traces", "traces.go")

	if err := executeTemplate("traces.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating traces.go: %v\n", err)
		os.Exit(1)
	}
}

func generateConfigFile() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, ".env.telemetryflow")

	if err := executeTemplate("env.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating .env.telemetryflow: %v\n", err)
		os.Exit(1)
	}
}

func generateReadme() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "telemetry", "README.md")

	if err := executeTemplate("README.md.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating README.md: %v\n", err)
		os.Exit(1)
	}
}

func generateBasicExample() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "example_basic.go")

	if err := executeTemplate("example_basic.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating basic example: %v\n", err)
		os.Exit(1)
	}
}

func generateHTTPServerExample() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "example_http_server.go")

	if err := executeTemplate("example_http_server.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HTTP server example: %v\n", err)
		os.Exit(1)
	}
}

func generateGRPCServerExample() {
	data := newTemplateData()
	data.Port = "50051"
	outputPath := filepath.Join(outputDir, "example_grpc_server.go")

	if err := executeTemplate("example_grpc_server.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating gRPC server example: %v\n", err)
		os.Exit(1)
	}
}

func generateWorkerExample() {
	data := newTemplateData()
	outputPath := filepath.Join(outputDir, "example_worker.go")

	if err := executeTemplate("example_worker.go.tpl", data, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating worker example: %v\n", err)
		os.Exit(1)
	}
}

// ===== HELPER FUNCTIONS =====

func getModulePath() string {
	// In a real implementation, this would read go.mod
	// For now, return a placeholder
	if projectName != "" {
		return strings.ToLower(strings.ReplaceAll(projectName, " ", "-"))
	}
	return "your-module"
}
