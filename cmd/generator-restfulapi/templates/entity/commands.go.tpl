// Package command contains CQRS commands for {{.EntityName}}.
package command

import (
	"github.com/google/uuid"
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
