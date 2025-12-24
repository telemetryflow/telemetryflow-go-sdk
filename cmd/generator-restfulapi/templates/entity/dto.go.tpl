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
{{- if and (ne .DBColumn "created_at") (ne .DBColumn "updated_at")}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}"`
{{- end}}
{{- end}}
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// From{{.EntityName}} converts entity to response DTO
func From{{.EntityName}}(e *entity.{{.EntityName}}) {{.EntityName}}Response {
	return {{.EntityName}}Response{
		ID:        e.ID,
{{- range .EntityFields}}
{{- if and (ne .DBColumn "created_at") (ne .DBColumn "updated_at")}}
		{{.Name}}: e.{{.Name}},
{{- end}}
{{- end}}
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// {{.EntityName}}ToResponse converts entity pointer to response DTO pointer
func {{.EntityName}}ToResponse(e *entity.{{.EntityName}}) *{{.EntityName}}Response {
	if e == nil {
		return nil
	}
	resp := From{{.EntityName}}(e)
	return &resp
}

// From{{.EntityNamePlural | pascal}} converts entities to response DTOs
func From{{.EntityNamePlural | pascal}}(entities []entity.{{.EntityName}}) []{{.EntityName}}Response {
	responses := make([]{{.EntityName}}Response, len(entities))
	for i, e := range entities {
		responses[i] = From{{.EntityName}}(&e)
	}
	return responses
}

// {{.EntityName}}ListResponse represents the list {{.EntityNameLower}} API response
type {{.EntityName}}ListResponse struct {
	Data   []*{{.EntityName}}Response `json:"data"`
	Total  int64                       `json:"total"`
	Offset int                         `json:"offset"`
	Limit  int                         `json:"limit"`
}

// Create{{.EntityName}}Request represents the create {{.EntityNameLower}} request
type Create{{.EntityName}}Request struct {
{{- range .EntityFields}}
{{- if and (ne .DBColumn "created_at") (ne .DBColumn "updated_at")}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
{{- end}}
}

// Update{{.EntityName}}Request represents the update {{.EntityNameLower}} request
type Update{{.EntityName}}Request struct {
{{- range .EntityFields}}
{{- if and (ne .DBColumn "created_at") (ne .DBColumn "updated_at")}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
{{- end}}
}
