// Package query contains CQRS queries for {{.EntityName}}.
package query

import (
	"github.com/google/uuid"
)

// Get{{.EntityName}}ByIDQuery represents the get {{.EntityNameLower}} by ID query
type Get{{.EntityName}}ByIDQuery struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

// Validate validates the query
func (q *Get{{.EntityName}}ByIDQuery) Validate() error {
	if q.ID == uuid.Nil {
		return ErrInvalidID
	}
	return nil
}

// List{{.EntityNamePlural | pascal}}Query represents the list {{.EntityNamePlural}} query
type List{{.EntityNamePlural | pascal}}Query struct {
	Page     int    `json:"page" query:"page"`
	PageSize int    `json:"page_size" query:"page_size"`
	SortBy   string `json:"sort_by" query:"sort_by"`
	SortDir  string `json:"sort_dir" query:"sort_dir"`
	Search   string `json:"search" query:"search"`
}

// Validate validates the query
func (q *List{{.EntityNamePlural | pascal}}Query) Validate() error {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 10
	}
	if q.SortDir != "asc" && q.SortDir != "desc" {
		q.SortDir = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	return nil
}

// Offset returns the offset for pagination
func (q *List{{.EntityNamePlural | pascal}}Query) Offset() int {
	return (q.Page - 1) * q.PageSize
}

// GetAll{{.EntityNamePlural}}Query represents the get all {{.EntityNamePlural}} query with pagination
type GetAll{{.EntityNamePlural}}Query struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *GetAll{{.EntityNamePlural}}Query) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}

// Search{{.EntityNamePlural}}Query represents the search {{.EntityNamePlural}} query
type Search{{.EntityNamePlural}}Query struct {
	Query  string `json:"query" query:"query"`
	Offset int    `json:"offset" query:"offset"`
	Limit  int    `json:"limit" query:"limit"`
}

// Validate validates the query
func (q *Search{{.EntityNamePlural}}Query) Validate() error {
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 10
	}
	return nil
}
