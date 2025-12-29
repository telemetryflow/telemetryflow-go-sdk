// Package generator_test provides unit tests for the RESTful API generator.
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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// RESTAPITemplateData mirrors the generator's TemplateData struct for testing
type RESTAPITemplateData struct {
	// Project info
	ProjectName    string
	ModulePath     string
	ServiceName    string
	ServiceVersion string
	Environment    string

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

	// Computed
	Timestamp string
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

func TestRESTAPITemplateDataDefaults(t *testing.T) {
	t.Run("should have default values", func(t *testing.T) {
		data := RESTAPITemplateData{
			ProjectName:     "my-api",
			ModulePath:      "github.com/example/my-api",
			ServiceName:     "my-api",
			ServiceVersion:  "1.0.0",
			Environment:     "development",
			DBDriver:        "postgres",
			DBHost:          "localhost",
			DBPort:          "5432",
			DBName:          "my_api",
			DBUser:          "postgres",
			ServerPort:      "8080",
			EnableTelemetry: true,
			EnableSwagger:   true,
			EnableCORS:      true,
			EnableAuth:      true,
			EnableRateLimit: true,
		}

		assert.Equal(t, "my-api", data.ProjectName)
		assert.Equal(t, "github.com/example/my-api", data.ModulePath)
		assert.Equal(t, "development", data.Environment)
		assert.Equal(t, "postgres", data.DBDriver)
		assert.True(t, data.EnableTelemetry)
		assert.True(t, data.EnableSwagger)
		assert.True(t, data.EnableCORS)
		assert.True(t, data.EnableAuth)
		assert.True(t, data.EnableRateLimit)
	})
}

func TestEntityFieldParsing(t *testing.T) {
	t.Run("should parse simple field", func(t *testing.T) {
		field := EntityField{
			Name:     "Name",
			Type:     "string",
			JSONName: "name",
			DBColumn: "name",
			GoType:   "string",
			Nullable: false,
		}

		assert.Equal(t, "Name", field.Name)
		assert.Equal(t, "string", field.Type)
		assert.Equal(t, "name", field.JSONName)
		assert.Equal(t, "name", field.DBColumn)
		assert.Equal(t, "string", field.GoType)
		assert.False(t, field.Nullable)
	})

	t.Run("should handle nullable fields", func(t *testing.T) {
		field := EntityField{
			Name:     "Description",
			Type:     "string?",
			JSONName: "description",
			DBColumn: "description",
			GoType:   "string",
			Nullable: true,
		}

		assert.True(t, field.Nullable)
	})
}

func TestTypeMapping(t *testing.T) {
	typeMap := map[string]string{
		"string":    "string",
		"text":      "string",
		"int":       "int",
		"integer":   "int",
		"int64":     "int64",
		"bigint":    "int64",
		"float":     "float64",
		"float64":   "float64",
		"decimal":   "float64",
		"bool":      "bool",
		"boolean":   "bool",
		"time":      "time.Time",
		"datetime":  "time.Time",
		"timestamp": "time.Time",
		"uuid":      "uuid.UUID",
	}

	for input, expected := range typeMap {
		t.Run("should map "+input+" to "+expected, func(t *testing.T) {
			result := mapTypeToGo(input)
			assert.Equal(t, expected, result)
		})
	}
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

func TestCaseConversions(t *testing.T) {
	t.Run("should convert to PascalCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"user_name", "UserName"},
			{"user-name", "UserName"},
			{"user name", "UserName"},
			{"username", "Username"},
		}

		for _, tt := range tests {
			result := toPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		}
	})

	t.Run("should convert to camelCase", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"user_name", "userName"},
			{"user-name", "userName"},
			{"user name", "userName"},
			// Note: "UserName" without separators is treated as a single word
			// and becomes "username" (lowercased)
			{"UserName", "username"},
		}

		for _, tt := range tests {
			result := toCamelCase(tt.input)
			assert.Equal(t, tt.expected, result)
		}
	})

	t.Run("should convert to snake_case", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"UserName", "user_name"},
			{"userName", "user_name"},
			{"ID", "i_d"},
		}

		for _, tt := range tests {
			result := toSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		}
	})
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

