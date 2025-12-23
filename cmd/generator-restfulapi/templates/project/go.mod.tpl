module {{.ModulePath}}

go 1.22

require (
	github.com/google/uuid v1.6.0
	github.com/labstack/echo/v4 v4.12.0
	github.com/go-playground/validator/v10 v10.22.0
	github.com/spf13/viper v1.19.0
{{- if eq .DBDriver "postgres"}}
	github.com/lib/pq v1.10.9
	github.com/jackc/pgx/v5 v5.6.0
{{- else if eq .DBDriver "mysql"}}
	github.com/go-sql-driver/mysql v1.8.1
{{- else if eq .DBDriver "sqlite"}}
	github.com/mattn/go-sqlite3 v1.14.22
{{- end}}
{{- if .EnableTelemetry}}
	github.com/telemetryflow/telemetryflow-go-sdk v1.0.0
{{- end}}
{{- if .EnableAuth}}
	github.com/golang-jwt/jwt/v5 v5.2.1
{{- end}}
{{- if .EnableRateLimit}}
	golang.org/x/time v0.5.0
{{- end}}
)
