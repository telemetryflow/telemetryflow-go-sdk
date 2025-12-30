// Package entity contains domain entities.
package entity

import (
{{- $hasTime := false}}
{{- $hasUUID := false}}
{{- range .EntityFields}}
{{- if eq .GoType "time.Time"}}
{{- $hasTime = true}}
{{- end}}
{{- if eq .GoType "uuid.UUID"}}
{{- $hasUUID = true}}
{{- end}}
{{- end}}
{{- if $hasTime}}
	"time"
{{- end}}

{{- if $hasUUID}}
	"github.com/google/uuid"
{{- end}}
)

// {{.EntityName}} represents the {{.EntityNameLower}} domain entity
type {{.EntityName}} struct {
	Base
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" gorm:"{{if eq .GoType "uuid.UUID"}}type:uuid;{{end}}{{if eq .DBColumn "status"}}type:varchar(50);{{end}}not null{{if eq .Type "string"}};index{{end}}"`
{{- end}}
}

// TableName returns the table name for GORM
func ({{.EntityName}}) TableName() string {
	return "{{snake .EntityNamePlural}}"
}

// New{{.EntityName}} creates a new {{.EntityName}} entity
func New{{.EntityName}}({{range $i, $f := .EntityFields}}{{if $i}}, {{end}}{{$f.JSONName}} {{$f.GoType}}{{end}}) *{{.EntityName}} {
	return &{{.EntityName}}{
		Base: NewBase(),
{{- range .EntityFields}}
		{{.Name}}: {{.JSONName}},
{{- end}}
	}
}

// Update updates the {{.EntityNameLower}} fields
func (e *{{.EntityName}}) Update({{range $i, $f := .EntityFields}}{{if $i}}, {{end}}{{$f.JSONName}} {{$f.GoType}}{{end}}) {
{{- range .EntityFields}}
	e.{{.Name}} = {{.JSONName}}
{{- end}}
	e.MarkUpdated()
}

// Validate validates the entity
func (e *{{.EntityName}}) Validate() error {
	// Add validation logic here
	return nil
}
