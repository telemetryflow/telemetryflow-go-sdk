// Package main provides the TelemetryFlow RESTful API Generator CLI.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed all:templates
var templateFS embed.FS

var (
	// Project configuration
	projectName    string
	modulePath     string
	serviceName    string
	serviceVersion string
	environment    string

	// Database configuration
	dbDriver string
	dbHost   string
	dbPort   string
	dbName   string
	dbUser   string

	// Server configuration
	serverPort string

	// Feature flags
	enableTelemetry bool
	enableSwagger   bool
	enableCORS      bool
	enableAuth      bool
	enableRateLimit bool

	// Entity configuration
	entityName   string
	entityFields string

	// Output configuration
	outputDir   string
	templateDir string
	noBanner    bool
)

// safePath validates and returns a clean path that is safe to use.
// It ensures the resolved path doesn't escape the base directory through path traversal.
func safePath(baseDir, relativePath string) (string, error) {
	// Clean and resolve the base directory
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base directory: %w", err)
	}

	// Clean and join the paths
	cleanPath := filepath.Clean(filepath.Join(absBase, relativePath))

	// Verify the result is still within the base directory
	if !strings.HasPrefix(cleanPath, absBase) {
		return "", fmt.Errorf("path traversal detected: %s escapes base directory", relativePath)
	}

	return cleanPath, nil
}

