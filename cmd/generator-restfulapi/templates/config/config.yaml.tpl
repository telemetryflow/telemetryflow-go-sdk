# {{.ProjectName}} Configuration

server:
  port: "{{.ServerPort}}"
  read_timeout: 15s
  write_timeout: 15s

database:
  driver: {{.DBDriver}}
  host: {{.DBHost}}
  port: "{{.DBPort}}"
  name: {{.DBName}}
  user: {{.DBUser}}
  # password: from environment variable DB_PASSWORD
  ssl_mode: disable
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: 5m

{{- if .EnableAuth}}

jwt:
  # secret: from environment variable JWT_SECRET
  expiration: 24h
  refresh_expiration: 168h
{{- end}}

{{- if .EnableRateLimit}}

ratelimit:
  requests: 100
  window: 1m
{{- end}}

{{- if .EnableTelemetry}}

telemetry:
  # api_key_id: from environment variable TELEMETRYFLOW_API_KEY_ID
  # api_key_secret: from environment variable TELEMETRYFLOW_API_KEY_SECRET
  endpoint: api.telemetryflow.id:4317
  service_name: {{.ServiceName}}
  service_version: {{.ServiceVersion}}
{{- end}}

log:
  level: info
  format: json
