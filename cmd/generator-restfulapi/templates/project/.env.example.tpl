# {{.ProjectName}} Environment Configuration

# Server
SERVER_PORT={{.ServerPort}}
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s

# Environment
ENV={{.Environment}}

# Database
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

# JWT Authentication
JWT_SECRET=your_jwt_secret_here_min_32_chars
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=168h
{{- end}}

{{- if .EnableRateLimit}}

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
{{- end}}

{{- if .EnableTelemetry}}

# TelemetryFlow
TELEMETRYFLOW_API_KEY_ID=tfk_your_key_id
TELEMETRYFLOW_API_KEY_SECRET=tfs_your_key_secret
TELEMETRYFLOW_ENDPOINT=api.telemetryflow.id:4317
TELEMETRYFLOW_SERVICE_NAME={{.ServiceName}}
TELEMETRYFLOW_SERVICE_VERSION={{.ServiceVersion}}
{{- end}}

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
