module {{.ModulePath}}

go 1.24.0

toolchain go1.24.11

require (
	github.com/go-playground/validator/v10 v10.30.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.14.0
	github.com/spf13/viper v1.21.0
	github.com/stretchr/testify v1.11.1
{{- if eq .DBDriver "postgres"}}
	gorm.io/driver/postgres v1.5.11
{{- else if eq .DBDriver "mysql"}}
	gorm.io/driver/mysql v1.5.7
{{- else if eq .DBDriver "sqlite"}}
	gorm.io/driver/sqlite v1.5.6
{{- end}}
	gorm.io/gorm v1.25.12
{{- if .EnableTelemetry}}
	github.com/telemetryflow/telemetryflow-go-sdk v1.1.2
{{- end}}
{{- if .EnableAuth}}
	github.com/golang-jwt/jwt/v5 v5.3.0
{{- end}}
{{- if .EnableRateLimit}}
	golang.org/x/time v0.14.0
{{- end}}
)
