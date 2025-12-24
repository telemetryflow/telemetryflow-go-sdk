// Package command contains CQRS commands for {{.EntityName}}.
package command

import (
{{- $hasTime := false}}
{{- range .EntityFields}}
{{- if eq .GoType "time.Time"}}
{{- $hasTime = true}}
{{- end}}
{{- end}}
{{- if $hasTime}}
	"time"
{{- end}}

	"github.com/google/uuid"
	"{{.ModulePath}}/internal/domain/entity"
)

// Create{{.EntityName}}Command represents the create {{.EntityNameLower}} command
type Create{{.EntityName}}Command struct {
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
}

// Validate validates the create command
func (c *Create{{.EntityName}}Command) Validate() error {
	// Add validation logic
	return nil
}

// ToEntity converts the command to an entity
func (c *Create{{.EntityName}}Command) ToEntity() *entity.{{.EntityName}} {
	return entity.New{{.EntityName}}({{range $i, $f := .EntityFields}}{{if $i}}, {{end}}c.{{$f.Name}}{{end}})
}

// Update{{.EntityName}}Command represents the update {{.EntityNameLower}} command
type Update{{.EntityName}}Command struct {
	ID uuid.UUID `json:"id" validate:"required"`
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" validate:"required"`
{{- end}}
}

// Validate validates the update command
func (c *Update{{.EntityName}}Command) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// ToEntity converts the command to an entity
func (c *Update{{.EntityName}}Command) ToEntity() *entity.{{.EntityName}} {
	e := entity.New{{.EntityName}}({{range $i, $f := .EntityFields}}{{if $i}}, {{end}}c.{{$f.Name}}{{end}})
	e.ID = c.ID
	return e
}

// Delete{{.EntityName}}Command represents the delete {{.EntityNameLower}} command
type Delete{{.EntityName}}Command struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the delete command
func (c *Delete{{.EntityName}}Command) Validate() error {
	if c.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}
