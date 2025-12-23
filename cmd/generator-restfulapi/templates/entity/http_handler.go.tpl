// Package handler provides HTTP handlers for {{.EntityName}}.
package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"{{.ModulePath}}/internal/application/command"
	"{{.ModulePath}}/internal/application/dto"
	"{{.ModulePath}}/internal/application/handler"
	"{{.ModulePath}}/internal/application/query"
	"{{.ModulePath}}/pkg/response"
)

// {{.EntityName}}Handler handles {{.EntityNameLower}} HTTP requests
type {{.EntityName}}Handler struct {
	commandHandler *handler.{{.EntityName}}CommandHandler
	queryHandler   *handler.{{.EntityName}}QueryHandler
}

// New{{.EntityName}}Handler creates a new {{.EntityNameLower}} handler
func New{{.EntityName}}Handler(
	cmdHandler *handler.{{.EntityName}}CommandHandler,
	qryHandler *handler.{{.EntityName}}QueryHandler,
) *{{.EntityName}}Handler {
	return &{{.EntityName}}Handler{
		commandHandler: cmdHandler,
		queryHandler:   qryHandler,
	}
}

// RegisterRoutes registers {{.EntityNameLower}} routes
func (h *{{.EntityName}}Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/{{.EntityNamePlural}}", h.Create)
	g.GET("/{{.EntityNamePlural}}", h.List)
	g.GET("/{{.EntityNamePlural}}/:id", h.GetByID)
	g.PUT("/{{.EntityNamePlural}}/:id", h.Update)
	g.DELETE("/{{.EntityNamePlural}}/:id", h.Delete)
}

// Create handles POST /{{.EntityNamePlural}}
func (h *{{.EntityName}}Handler) Create(c echo.Context) error {
	var req dto.Create{{.EntityName}}Request
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.Create{{.EntityName}}Command{
{{- range .EntityFields}}
		{{.Name}}: req.{{.Name}},
{{- end}}
	}

	result, err := h.commandHandler.HandleCreate(c.Request().Context(), cmd)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Created(c, result, "{{.EntityName}} created successfully")
}

// List handles GET /{{.EntityNamePlural}}
func (h *{{.EntityName}}Handler) List(c echo.Context) error {
	var q query.List{{.EntityNamePlural | pascal}}Query
	if err := c.Bind(&q); err != nil {
		return response.BadRequest(c, "Invalid query parameters")
	}
	q.Validate()

	result, err := h.queryHandler.HandleList(c.Request().Context(), &q)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Paginated(c, result.Items, result.TotalCount, q.Page, q.PageSize)
}

// GetByID handles GET /{{.EntityNamePlural}}/:id
func (h *{{.EntityName}}Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	q := &query.Get{{.EntityName}}ByIDQuery{ID: id}
	result, err := h.queryHandler.HandleGetByID(c.Request().Context(), q)
	if err != nil {
		return response.NotFound(c, "{{.EntityName}} not found")
	}

	return response.Success(c, result, "")
}

// Update handles PUT /{{.EntityNamePlural}}/:id
func (h *{{.EntityName}}Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	var req dto.Update{{.EntityName}}Request
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cmd := &command.Update{{.EntityName}}Command{
		ID: id,
{{- range .EntityFields}}
		{{.Name}}: req.{{.Name}},
{{- end}}
	}

	result, err := h.commandHandler.HandleUpdate(c.Request().Context(), cmd)
	if err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.Success(c, result, "{{.EntityName}} updated successfully")
}

// Delete handles DELETE /{{.EntityNamePlural}}/:id
func (h *{{.EntityName}}Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, "Invalid ID format")
	}

	cmd := &command.Delete{{.EntityName}}Command{ID: id}
	if err := h.commandHandler.HandleDelete(c.Request().Context(), cmd); err != nil {
		return response.InternalError(c, err.Error())
	}

	return response.NoContent(c)
}
