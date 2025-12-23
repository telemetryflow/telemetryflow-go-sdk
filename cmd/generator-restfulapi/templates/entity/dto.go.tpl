// Package dto contains DTOs for {{.EntityName}}.
package dto

import (
	"time"

	"github.com/google/uuid"
	"{{.ModulePath}}/internal/domain/entity"
)

// {{.EntityName}}Response represents the {{.EntityNameLower}} API response
type {{.EntityName}}Response struct {
	ID        uuid.UUID  `json:"id"`
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}"`
{{- end}}
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// From{{.EntityName}} converts entity to response DTO
func From{{.EntityName}}(e *entity.{{.EntityName}}) {{.EntityName}}Response {
	return {{.EntityName}}Response{
		ID:        e.ID,
{{- range .EntityFields}}
		{{.Name}}: e.{{.Name}},
{{- end}}
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// From{{.EntityNamePlural | pascal}} converts entities to response DTOs
func From{{.EntityNamePlural | pascal}}(entities []entity.{{.EntityName}}) []{{.EntityName}}Response {
	responses := make([]{{.EntityName}}Response, len(entities))
	for i, e := range entities {
		responses[i] = From{{.EntityName}}(&e)
	}
	return responses
}

// Create{{.EntityName}}Request represents the create {{.EntityNameLower}} request
type Create{{.EntityName}}Request struct {
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
}

// Update{{.EntityName}}Request represents the update {{.EntityNameLower}} request
type Update{{.EntityName}}Request struct {
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
}