// safeReadFile reads a file after validating the path is safe.
// This addresses gosec G304 (potential file inclusion via variable).
func safeReadFile(filePath string) ([]byte, error) {
	// Resolve to absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Read the file
	content, err := os.ReadFile(absPath) // #nosec G304 - path is cleaned and resolved to absolute
	if err != nil {
		return nil, err
	}

	return content, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "telemetryflow-restapi",
		Short: "TelemetryFlow RESTful API Generator",
		Long: `Generate a complete DDD + CQRS pattern RESTful API project with:
  - Domain-Driven Design architecture
  - CQRS (Command Query Responsibility Segregation) pattern
  - Echo web framework
  - Full CRUD operations
  - OpenAPI/Swagger documentation
  - Postman BDD collections
  - ERD and DFD diagrams
  - TelemetryFlow observability integration`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if !noBanner && cmd.Name() != "version" {
				cfg := banner.RESTfulAPIGeneratorConfig()
				cfg.Version = version.Version
				cfg.GitCommit = version.GitCommit
				cfg.BuildTime = version.BuildTime
				cfg.GoVersion = version.GoVersion()
				cfg.Platform = version.Platform()
				banner.PrintCompact(cfg)
			}
		},
	}

	// New command - create new project
	var newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new RESTful API project",
		Long:  "Generate a complete DDD + CQRS RESTful API project structure",
		Run:   runNew,
	}

	// Entity command - add new entity
	var entityCmd = &cobra.Command{
		Use:   "entity",
		Short: "Add a new entity with full CRUD",
		Long:  "Generate domain entity, repository, commands, queries, handlers, and API endpoints",
		Run:   runEntity,
	}

	// Docs command - generate documentation
	var docsCmd = &cobra.Command{
		Use:   "docs",
		Short: "Generate API documentation",
		Long:  "Generate OpenAPI spec, Swagger UI, ERD, DFD, and Postman collection",
		Run:   runDocs,
	}

	// Version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := banner.RESTfulAPIGeneratorConfig()
			cfg.Version = version.Version
			cfg.GitCommit = version.GitCommit
			cfg.BuildTime = version.BuildTime
			cfg.GoVersion = version.GoVersion()
			cfg.Platform = version.Platform()
			banner.Print(cfg)
		},
	}

	// New command flags
	newCmd.Flags().StringVarP(&projectName, "name", "n", "", "Project name (required)")
	newCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Go module path (e.g., github.com/user/project)")
	newCmd.Flags().StringVar(&serviceName, "service", "", "Service name (defaults to project name)")
	newCmd.Flags().StringVar(&serviceVersion, "version", "1.0.0", "Service version")
	newCmd.Flags().StringVar(&environment, "env", "development", "Environment (development, staging, production)")
	newCmd.Flags().StringVar(&dbDriver, "db-driver", "postgres", "Database driver (postgres, mysql, sqlite)")
	newCmd.Flags().StringVar(&dbHost, "db-host", "localhost", "Database host")
	newCmd.Flags().StringVar(&dbPort, "db-port", "5432", "Database port")
	newCmd.Flags().StringVar(&dbName, "db-name", "", "Database name (defaults to project name)")
	newCmd.Flags().StringVar(&dbUser, "db-user", "postgres", "Database user")
	newCmd.Flags().StringVar(&serverPort, "port", "8080", "Server port")
	newCmd.Flags().BoolVar(&enableTelemetry, "telemetry", true, "Enable TelemetryFlow integration")
	newCmd.Flags().BoolVar(&enableSwagger, "swagger", true, "Enable Swagger documentation")
	newCmd.Flags().BoolVar(&enableCORS, "cors", true, "Enable CORS middleware")
	newCmd.Flags().BoolVar(&enableAuth, "auth", true, "Enable JWT authentication")
	newCmd.Flags().BoolVar(&enableRateLimit, "rate-limit", true, "Enable rate limiting")
	newCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	_ = newCmd.MarkFlagRequired("name")

	// Entity command flags
	entityCmd.Flags().StringVarP(&entityName, "name", "n", "", "Entity name (e.g., User, Product)")
	entityCmd.Flags().StringVarP(&entityFields, "fields", "f", "", "Entity fields (e.g., 'name:string,email:string,age:int')")
	entityCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Project root directory")
	_ = entityCmd.MarkFlagRequired("name")

	// Docs command flags
	docsCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Project root directory")

	// Global flags
	rootCmd.PersistentFlags().StringVar(&templateDir, "template-dir", "", "Custom template directory")
	rootCmd.PersistentFlags().BoolVar(&noBanner, "no-banner", false, "Disable banner output")

	rootCmd.AddCommand(newCmd, entityCmd, docsCmd, versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// TemplateData holds all data passed to templates
type TemplateData struct {
	// Project info
	ProjectName    string
	ModulePath     string
	ServiceName    string
	ServiceVersion string
	Environment    string
	EnvPrefix      string

	// Database
	DBDriver string
	DBHost   string
	DBPort   string
	DBName   string
	DBUser   string

	// Server
	ServerPort string

	// Features
	EnableTelemetry bool
	EnableSwagger   bool
	EnableCORS      bool
	EnableAuth      bool
	EnableRateLimit bool

	// Entity (for entity generation)
	EntityName       string
	EntityNameLower  string
	EntityNamePlural string
	EntityFields     []EntityField

	// Entities (for documentation generation)
	Entities []EntityInfo

	// Computed
	Timestamp string
}

// EntityInfo represents entity information for documentation
type EntityInfo struct {
	Name       string
	PluralName string
	Fields     []EntityField
}

// EntityField represents a field in an entity
type EntityField struct {
	Name     string
	Type     string
	JSONName string
	DBColumn string
	GoType   string
	Nullable bool
}

func newTemplateData() TemplateData {
	if serviceName == "" {
		serviceName = projectName
	}
	if dbName == "" {
		dbName = strings.ToLower(strings.ReplaceAll(projectName, "-", "_"))
	}
	if modulePath == "" {
		modulePath = fmt.Sprintf("github.com/example/%s", strings.ToLower(projectName))
	}

	// Generate environment prefix from project name
	envPrefix := strings.ToUpper(strings.ReplaceAll(projectName, "-", "_"))

	return TemplateData{
		ProjectName:     projectName,
		ModulePath:      modulePath,
		ServiceName:     serviceName,
		ServiceVersion:  serviceVersion,
		Environment:     environment,
		EnvPrefix:       envPrefix,
		DBDriver:        dbDriver,
		DBHost:          dbHost,
		DBPort:          dbPort,
		DBName:          dbName,
		DBUser:          dbUser,
		ServerPort:      serverPort,
		EnableTelemetry: enableTelemetry,
		EnableSwagger:   enableSwagger,
		EnableCORS:      enableCORS,
		EnableAuth:      enableAuth,
		EnableRateLimit: enableRateLimit,
	}
}

func runNew(cmd *cobra.Command, args []string) {
	fmt.Printf("Creating new RESTful API project: %s\n", projectName)

	data := newTemplateData()

	// Create directory structure
	dirs := []string{
		// Root
		filepath.Join(outputDir, projectName),
		// CMD
		filepath.Join(outputDir, projectName, "cmd", "api"),
		// Internal - Domain layer
		filepath.Join(outputDir, projectName, "internal", "domain", "entity"),
		filepath.Join(outputDir, projectName, "internal", "domain", "repository"),
		filepath.Join(outputDir, projectName, "internal", "domain", "valueobject"),
		// Internal - Application layer (CQRS)
		filepath.Join(outputDir, projectName, "internal", "application", "command"),
		filepath.Join(outputDir, projectName, "internal", "application", "query"),
		filepath.Join(outputDir, projectName, "internal", "application", "handler"),
		filepath.Join(outputDir, projectName, "internal", "application", "dto"),
		// Internal - Infrastructure layer
		filepath.Join(outputDir, projectName, "internal", "infrastructure", "persistence"),
		filepath.Join(outputDir, projectName, "internal", "infrastructure", "http"),
		filepath.Join(outputDir, projectName, "internal", "infrastructure", "http", "middleware"),
		filepath.Join(outputDir, projectName, "internal", "infrastructure", "http", "handler"),
		filepath.Join(outputDir, projectName, "internal", "infrastructure", "config"),
		// PKG
		filepath.Join(outputDir, projectName, "pkg", "logger"),
		filepath.Join(outputDir, projectName, "pkg", "validator"),
		filepath.Join(outputDir, projectName, "pkg", "response"),
		// Telemetry
		filepath.Join(outputDir, projectName, "telemetry"),
		filepath.Join(outputDir, projectName, "telemetry", "metrics"),
		filepath.Join(outputDir, projectName, "telemetry", "logs"),
		filepath.Join(outputDir, projectName, "telemetry", "traces"),
		// Docs
		filepath.Join(outputDir, projectName, "docs", "api"),
		filepath.Join(outputDir, projectName, "docs", "diagrams"),
		filepath.Join(outputDir, projectName, "docs", "postman"),
		// Configs
		filepath.Join(outputDir, projectName, "configs"),
		// Migrations
		filepath.Join(outputDir, projectName, "migrations"),
		// Scripts
		filepath.Join(outputDir, projectName, "scripts"),
		// Tests
		filepath.Join(outputDir, projectName, "tests", "unit"),
		filepath.Join(outputDir, projectName, "tests", "integration"),
		filepath.Join(outputDir, projectName, "tests", "e2e"),
		filepath.Join(outputDir, projectName, "tests", "mocks"),
		filepath.Join(outputDir, projectName, "tests", "fixtures"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0750); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create directory %s: %v\n", dir, err)
			os.Exit(1)
		}
	}

	// Generate files
	projectRoot := filepath.Join(outputDir, projectName)

	generateFromTemplate("project/go.mod.tpl", data, filepath.Join(projectRoot, "go.mod"))
	generateFromTemplate("project/main.go.tpl", data, filepath.Join(projectRoot, "cmd", "api", "main.go"))
	generateFromTemplate("project/Makefile.tpl", data, filepath.Join(projectRoot, "Makefile"))
	generateFromTemplate("project/README.md.tpl", data, filepath.Join(projectRoot, "README.md"))
	generateFromTemplate("project/Dockerfile.tpl", data, filepath.Join(projectRoot, "Dockerfile"))
	generateFromTemplate("project/docker-compose.yml.tpl", data, filepath.Join(projectRoot, "docker-compose.yml"))
	generateFromTemplate("project/.env.example.tpl", data, filepath.Join(projectRoot, ".env.example"))
	generateFromTemplate("project/.gitignore.tpl", data, filepath.Join(projectRoot, ".gitignore"))

	// Config
	generateFromTemplate("config/config.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "config", "config.go"))
	generateFromTemplate("config/config.yaml.tpl", data, filepath.Join(projectRoot, "configs", "config.yaml"))

	// Domain base
	generateFromTemplate("domain/entity_base.go.tpl", data, filepath.Join(projectRoot, "internal", "domain", "entity", "base.go"))
	generateFromTemplate("domain/repository_base.go.tpl", data, filepath.Join(projectRoot, "internal", "domain", "repository", "base.go"))

	// Application base (CQRS)
	generateFromTemplate("application/command_base.go.tpl", data, filepath.Join(projectRoot, "internal", "application", "command", "base.go"))
	generateFromTemplate("application/query_base.go.tpl", data, filepath.Join(projectRoot, "internal", "application", "query", "base.go"))
	generateFromTemplate("application/handler_base.go.tpl", data, filepath.Join(projectRoot, "internal", "application", "handler", "base.go"))
	generateFromTemplate("application/dto_base.go.tpl", data, filepath.Join(projectRoot, "internal", "application", "dto", "base.go"))

	// Infrastructure - HTTP
	generateFromTemplate("infrastructure/server.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "server.go"))
	generateFromTemplate("infrastructure/router.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "router.go"))
	generateFromTemplate("infrastructure/middleware_logger.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "middleware", "logger.go"))
	generateFromTemplate("infrastructure/middleware_auth.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "middleware", "auth.go"))
	generateFromTemplate("infrastructure/middleware_cors.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "middleware", "cors.go"))
	generateFromTemplate("infrastructure/middleware_ratelimit.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "middleware", "ratelimit.go"))
	generateFromTemplate("infrastructure/health_handler.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "handler", "health.go"))
	generateFromTemplate("infrastructure/home_handler.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "handler", "home.go"))

	// Swagger documentation handlers
	if enableSwagger {
		generateFromTemplate("infrastructure/swagger_handler.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "handler", "swagger.go"))
		generateFromTemplate("infrastructure/swagger_ui.html.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "http", "handler", "swagger_ui.html"))
	}

	// Infrastructure - Persistence
	generateFromTemplate("infrastructure/database.go.tpl", data, filepath.Join(projectRoot, "internal", "infrastructure", "persistence", "database.go"))

	// PKG
	generateFromTemplate("pkg/logger.go.tpl", data, filepath.Join(projectRoot, "pkg", "logger", "logger.go"))
	generateFromTemplate("pkg/validator.go.tpl", data, filepath.Join(projectRoot, "pkg", "validator", "validator.go"))
	generateFromTemplate("pkg/response.go.tpl", data, filepath.Join(projectRoot, "pkg", "response", "response.go"))

	// Create safefile package directory
	if err := os.MkdirAll(filepath.Join(projectRoot, "pkg", "safefile"), 0750); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create safefile directory: %v\n", err)
	}
	generateFromTemplate("pkg/safefile.go.tpl", data, filepath.Join(projectRoot, "pkg", "safefile", "safefile.go"))

	// Telemetry
	if enableTelemetry {
		generateFromTemplate("telemetry/init.go.tpl", data, filepath.Join(projectRoot, "telemetry", "init.go"))
		generateFromTemplate("telemetry/metrics.go.tpl", data, filepath.Join(projectRoot, "telemetry", "metrics", "metrics.go"))
		generateFromTemplate("telemetry/logs.go.tpl", data, filepath.Join(projectRoot, "telemetry", "logs", "logs.go"))
		generateFromTemplate("telemetry/traces.go.tpl", data, filepath.Join(projectRoot, "telemetry", "traces", "traces.go"))
	}

	// Docs
	generateFromTemplate("docs/openapi.yaml.tpl", data, filepath.Join(projectRoot, "docs", "api", "openapi.yaml"))
	generateFromTemplate("docs/swagger.json.tpl", data, filepath.Join(projectRoot, "docs", "api", "swagger.json"))
	generateFromTemplate("docs/embed.go.tpl", data, filepath.Join(projectRoot, "docs", "api", "embed.go"))
	generateFromTemplate("docs/erd.md.tpl", data, filepath.Join(projectRoot, "docs", "diagrams", "ERD.md"))
	generateFromTemplate("docs/dfd.md.tpl", data, filepath.Join(projectRoot, "docs", "diagrams", "DFD.md"))
	generateFromTemplate("docs/postman_collection.json.tpl", data, filepath.Join(projectRoot, "docs", "postman", "collection.json"))
	generateFromTemplate("docs/postman_environment.json.tpl", data, filepath.Join(projectRoot, "docs", "postman", "environment.json"))

	// Migrations
	generateFromTemplate("migrations/000001_init.up.sql.tpl", data, filepath.Join(projectRoot, "migrations", "000001_init.up.sql"))
	generateFromTemplate("migrations/000001_init.down.sql.tpl", data, filepath.Join(projectRoot, "migrations", "000001_init.down.sql"))

	// Scripts
	generateFromTemplate("scripts/run.sh.tpl", data, filepath.Join(projectRoot, "scripts", "run.sh"))
	generateFromTemplate("scripts/test.sh.tpl", data, filepath.Join(projectRoot, "scripts", "test.sh"))

	// Tests
	generateFromTemplate("tests/unit_test.go.tpl", data, filepath.Join(projectRoot, "tests", "unit", "example_test.go"))
	generateFromTemplate("tests/integration_test.go.tpl", data, filepath.Join(projectRoot, "tests", "integration", "api_test.go"))
	generateFromTemplate("tests/e2e_test.go.tpl", data, filepath.Join(projectRoot, "tests", "e2e", "e2e_test.go"))
	generateFromTemplate("tests/mocks.go.tpl", data, filepath.Join(projectRoot, "tests", "mocks", "mocks.go"))
	generateFromTemplate("tests/fixtures.go.tpl", data, filepath.Join(projectRoot, "tests", "fixtures", "fixtures.go"))

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("Project created successfully!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. cd %s\n", projectName)
	fmt.Printf("  2. cp .env.example .env\n")
	fmt.Printf("  3. Edit .env with your configuration\n")
	fmt.Printf("  4. go mod tidy\n")
	fmt.Printf("  5. make run\n")
	fmt.Println("\nTo add a new entity:")
	fmt.Printf("  telemetryflow-restapi entity -n User -f 'name:string,email:string,password:string'\n")
	fmt.Println("\nDocumentation:")
	fmt.Printf("  - OpenAPI: docs/api/openapi.yaml\n")
	fmt.Printf("  - Swagger: docs/api/swagger.json\n")
	fmt.Printf("  - ERD: docs/diagrams/ERD.md\n")
	fmt.Printf("  - DFD: docs/diagrams/DFD.md\n")
	fmt.Printf("  - Postman: docs/postman/collection.json\n")
}

func runEntity(cmd *cobra.Command, args []string) {
	fmt.Printf("Adding entity: %s\n", entityName)

	// Use the module path from the existing go.mod file
	// Use safePath to validate the path is within output directory
	goModPath, pathErr := safePath(outputDir, "go.mod")
	if pathErr != nil {
		fmt.Fprintf(os.Stderr, "Warning: Invalid go.mod path: %v\n", pathErr)
	} else if content, err := safeReadFile(goModPath); err == nil {
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module"))
				break
			}
		}
	}

	data := newTemplateData()
	data.EntityName = toPascalCase(entityName)
	data.EntityNameLower = strings.ToLower(entityName)
	data.EntityNamePlural = pluralize(strings.ToLower(entityName))
	data.EntityFields = parseFields(entityFields)

	// Generate entity files
	generateFromTemplate("entity/entity.go.tpl", data, filepath.Join(outputDir, "internal", "domain", "entity", data.EntityNameLower+".go"))
	generateFromTemplate("entity/repository.go.tpl", data, filepath.Join(outputDir, "internal", "domain", "repository", data.EntityNameLower+"_repository.go"))

	// Generate CQRS files
	generateFromTemplate("entity/commands.go.tpl", data, filepath.Join(outputDir, "internal", "application", "command", data.EntityNameLower+"_commands.go"))
	generateFromTemplate("entity/queries.go.tpl", data, filepath.Join(outputDir, "internal", "application", "query", data.EntityNameLower+"_queries.go"))
	generateFromTemplate("entity/command_handler.go.tpl", data, filepath.Join(outputDir, "internal", "application", "handler", data.EntityNameLower+"_command_handler.go"))
	generateFromTemplate("entity/query_handler.go.tpl", data, filepath.Join(outputDir, "internal", "application", "handler", data.EntityNameLower+"_query_handler.go"))
	generateFromTemplate("entity/dto.go.tpl", data, filepath.Join(outputDir, "internal", "application", "dto", data.EntityNameLower+"_dto.go"))

	// Generate infrastructure files
	generateFromTemplate("entity/persistence.go.tpl", data, filepath.Join(outputDir, "internal", "infrastructure", "persistence", data.EntityNameLower+"_repository.go"))
	generateFromTemplate("entity/http_handler.go.tpl", data, filepath.Join(outputDir, "internal", "infrastructure", "http", "handler", data.EntityNameLower+"_handler.go"))

	// Generate migration
	generateFromTemplate("entity/migration.up.sql.tpl", data, filepath.Join(outputDir, "migrations", fmt.Sprintf("000002_create_%s.up.sql", data.EntityNamePlural)))
	generateFromTemplate("entity/migration.down.sql.tpl", data, filepath.Join(outputDir, "migrations", fmt.Sprintf("000002_create_%s.down.sql", data.EntityNamePlural)))

	fmt.Println("\nEntity created successfully!")
	fmt.Println("\nDon't forget to:")
	fmt.Printf("  1. Register routes in internal/infrastructure/http/router.go\n")
	fmt.Printf("  2. Register repository in dependency injection\n")
	fmt.Printf("  3. Run migrations: make migrate-up\n")
}

func runDocs(cmd *cobra.Command, args []string) {
	fmt.Println("Generating documentation...")

	data := newTemplateData()

	generateFromTemplate("docs/openapi.yaml.tpl", data, filepath.Join(outputDir, "docs", "api", "openapi.yaml"))
	generateFromTemplate("docs/swagger.json.tpl", data, filepath.Join(outputDir, "docs", "api", "swagger.json"))
	generateFromTemplate("docs/erd.md.tpl", data, filepath.Join(outputDir, "docs", "diagrams", "ERD.md"))
	generateFromTemplate("docs/dfd.md.tpl", data, filepath.Join(outputDir, "docs", "diagrams", "DFD.md"))
	generateFromTemplate("docs/postman_collection.json.tpl", data, filepath.Join(outputDir, "docs", "postman", "collection.json"))
	generateFromTemplate("docs/postman_environment.json.tpl", data, filepath.Join(outputDir, "docs", "postman", "environment.json"))

	fmt.Println("Documentation generated successfully!")
}

// Template helpers

func generateFromTemplate(templateName string, data interface{}, outputPath string) {
	tmpl, err := loadTemplate(templateName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Template %s not found, skipping: %v\n", templateName, err)
		return
	}

	// Sanitize the output path to prevent path traversal (G304)
	cleanPath := filepath.Clean(outputPath)

	if err := os.MkdirAll(filepath.Dir(cleanPath), 0750); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create directory for %s: %v\n", cleanPath, err)
		return
	}

	f, err := os.Create(cleanPath) // #nosec G304 - path is sanitized above
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file %s: %v\n", cleanPath, err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close file %s: %v\n", cleanPath, err)
		}
	}()

	if err := tmpl.Execute(f, data); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute template %s: %v\n", templateName, err)
		return
	}

	fmt.Printf("Generated: %s\n", cleanPath)
}

func loadTemplate(name string) (*template.Template, error) {
	var content []byte
	var err error

	titleCaser := cases.Title(language.English)
	funcMap := template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      titleCaser.String,
		"pascal":     toPascalCase,
		"camel":      toCamelCase,
		"snake":      toSnakeCase,
		"plural":     pluralize,
		"contains":   strings.Contains,
		"replace":    strings.ReplaceAll,
		"trimSuffix": strings.TrimSuffix,
		"trimPrefix": strings.TrimPrefix,
		"add":        func(a, b int) int { return a + b },
		"isLast":     func(i int, slice interface{}) bool { return isLastIndex(i, slice) },
	}

	if templateDir != "" {
		// Use safePath to validate the template path
		filePath, pathErr := safePath(templateDir, name)
		if pathErr != nil {
			return nil, fmt.Errorf("invalid template path %s: %w", name, pathErr)
		}
		content, err = safeReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", filePath, err)
		}
	} else {
		content, err = templateFS.ReadFile("templates/" + name)
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded template %s: %w", name, err)
		}
	}

	tmpl, err := template.New(name).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return tmpl, nil
}

