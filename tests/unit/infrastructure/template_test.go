// Package infrastructure_test provides unit tests for infrastructure components.
//
// TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
// Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
package infrastructure_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test template functions used by generators

func TestTemplateFunctions(t *testing.T) {
	funcMap := template.FuncMap{
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"title":      strings.Title,
		"pascal":     toPascalCase,
		"camel":      toCamelCase,
		"snake":      toSnakeCase,
		"plural":     pluralize,
		"contains":   strings.Contains,
		"replace":    strings.ReplaceAll,
		"trimSuffix": strings.TrimSuffix,
		"trimPrefix": strings.TrimPrefix,
	}

	t.Run("lower function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{lower .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "UserName"})
		require.NoError(t, err)

		assert.Equal(t, "username", buf.String())
	})

	t.Run("upper function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{upper .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "UserName"})
		require.NoError(t, err)

		assert.Equal(t, "USERNAME", buf.String())
	})

	t.Run("pascal function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{pascal .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "user_name"})
		require.NoError(t, err)

		assert.Equal(t, "UserName", buf.String())
	})

	t.Run("camel function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{camel .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "user_name"})
		require.NoError(t, err)

		assert.Equal(t, "userName", buf.String())
	})

	t.Run("snake function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{snake .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "UserName"})
		require.NoError(t, err)

		assert.Equal(t, "user_name", buf.String())
	})

	t.Run("plural function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{plural .Name}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "user"})
		require.NoError(t, err)

		assert.Equal(t, "users", buf.String())
	})

	t.Run("contains function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{if contains .Name "user"}}yes{{else}}no{{end}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "username"})
		require.NoError(t, err)

		assert.Equal(t, "yes", buf.String())
	})

	t.Run("replace function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{replace .Name "-" "_"}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "user-name"})
		require.NoError(t, err)

		assert.Equal(t, "user_name", buf.String())
	})

	t.Run("trimSuffix function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{trimSuffix .Name "_id"}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "user_id"})
		require.NoError(t, err)

		assert.Equal(t, "user", buf.String())
	})

	t.Run("trimPrefix function", func(t *testing.T) {
		tmpl, err := template.New("test").Funcs(funcMap).Parse(`{{trimPrefix .Name "tbl_"}}`)
		require.NoError(t, err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, map[string]string{"Name": "tbl_users"})
		require.NoError(t, err)

		assert.Equal(t, "users", buf.String())
	})
}

func TestTemplateExecution(t *testing.T) {
	t.Run("should execute template with complete data", func(t *testing.T) {
		tmplStr := `package {{.PackageName}}

type {{.TypeName}} struct {
	ID   string
	Name string
}
`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := map[string]string{
			"PackageName": "entity",
			"TypeName":    "User",
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "package entity")
		assert.Contains(t, result, "type User struct")
	})

	t.Run("should handle missing optional fields", func(t *testing.T) {
		tmplStr := `{{.Required}}{{if .Optional}}-{{.Optional}}{{end}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := map[string]string{
			"Required": "value",
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		require.NoError(t, err)

		assert.Equal(t, "value", buf.String())
	})

	t.Run("should handle conditional blocks", func(t *testing.T) {
		tmplStr := `{{if .EnableFeature}}Feature enabled{{else}}Feature disabled{{end}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		t.Run("feature enabled", func(t *testing.T) {
			data := map[string]bool{"EnableFeature": true}
			var buf bytes.Buffer
			err := tmpl.Execute(&buf, data)
			require.NoError(t, err)
			assert.Equal(t, "Feature enabled", buf.String())
		})

		t.Run("feature disabled", func(t *testing.T) {
			data := map[string]bool{"EnableFeature": false}
			var buf bytes.Buffer
			err := tmpl.Execute(&buf, data)
			require.NoError(t, err)
			assert.Equal(t, "Feature disabled", buf.String())
		})
	})

	t.Run("should handle range blocks", func(t *testing.T) {
		tmplStr := `{{range .Items}}{{.Name}},{{end}}`
		tmpl, err := template.New("test").Parse(tmplStr)
		require.NoError(t, err)

		data := map[string][]map[string]string{
			"Items": {
				{"Name": "Item1"},
				{"Name": "Item2"},
				{"Name": "Item3"},
			},
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, data)
		require.NoError(t, err)

		assert.Equal(t, "Item1,Item2,Item3,", buf.String())
	})
}

func TestTemplateFileOperations(t *testing.T) {
	t.Run("should create output file with correct permissions", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "test", "output.go")

		// Create directory structure
		err := os.MkdirAll(filepath.Dir(outputPath), 0755)
		require.NoError(t, err)

		// Create file
		content := "package test\n"
		err = os.WriteFile(outputPath, []byte(content), 0644)
		require.NoError(t, err)

		// Verify file exists and has correct content
		readContent, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Equal(t, content, string(readContent))

		// Verify file permissions
		info, err := os.Stat(outputPath)
		require.NoError(t, err)
		assert.True(t, info.Mode().Perm()&0644 == 0644)
	})

	t.Run("should overwrite existing file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputPath := filepath.Join(tmpDir, "output.go")

		// Create initial file
		initialContent := "initial content"
		err := os.WriteFile(outputPath, []byte(initialContent), 0644)
		require.NoError(t, err)

		// Overwrite with new content
		newContent := "new content"
		err = os.WriteFile(outputPath, []byte(newContent), 0644)
		require.NoError(t, err)

		// Verify new content
		readContent, err := os.ReadFile(outputPath)
		require.NoError(t, err)
		assert.Equal(t, newContent, string(readContent))
	})

	t.Run("should handle nested directory creation", func(t *testing.T) {
		tmpDir := t.TempDir()
		nestedPath := filepath.Join(tmpDir, "a", "b", "c", "file.go")

		err := os.MkdirAll(filepath.Dir(nestedPath), 0755)
		require.NoError(t, err)

		err = os.WriteFile(nestedPath, []byte("content"), 0644)
		require.NoError(t, err)

		assert.FileExists(t, nestedPath)
	})
}

// Helper functions (same as in main.go)

func toPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	for i, word := range words {
		words[i] = strings.Title(strings.ToLower(word))
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