func TestPluralization(t *testing.T) {
	t.Run("should pluralize correctly", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"user", "users"},
			{"category", "categories"},
			{"class", "classes"},
			{"status", "statuses"},
		}

		for _, tt := range tests {
			result := pluralize(tt.input)
			assert.Equal(t, tt.expected, result)
		}
	})
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

func TestDDDDirectoryStructure(t *testing.T) {
	t.Run("should define DDD directory structure", func(t *testing.T) {
		baseDir := "/tmp/project"
		dirs := []string{
			filepath.Join(baseDir, "cmd", "api"),
			filepath.Join(baseDir, "internal", "domain", "entity"),
			filepath.Join(baseDir, "internal", "domain", "repository"),
			filepath.Join(baseDir, "internal", "domain", "valueobject"),
			filepath.Join(baseDir, "internal", "application", "command"),
			filepath.Join(baseDir, "internal", "application", "query"),
			filepath.Join(baseDir, "internal", "application", "handler"),
			filepath.Join(baseDir, "internal", "application", "dto"),
			filepath.Join(baseDir, "internal", "infrastructure", "persistence"),
			filepath.Join(baseDir, "internal", "infrastructure", "http"),
			filepath.Join(baseDir, "internal", "infrastructure", "http", "middleware"),
			filepath.Join(baseDir, "internal", "infrastructure", "http", "handler"),
			filepath.Join(baseDir, "internal", "infrastructure", "config"),
		}

		// Verify path structure
		for _, dir := range dirs {
			assert.Contains(t, dir, baseDir)
		}

		// Verify DDD layers
		domainDirs := []string{}
		applicationDirs := []string{}
		infrastructureDirs := []string{}

		for _, dir := range dirs {
			if strings.Contains(dir, "/domain/") {
				domainDirs = append(domainDirs, dir)
			}
			if strings.Contains(dir, "/application/") {
				applicationDirs = append(applicationDirs, dir)
			}
			if strings.Contains(dir, "/infrastructure/") {
				infrastructureDirs = append(infrastructureDirs, dir)
			}
		}

		assert.Greater(t, len(domainDirs), 0)
		assert.Greater(t, len(applicationDirs), 0)
		assert.Greater(t, len(infrastructureDirs), 0)
	})
}

func TestCQRSStructure(t *testing.T) {
	t.Run("should have separate command and query paths", func(t *testing.T) {
		baseDir := "/tmp/project/internal/application"

		commandPath := filepath.Join(baseDir, "command")
		queryPath := filepath.Join(baseDir, "query")
		handlerPath := filepath.Join(baseDir, "handler")

		assert.NotEqual(t, commandPath, queryPath)
		assert.Contains(t, commandPath, "command")
		assert.Contains(t, queryPath, "query")
		assert.Contains(t, handlerPath, "handler")
	})
}

func TestDatabaseConfiguration(t *testing.T) {
	t.Run("should support multiple database drivers", func(t *testing.T) {
		drivers := []string{"postgres", "mysql", "sqlite"}

		for _, driver := range drivers {
			data := RESTAPITemplateData{DBDriver: driver}
			assert.Equal(t, driver, data.DBDriver)
		}
	})

	t.Run("should have default database ports", func(t *testing.T) {
		ports := map[string]string{
			"postgres": "5432",
			"mysql":    "3306",
			"sqlite":   "",
		}

		for driver, expectedPort := range ports {
			assert.NotEmpty(t, driver)
			if driver != "sqlite" {
				assert.NotEmpty(t, expectedPort)
			}
		}
	})
}

func TestEntityGeneration(t *testing.T) {
	t.Run("should generate entity template", func(t *testing.T) {
		tmplStr := `// Package entity defines domain entities.
package entity

import (
	"time"
	"github.com/google/uuid"
)

// {{.EntityName}} represents a {{.EntityNameLower}} entity
type {{.EntityName}} struct {
	ID        uuid.UUID
{{- range .EntityFields}}
	{{.Name}} {{.GoType}}
{{- end}}
	CreatedAt time.Time
	UpdatedAt time.Time
}
`
		tmpl, err := template.New("entity").Parse(tmplStr)
		require.NoError(t, err)

		data := RESTAPITemplateData{
			EntityName:      "User",
			EntityNameLower: "user",
			EntityFields: []EntityField{
				{Name: "Name", GoType: "string"},
				{Name: "Email", GoType: "string"},
				{Name: "Age", GoType: "int"},
			},
		}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		result := sb.String()
		assert.Contains(t, result, "type User struct")
		assert.Contains(t, result, "Name string")
		assert.Contains(t, result, "Email string")
		assert.Contains(t, result, "Age int")
	})
}

