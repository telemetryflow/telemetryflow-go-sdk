{
  "openapi": "3.0.3",
  "info": {
    "title": "{{.ProjectName}} API",
    "description": "{{.ServiceName}} - RESTful API with DDD + CQRS Pattern\n\nThis API provides endpoints for managing resources with full observability support (traces, metrics, logs).\n\n## Architecture\nThis API follows Domain-Driven Design (DDD) with CQRS pattern.\n\n## Authentication\n{{- if .EnableAuth}}API uses JWT Bearer token authentication. Include the token in the Authorization header:\n```\nAuthorization: Bearer <token>\n```{{- else}}Authentication is not enabled for this API.{{- end}}",
    "version": "{{.ServiceVersion}}",
    "contact": {
      "name": "API Support",
      "email": "support@example.com"
    },
    "license": {
      "name": "Apache 2.0",
      "url": "https://www.apache.org/licenses/LICENSE-2.0"
    }
  },
  "servers": [
    {
      "url": "http://localhost:{{.ServerPort}}",
      "description": "Development server"
    },
    {
      "url": "https://api.example.com",
      "description": "Production server"
    }
  ],
  "tags": [
    {
      "name": "Health",
      "description": "Health check endpoints"
    }{{range .Entities}},
    {
      "name": "{{.PluralName}}",
      "description": "{{.Name}} management endpoints"
    }{{end}}
  ],
  "paths": {
    "/health": {
      "get": {
        "tags": ["Health"],
        "summary": "Health check",
        "description": "Check if the service is healthy",
        "operationId": "health",
        "security": [],
        "responses": {
          "200": {
            "description": "Service is healthy",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponse"
                }
              }
            }
          }
        }
      }
    },
    "/ready": {
      "get": {
        "tags": ["Health"],
        "summary": "Readiness check",
        "description": "Check if the service is ready to accept requests",
        "operationId": "ready",
        "security": [],
        "responses": {
          "200": {
            "description": "Service is ready",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/HealthResponse"
                }
              }
            }
          }
        }
      }
    }{{range .Entities}},
    "/api/v1/{{.PluralNameLower}}": {
      "get": {
        "tags": ["{{.PluralName}}"],
        "summary": "List {{.PluralNameLower}}",
        "description": "Get a paginated list of all {{.PluralNameLower}}",
        "operationId": "list{{.PluralName}}",
        "parameters": [
          {
            "name": "limit",
            "in": "query",
            "description": "Number of items to return",
            "schema": {
              "type": "integer",
              "default": 10,
              "minimum": 1,
              "maximum": 100
            }
          },
          {
            "name": "offset",
            "in": "query",
            "description": "Number of items to skip",
            "schema": {
              "type": "integer",
              "default": 0,
              "minimum": 0
            }
          }
        ],
        "responses": {
          "200": {
            "description": "List of {{.PluralNameLower}}",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/{{.Name}}ListResponse"
                }
              }
            }
          },
          "401": {
            "$ref": "#/components/responses/Unauthorized"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      },
      "post": {
        "tags": ["{{.PluralName}}"],
        "summary": "Create {{.NameLower}}",
        "description": "Create a new {{.NameLower}}",
        "operationId": "create{{.Name}}",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Create{{.Name}}Request"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "{{.Name}} created",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SuccessResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "401": {
            "$ref": "#/components/responses/Unauthorized"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      }
    },
    "/api/v1/{{.PluralNameLower}}/{id}": {
      "get": {
        "tags": ["{{.PluralName}}"],
        "summary": "Get {{.NameLower}} by ID",
        "description": "Get a specific {{.NameLower}} by its ID",
        "operationId": "get{{.Name}}ById",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "{{.Name}} ID (UUID)",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "{{.Name}} details",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/{{.Name}}Response"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "401": {
            "$ref": "#/components/responses/Unauthorized"
          },
          "404": {
            "$ref": "#/components/responses/NotFound"
          }
        }
      },
      "put": {
        "tags": ["{{.PluralName}}"],
        "summary": "Update {{.NameLower}}",
        "description": "Update an existing {{.NameLower}}",
        "operationId": "update{{.Name}}",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "{{.Name}} ID (UUID)",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/Update{{.Name}}Request"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "{{.Name}} updated",
            "content": {
              "application/json": {
                "schema": {
                  "$ref": "#/components/schemas/SuccessResponse"
                }
              }
            }
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "401": {
            "$ref": "#/components/responses/Unauthorized"
          },
          "404": {
            "$ref": "#/components/responses/NotFound"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      },
      "delete": {
        "tags": ["{{.PluralName}}"],
        "summary": "Delete {{.NameLower}}",
        "description": "Delete a {{.NameLower}}",
        "operationId": "delete{{.Name}}",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "{{.Name}} ID (UUID)",
            "schema": {
              "type": "string",
              "format": "uuid"
            }
          }
        ],
        "responses": {
          "204": {
            "description": "{{.Name}} deleted"
          },
          "400": {
            "$ref": "#/components/responses/BadRequest"
          },
          "401": {
            "$ref": "#/components/responses/Unauthorized"
          },
          "500": {
            "$ref": "#/components/responses/InternalError"
          }
        }
      }
    }{{end}}
  },
  "components": {
    "schemas": {
      "HealthResponse": {
        "type": "object",
        "properties": {
          "status": {
            "type": "string",
            "example": "healthy"
          },
          "timestamp": {
            "type": "string",
            "format": "date-time"
          },
          "uptime": {
            "type": "string",
            "example": "1h30m45s"
          }
        }
      },
      "SuccessResponse": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": true
          },
          "message": {
            "type": "string",
            "example": "Operation completed successfully"
          },
          "data": {
            "type": "object"
          }
        }
      },
      "ErrorResponse": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": false
          },
          "error": {
            "type": "object",
            "properties": {
              "code": {
                "type": "string",
                "example": "BAD_REQUEST"
              },
              "message": {
                "type": "string",
                "example": "Invalid request"
              },
              "details": {
                "type": "object",
                "additionalProperties": {
                  "type": "string"
                }
              }
            }
          }
        }
      },
      "PaginationMeta": {
        "type": "object",
        "properties": {
          "page": {
            "type": "integer",
            "example": 1
          },
          "page_size": {
            "type": "integer",
            "example": 10
          },
          "total_count": {
            "type": "integer",
            "example": 100
          },
          "total_pages": {
            "type": "integer",
            "example": 10
          }
        }
      }{{range .Entities}},
      "{{.Name}}": {
        "type": "object",
        "properties": {
          "id": {
            "type": "string",
            "format": "uuid",
            "example": "550e8400-e29b-41d4-a716-446655440000"
          }{{range .Fields}},
          "{{.NameSnake}}": {
{{- if eq .Type "string"}}
            "type": "string"{{if .Example}},
            "example": "{{.Example}}"{{end}}
{{- else if eq .Type "int" "int32" "int64"}}
            "type": "integer"{{if .Example}},
            "example": {{.Example}}{{end}}
{{- else if eq .Type "float32" "float64"}}
            "type": "number",
            "format": "double"{{if .Example}},
            "example": {{.Example}}{{end}}
{{- else if eq .Type "bool"}}
            "type": "boolean"{{if .Example}},
            "example": {{.Example}}{{end}}
{{- else if eq .Type "time.Time"}}
            "type": "string",
            "format": "date-time"
{{- else if eq .Type "uuid.UUID"}}
            "type": "string",
            "format": "uuid"
{{- end}}
          }{{end}},
          "created_at": {
            "type": "string",
            "format": "date-time"
          },
          "updated_at": {
            "type": "string",
            "format": "date-time"
          }
        }
      },
      "{{.Name}}Response": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": true
          },
          "data": {
            "$ref": "#/components/schemas/{{.Name}}"
          }
        }
      },
      "{{.Name}}ListResponse": {
        "type": "object",
        "properties": {
          "success": {
            "type": "boolean",
            "example": true
          },
          "data": {
            "type": "array",
            "items": {
              "$ref": "#/components/schemas/{{.Name}}"
            }
          },
          "meta": {
            "$ref": "#/components/schemas/PaginationMeta"
          }
        }
      },
      "Create{{.Name}}Request": {
        "type": "object",
        "required": [{{$first := true}}{{range .Fields}}{{if .Required}}{{if not $first}}, {{end}}"{{.NameSnake}}"{{$first = false}}{{end}}{{end}}],
        "properties": {
{{- range $i, $f := .Fields}}
          "{{$f.NameSnake}}": {
{{- if eq $f.Type "string"}}
            "type": "string"{{if $f.Example}},
            "example": "{{$f.Example}}"{{end}}
{{- else if eq $f.Type "int" "int32" "int64"}}
            "type": "integer"{{if $f.Minimum}},
            "minimum": {{$f.Minimum}}{{end}}{{if $f.Example}},
            "example": {{$f.Example}}{{end}}
{{- else if eq $f.Type "float32" "float64"}}
            "type": "number",
            "format": "double"{{if $f.Example}},
            "example": {{$f.Example}}{{end}}
{{- else if eq $f.Type "bool"}}
            "type": "boolean"{{if $f.Example}},
            "example": {{$f.Example}}{{end}}
{{- else if eq $f.Type "uuid.UUID"}}
            "type": "string",
            "format": "uuid"
{{- end}}
          }{{if not (isLast $i $.Fields)}},{{end}}
{{- end}}
        }
      },
      "Update{{.Name}}Request": {
        "type": "object",
        "required": [{{$first := true}}{{range .Fields}}{{if .Required}}{{if not $first}}, {{end}}"{{.NameSnake}}"{{$first = false}}{{end}}{{end}}],
        "properties": {
{{- range $i, $f := .Fields}}
          "{{$f.NameSnake}}": {
{{- if eq $f.Type "string"}}
            "type": "string"
{{- else if eq $f.Type "int" "int32" "int64"}}
            "type": "integer"{{if $f.Minimum}},
            "minimum": {{$f.Minimum}}{{end}}
{{- else if eq $f.Type "float32" "float64"}}
            "type": "number",
            "format": "double"
{{- else if eq $f.Type "bool"}}
            "type": "boolean"
{{- else if eq $f.Type "uuid.UUID"}}
            "type": "string",
            "format": "uuid"
{{- end}}
          }{{if not (isLast $i $.Fields)}},{{end}}
{{- end}}
        }
      }{{end}}
    },
    "responses": {
      "BadRequest": {
        "description": "Bad request",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            },
            "example": {
              "success": false,
              "error": {
                "code": "BAD_REQUEST",
                "message": "Invalid request body"
              }
            }
          }
        }
      },
      "Unauthorized": {
        "description": "Unauthorized",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            },
            "example": {
              "success": false,
              "error": {
                "code": "UNAUTHORIZED",
                "message": "Invalid or missing authentication token"
              }
            }
          }
        }
      },
      "NotFound": {
        "description": "Resource not found",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            },
            "example": {
              "success": false,
              "error": {
                "code": "NOT_FOUND",
                "message": "Resource not found"
              }
            }
          }
        }
      },
      "InternalError": {
        "description": "Internal server error",
        "content": {
          "application/json": {
            "schema": {
              "$ref": "#/components/schemas/ErrorResponse"
            },
            "example": {
              "success": false,
              "error": {
                "code": "INTERNAL_ERROR",
                "message": "An internal error occurred"
              }
            }
          }
        }
      }
    }{{if .EnableAuth}},
    "securitySchemes": {
      "bearerAuth": {
        "type": "http",
        "scheme": "bearer",
        "bearerFormat": "JWT",
        "description": "Enter your JWT token"
      }
    }{{end}}
  }{{if .EnableAuth}},
  "security": [
    {
      "bearerAuth": []
    }
  ]{{end}}
}
