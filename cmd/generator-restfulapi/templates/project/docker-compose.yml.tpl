#================================================================================================
# {{.ProjectName}} - Docker Compose Configuration
# RESTful API with DDD + CQRS Pattern
#================================================================================================

#================================================================================================
# NETWORK SETUP
#================================================================================================
networks:
  {{.ProjectName | lower}}_net:
    name: {{.ProjectName | lower}}_net
    driver: bridge

#================================================================================================
# VOLUME SETUP
#================================================================================================
volumes:
{{- if eq .DBDriver "postgres"}}
  vol_postgres_data:
    driver: local
{{- else if eq .DBDriver "mysql"}}
  vol_mysql_data:
    driver: local
{{- end}}

#================================================================================================
# SERVICES
#================================================================================================
services:
  #----------------------------------------------------------------------------------------------
  # DATABASE
  #----------------------------------------------------------------------------------------------
{{- if eq .DBDriver "postgres"}}
  postgres:
    profiles: ["db", "all"]
    platform: linux/amd64
    image: postgres:${POSTGRES_VERSION:-16-alpine}
    container_name: ${CONTAINER_POSTGRES:-{{.ProjectName | lower}}_postgres}
    restart: unless-stopped
    ports:
      - "${PORT_POSTGRES:-{{.DBPort}}}:5432"
    environment:
      - TZ=${TZ:-UTC}
      - POSTGRES_USER=${DB_USER:-{{.DBUser}}}
      - POSTGRES_PASSWORD=${DB_PASSWORD:-password}
      - POSTGRES_DB=${DB_NAME:-{{.DBName}}}
      - PGDATA=/var/lib/postgresql/data/pgdata
    volumes:
      - vol_postgres_data:/var/lib/postgresql/data
    networks:
      - {{.ProjectName | lower}}_net
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-{{.DBUser}}} -d ${DB_NAME:-{{.DBName}}}"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
{{- else if eq .DBDriver "mysql"}}
  mysql:
    profiles: ["db", "all"]
    platform: linux/amd64
    image: mysql:${MYSQL_VERSION:-8.0}
    container_name: ${CONTAINER_MYSQL:-{{.ProjectName | lower}}_mysql}
    restart: unless-stopped
    ports:
      - "${PORT_MYSQL:-{{.DBPort}}}:3306"
    environment:
      - TZ=${TZ:-UTC}
      - MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD:-rootpassword}
      - MYSQL_USER=${DB_USER:-{{.DBUser}}}
      - MYSQL_PASSWORD=${DB_PASSWORD:-password}
      - MYSQL_DATABASE=${DB_NAME:-{{.DBName}}}
    volumes:
      - vol_mysql_data:/var/lib/mysql
    networks:
      - {{.ProjectName | lower}}_net
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
{{- end}}

  #----------------------------------------------------------------------------------------------
  # API - {{.ProjectName}} RESTful API
  #----------------------------------------------------------------------------------------------
  api:
    profiles: ["app", "all"]
    platform: linux/amd64
    build:
      context: .
      dockerfile: Dockerfile
    container_name: ${CONTAINER_API:-{{.ProjectName | lower}}_api}
    restart: unless-stopped
    ports:
      - "${PORT:-{{.ServerPort}}}:{{.ServerPort}}"
    environment:
      - TZ=${TZ:-UTC}

      # Server
      - SERVER_PORT={{.ServerPort}}
      - SERVER_READ_TIMEOUT=${SERVER_READ_TIMEOUT:-15s}
      - SERVER_WRITE_TIMEOUT=${SERVER_WRITE_TIMEOUT:-15s}

      # Database
      - DB_DRIVER={{.DBDriver}}
{{- if eq .DBDriver "postgres"}}
      - DB_HOST=postgres
      - DB_PORT=5432
{{- else if eq .DBDriver "mysql"}}
      - DB_HOST=mysql
      - DB_PORT=3306
{{- end}}
      - DB_NAME=${DB_NAME:-{{.DBName}}}
      - DB_USER=${DB_USER:-{{.DBUser}}}
      - DB_PASSWORD=${DB_PASSWORD:-password}
      - DB_SSL_MODE=${DB_SSL_MODE:-disable}
      - DB_MAX_OPEN_CONNS=${DB_MAX_OPEN_CONNS:-25}
      - DB_MAX_IDLE_CONNS=${DB_MAX_IDLE_CONNS:-5}
      - DB_CONN_MAX_LIFETIME=${DB_CONN_MAX_LIFETIME:-5m}