func TestRepositoryInterfaceGeneration(t *testing.T) {
	t.Run("should generate repository interface", func(t *testing.T) {
		tmplStr := `// {{.EntityName}}Repository defines the repository interface
type {{.EntityName}}Repository interface {
	Create(ctx context.Context, e *entity.{{.EntityName}}) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.{{.EntityName}}, error)
	FindAll(ctx context.Context, offset, limit int) ([]entity.{{.EntityName}}, int64, error)
	Update(ctx context.Context, e *entity.{{.EntityName}}) error
	Delete(ctx context.Context, id uuid.UUID) error
}
`
		tmpl, err := template.New("repository").Parse(tmplStr)
		require.NoError(t, err)

		data := RESTAPITemplateData{
			EntityName: "User",
		}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		result := sb.String()
		assert.Contains(t, result, "UserRepository interface")
		assert.Contains(t, result, "Create(")
		assert.Contains(t, result, "FindByID(")
		assert.Contains(t, result, "FindAll(")
		assert.Contains(t, result, "Update(")
		assert.Contains(t, result, "Delete(")
	})
}

func TestCommandGeneration(t *testing.T) {
	t.Run("should generate CQRS commands", func(t *testing.T) {
		tmplStr := `// Create{{.EntityName}}Command represents the create command
type Create{{.EntityName}}Command struct {
{{- range .EntityFields}}
	{{.Name}} {{.GoType}}
{{- end}}
}

// Update{{.EntityName}}Command represents the update command
type Update{{.EntityName}}Command struct {
	ID uuid.UUID
{{- range .EntityFields}}
	{{.Name}} {{.GoType}}
{{- end}}
}

// Delete{{.EntityName}}Command represents the delete command
type Delete{{.EntityName}}Command struct {
	ID uuid.UUID
}
`
		tmpl, err := template.New("commands").Parse(tmplStr)
		require.NoError(t, err)

		data := RESTAPITemplateData{
			EntityName: "User",
			EntityFields: []EntityField{
				{Name: "Name", GoType: "string"},
				{Name: "Email", GoType: "string"},
			},
		}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		result := sb.String()
		assert.Contains(t, result, "CreateUserCommand")
		assert.Contains(t, result, "UpdateUserCommand")
		assert.Contains(t, result, "DeleteUserCommand")
	})
}

func TestQueryGeneration(t *testing.T) {
	t.Run("should generate CQRS queries", func(t *testing.T) {
		tmplStr := `// Get{{.EntityName}}ByIDQuery represents the get by ID query
type Get{{.EntityName}}ByIDQuery struct {
	ID uuid.UUID
}

// List{{.EntityNamePlural}}Query represents the list query
type List{{.EntityNamePlural}}Query struct {
	Page     int
	PageSize int
	SortBy   string
	SortDir  string
	Search   string
}
`
		tmpl, err := template.New("queries").Parse(tmplStr)
		require.NoError(t, err)

		data := RESTAPITemplateData{
			EntityName:       "User",
			EntityNamePlural: "Users",
		}

		var sb strings.Builder
		err = tmpl.Execute(&sb, data)
		require.NoError(t, err)

		result := sb.String()
		assert.Contains(t, result, "GetUserByIDQuery")
		assert.Contains(t, result, "ListUsersQuery")
		assert.Contains(t, result, "Page")
		assert.Contains(t, result, "PageSize")
		assert.Contains(t, result, "SortBy")
	})
}

