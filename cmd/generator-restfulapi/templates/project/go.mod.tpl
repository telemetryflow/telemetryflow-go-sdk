module {{.ModulePath}}

go 1.24

require (
	github.com/go-playground/validator/v10 v10.30.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.14.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
{{- if eq .DBDriver "postgres"}}
	github.com/lib/pq v1.10.9
{{- else if eq .DBDriver "mysql"}}
	github.com/go-sql-driver/mysql v1.8.1
{{- else if eq .DBDriver "sqlite"}}
	github.com/mattn/go-sqlite3 v1.14.22
{{- end}}
{{- if .EnableTelemetry}}
	github.com/telemetryflow/telemetryflow-go-sdk v1.0.1
{{- end}}
{{- if .EnableAuth}}
	github.com/golang-jwt/jwt/v5 v5.3.0
{{- end}}
{{- if .EnableRateLimit}}
	golang.org/x/time v0.14.0
{{- end}}
)