{{- if .EnableAuth}}
      # JWT
      - JWT_SECRET=${JWT_SECRET}
      - JWT_EXPIRATION=${JWT_EXPIRATION:-24h}
      - JWT_REFRESH_EXPIRATION=${JWT_REFRESH_EXPIRATION:-168h}
{{- end}}

{{- if .EnableRateLimit}}
      # Rate Limiting
      - RATE_LIMIT_REQUESTS=${RATE_LIMIT_REQUESTS:-100}
      - RATE_LIMIT_WINDOW=${RATE_LIMIT_WINDOW:-1m}
{{- end}}

{{- if .EnableTelemetry}}
      # TelemetryFlow / OpenTelemetry
      - TELEMETRYFLOW_API_KEY_ID=${TELEMETRYFLOW_API_KEY_ID}
      - TELEMETRYFLOW_API_KEY_SECRET=${TELEMETRYFLOW_API_KEY_SECRET}
      - TELEMETRYFLOW_ENDPOINT=otel-collector:4317
      - TELEMETRYFLOW_SERVICE_NAME=${TELEMETRYFLOW_SERVICE_NAME:-{{.ServiceName}}}
      - TELEMETRYFLOW_SERVICE_VERSION=${TELEMETRYFLOW_SERVICE_VERSION:-{{.ServiceVersion}}}
      - TELEMETRYFLOW_INSECURE=${TELEMETRYFLOW_INSECURE:-true}
{{- end}}

      # Logging
      - LOG_LEVEL=${LOG_LEVEL:-info}
      - LOG_FORMAT=${LOG_FORMAT:-json}
    depends_on:
{{- if eq .DBDriver "postgres"}}
      postgres:
        condition: service_healthy
{{- else if eq .DBDriver "mysql"}}
      mysql:
        condition: service_healthy
{{- end}}
{{- if .EnableTelemetry}}
      otel-collector:
        condition: service_started
{{- end}}
    networks:
      - {{.ProjectName | lower}}_net
    healthcheck:
      test: ["CMD-SHELL", "wget --no-verbose --tries=1 --spider http://localhost:{{.ServerPort}}/health || exit 1"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

{{- if .EnableTelemetry}}
  #----------------------------------------------------------------------------------------------
  # OTEL COLLECTOR - OpenTelemetry Collector
  #----------------------------------------------------------------------------------------------
  otel-collector:
    profiles: ["monitoring", "all"]
    platform: linux/amd64
    # OTEL Collector Community Contributor
    # image: otel/opentelemetry-collector-contrib:${OTEL_VERSION:-latest}
    # TelemetryFlow Collector (TFO-Collector)
    image: telemetryflow/telemetryflow-collector:${OTEL_VERSION:-latest}
    container_name: ${CONTAINER_OTEL:-{{.ProjectName | lower}}_otel}
    restart: unless-stopped
    command: ["--config=/etc/otel-collector-config.yaml"]
    volumes:
      - ./configs/otel/otel-collector-config.yaml:/etc/otel-collector-config.yaml:ro
    ports:
      - "${PORT_OTEL_GRPC:-4317}:4317"     # OTLP gRPC
      - "${PORT_OTEL_HTTP:-4318}:4318"     # OTLP HTTP
      - "${PORT_OTEL_METRICS:-8889}:8889"  # Prometheus metrics
      - "${PORT_OTEL_HEALTH:-13133}:13133" # Health check
    networks:
      - {{.ProjectName | lower}}_net
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:13133"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
    depends_on:
      - jaeger

  #----------------------------------------------------------------------------------------------
  # JAEGER - Distributed Tracing UI
  #----------------------------------------------------------------------------------------------
  jaeger:
    profiles: ["monitoring", "all"]
    platform: linux/amd64
    image: jaegertracing/jaeger:${JAEGER_VERSION:-2}
    container_name: ${CONTAINER_JAEGER:-{{.ProjectName | lower}}_jaeger}
    restart: unless-stopped
    ports:
      - "${PORT_JAEGER_UI:-16686}:16686"  # Jaeger UI
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    networks:
      - {{.ProjectName | lower}}_net
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:16686"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
{{- end}}
