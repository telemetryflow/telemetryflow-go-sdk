// Package repository defines repository interfaces.
package repository

import (
	"context"

	"github.com/google/uuid"
	"{{.ModulePath}}/internal/domain/entity"
)

// {{.EntityName}}Repository defines the repository interface for {{.EntityName}}
type {{.EntityName}}Repository interface {
	// Create creates a new {{.EntityNameLower}}
	Create(ctx context.Context, e *entity.{{.EntityName}}) error

	// FindByID finds a {{.EntityNameLower}} by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.{{.EntityName}}, error)

	// FindAll finds all {{.EntityNamePlural}} with pagination
	FindAll(ctx context.Context, offset, limit int) ([]entity.{{.EntityName}}, int64, error)

	// Update updates an existing {{.EntityNameLower}}
	Update(ctx context.Context, e *entity.{{.EntityName}}) error

	// Delete soft-deletes a {{.EntityNameLower}} by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// HardDelete permanently deletes a {{.EntityNameLower}}
	HardDelete(ctx context.Context, id uuid.UUID) error
{{- range .EntityFields}}
{{- if eq .Type "string"}}

	// FindBy{{.Name}} finds {{$.EntityNamePlural}} by {{.Name}}
	FindBy{{.Name}}(ctx context.Context, {{.JSONName}} string) ([]entity.{{$.EntityName}}, error)
{{- end}}
{{- end}}
}