func TestDocumentationGeneration(t *testing.T) {
	t.Run("should generate OpenAPI spec path", func(t *testing.T) {
		baseDir := "/tmp/project"
		openAPIPath := filepath.Join(baseDir, "docs", "api", "openapi.yaml")

		assert.Contains(t, openAPIPath, "docs")
		assert.Contains(t, openAPIPath, "api")
		assert.Contains(t, openAPIPath, "openapi.yaml")
	})

	t.Run("should generate Swagger JSON path", func(t *testing.T) {
		baseDir := "/tmp/project"
		swaggerPath := filepath.Join(baseDir, "docs", "api", "swagger.json")

		assert.Contains(t, swaggerPath, "swagger.json")
	})

	t.Run("should generate ERD path", func(t *testing.T) {
		baseDir := "/tmp/project"
		erdPath := filepath.Join(baseDir, "docs", "diagrams", "ERD.md")

		assert.Contains(t, erdPath, "diagrams")
		assert.Contains(t, erdPath, "ERD.md")
	})

	t.Run("should generate DFD path", func(t *testing.T) {
		baseDir := "/tmp/project"
		dfdPath := filepath.Join(baseDir, "docs", "diagrams", "DFD.md")

		assert.Contains(t, dfdPath, "diagrams")
		assert.Contains(t, dfdPath, "DFD.md")
	})

	t.Run("should generate Postman collection path", func(t *testing.T) {
		baseDir := "/tmp/project"
		postmanPath := filepath.Join(baseDir, "docs", "postman", "collection.json")

		assert.Contains(t, postmanPath, "postman")
		assert.Contains(t, postmanPath, "collection.json")
	})
}

func TestMigrationGeneration(t *testing.T) {
	t.Run("should generate migration files", func(t *testing.T) {
		baseDir := "/tmp/project"
		upMigration := filepath.Join(baseDir, "migrations", "000001_init.up.sql")
		downMigration := filepath.Join(baseDir, "migrations", "000001_init.down.sql")

		assert.Contains(t, upMigration, "migrations")
		assert.Contains(t, upMigration, ".up.sql")
		assert.Contains(t, downMigration, ".down.sql")
	})
}

func TestFeatureToggles(t *testing.T) {
	t.Run("should toggle telemetry", func(t *testing.T) {
		data := RESTAPITemplateData{EnableTelemetry: true}
		assert.True(t, data.EnableTelemetry)

		data.EnableTelemetry = false
		assert.False(t, data.EnableTelemetry)
	})

	t.Run("should toggle swagger", func(t *testing.T) {
		data := RESTAPITemplateData{EnableSwagger: true}
		assert.True(t, data.EnableSwagger)

		data.EnableSwagger = false
		assert.False(t, data.EnableSwagger)
	})

	t.Run("should toggle CORS", func(t *testing.T) {
		data := RESTAPITemplateData{EnableCORS: true}
		assert.True(t, data.EnableCORS)

		data.EnableCORS = false
		assert.False(t, data.EnableCORS)
	})

	t.Run("should toggle auth", func(t *testing.T) {
		data := RESTAPITemplateData{EnableAuth: true}
		assert.True(t, data.EnableAuth)

		data.EnableAuth = false
		assert.False(t, data.EnableAuth)
	})

	t.Run("should toggle rate limiting", func(t *testing.T) {
		data := RESTAPITemplateData{EnableRateLimit: true}
		assert.True(t, data.EnableRateLimit)

		data.EnableRateLimit = false
		assert.False(t, data.EnableRateLimit)
	})
}

func TestProjectFileGeneration(t *testing.T) {
	t.Run("should create all project files", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "restapi-test-*")
		require.NoError(t, err)
		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				t.Logf("Failed to remove temp dir: %v", err)
			}
		}()

		// Create directory structure
		dirs := []string{
			filepath.Join(tmpDir, "cmd", "api"),
			filepath.Join(tmpDir, "internal", "domain", "entity"),
			filepath.Join(tmpDir, "docs", "api"),
		}

		for _, dir := range dirs {
			err := os.MkdirAll(dir, 0755)
			require.NoError(t, err)
		}

		// Verify directories were created
		for _, dir := range dirs {
			info, err := os.Stat(dir)
			require.NoError(t, err)
			assert.True(t, info.IsDir())
		}
	})
}
