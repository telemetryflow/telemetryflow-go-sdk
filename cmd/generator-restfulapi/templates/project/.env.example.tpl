# {{.ProjectName}} Environment Configuration
# Copy this file to .env and update the values

#================================================================================================
# GENERAL
#================================================================================================
TZ=UTC
ENV={{.Environment}}

#================================================================================================
# SERVER
#================================================================================================
SERVER_PORT={{.ServerPort}}
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

#================================================================================================
# DATABASE ({{.DBDriver | upper}})
#================================================================================================
DB_DRIVER={{.DBDriver}}
DB_HOST={{.DBHost}}
DB_PORT={{.DBPort}}
DB_NAME={{.DBName}}
DB_USER={{.DBUser}}
DB_PASSWORD=your_password_here
DB_SSL_MODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m

{{- if .EnableAuth}}

#================================================================================================
# JWT AUTHENTICATION
#================================================================================================
JWT_SECRET=your_jwt_secret_here_min_64_chars_recommended
JWT_REFRESH_SECRET=your_refresh_secret_here_min_64_chars_recommended
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h
SESSION_SECRET=your_session_secret_here_min_64_chars_recommended
{{- end}}

{{- if .EnableRateLimit}}

#================================================================================================
# RATE LIMITING
#================================================================================================
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
{{- end}}

{{- if .EnableTelemetry}}

#================================================================================================
# TELEMETRYFLOW / OPENTELEMETRY
#================================================================================================
TELEMETRYFLOW_API_KEY_ID=tfk_your_key_id
TELEMETRYFLOW_API_KEY_SECRET=tfs_your_key_secret
TELEMETRYFLOW_ENDPOINT=localhost:4317
TELEMETRYFLOW_SERVICE_NAME={{.ServiceName}}
TELEMETRYFLOW_SERVICE_VERSION={{.ServiceVersion}}
TELEMETRYFLOW_INSECURE=true
{{- end}}

#================================================================================================
# LOGGING
#================================================================================================
LOG_LEVEL=info
LOG_FORMAT=json

#================================================================================================
# DOCKER COMPOSE - Container Settings
#================================================================================================
{{- if eq .DBDriver "postgres"}}
# PostgreSQL
POSTGRES_VERSION=16-alpine
CONTAINER_POSTGRES={{.ProjectName | lower}}_postgres
PORT_POSTGRES={{.DBPort}}
{{- else if eq .DBDriver "mysql"}}
# MySQL
MYSQL_VERSION=8.0
CONTAINER_MYSQL={{.ProjectName | lower}}_mysql
PORT_MYSQL={{.DBPort}}
{{- end}}

# API
CONTAINER_API={{.ProjectName | lower}}_api
PORT={{.ServerPort}}

{{- if .EnableTelemetry}}

# OpenTelemetry Collector
OTEL_VERSION=latest
CONTAINER_OTEL={{.ProjectName | lower}}_otel
PORT_OTEL_GRPC=4317
PORT_OTEL_HTTP=4318
PORT_OTEL_METRICS=8889
PORT_OTEL_HEALTH=13133

# Jaeger
JAEGER_VERSION=2
CONTAINER_JAEGER={{.ProjectName | lower}}_jaeger
PORT_JAEGER_UI=16686
{{- end}}
