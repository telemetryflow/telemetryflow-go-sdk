// Package entity contains domain entities.
package entity

import (
	"time"

	"github.com/google/uuid"
)

// {{.EntityName}} represents the {{.EntityNameLower}} domain entity
type {{.EntityName}} struct {
	Base
{{- range .EntityFields}}
	{{.Name}} {{.GoType}} `json:"{{.JSONName}}" db:"{{.DBColumn}}"`
{{- end}}
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