func parseFields(fieldsStr string) []EntityField {
	if fieldsStr == "" {
		return nil
	}

	var fields []EntityField
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		kv := strings.Split(strings.TrimSpace(part), ":")
		if len(kv) != 2 {
			continue
		}

		name := strings.TrimSpace(kv[0])
		fieldType := strings.TrimSpace(kv[1])

		fields = append(fields, EntityField{
			Name:     toPascalCase(name),
			Type:     fieldType,
			JSONName: toCamelCase(name),
			DBColumn: toSnakeCase(name),
			GoType:   mapTypeToGo(fieldType),
			Nullable: strings.HasSuffix(fieldType, "?"),
		})
	}

	return fields
}

func mapTypeToGo(t string) string {
	t = strings.TrimSuffix(t, "?")
	switch strings.ToLower(t) {
	case "string", "text":
		return "string"
	case "int", "integer":
		return "int"
	case "int64", "bigint":
		return "int64"
	case "float", "float64", "decimal":
		return "float64"
	case "bool", "boolean":
		return "bool"
	case "time", "datetime", "timestamp":
		return "time.Time"
	case "uuid":
		return "uuid.UUID"
	default:
		return "string"
	}
}

func toPascalCase(s string) string {
	titleCaser := cases.Title(language.English)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, word := range words {
		words[i] = titleCaser.String(strings.ToLower(word))
	}
	return strings.Join(words, "")
}

func toCamelCase(s string) string {
	pascal := toPascalCase(s)
	if len(pascal) == 0 {
		return ""
	}
	return strings.ToLower(pascal[:1]) + pascal[1:]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s + "es"
	}
	if strings.HasSuffix(s, "y") {
		return s[:len(s)-1] + "ies"
	}
	return s + "s"
}

// isLastIndex checks if the index is the last element in a slice
func isLastIndex(i int, slice interface{}) bool {
	switch s := slice.(type) {
	case []EntityField:
		return i == len(s)-1
	case []EntityInfo:
		return i == len(s)-1
	case []string:
		return i == len(s)-1
	default:
		return false
	}
}
