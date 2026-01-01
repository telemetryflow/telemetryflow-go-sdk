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
    ipam:
      config:
        - subnet: 172.152.0.0/16

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
      {{.ProjectName | lower}}_net:
        ipv4_address: ${CONTAINER_IP_POSTGRES:-172.152.152.20}
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
      {{.ProjectName | lower}}_net:
        ipv4_address: ${CONTAINER_IP_MYSQL:-172.152.152.20}
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
      # TelemetryFlow / OpenTelemetry (SDK v1.1.1+)
      - TELEMETRYFLOW_API_KEY_ID=${TELEMETRYFLOW_API_KEY_ID}
      - TELEMETRYFLOW_API_KEY_SECRET=${TELEMETRYFLOW_API_KEY_SECRET}
      - TELEMETRYFLOW_ENDPOINT=otel-collector:4317
      - TELEMETRYFLOW_SERVICE_NAME=${TELEMETRYFLOW_SERVICE_NAME:-{{.ServiceName}}}
      - TELEMETRYFLOW_SERVICE_VERSION=${TELEMETRYFLOW_SERVICE_VERSION:-{{.ServiceVersion}}}
      - TELEMETRYFLOW_SERVICE_NAMESPACE=${TELEMETRYFLOW_SERVICE_NAMESPACE:-telemetryflow}
      - TELEMETRYFLOW_COLLECTOR_ID=${TELEMETRYFLOW_COLLECTOR_ID}
      - TELEMETRYFLOW_ENVIRONMENT=${TELEMETRYFLOW_ENVIRONMENT:-production}
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
      {{.ProjectName | lower}}_net:
        ipv4_address: ${CONTAINER_IP_API:-172.152.152.10}
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
    # =============================================================================
    # OTEL Collector Community Contributor (Standard OTEL format)
    # image: otel/opentelemetry-collector-contrib:${OTEL_VERSION:-0.142.0}
    # command: ["--config=/etc/otelcol-contrib/config.yaml"]
    # =============================================================================
    # TelemetryFlow Collector OCB (TFO-Collector-OCB) - Standard OTEL format
    image: telemetryflow/telemetryflow-collector-ocb:${OTEL_VERSION:-latest}
    command: ["--config=/etc/tfo-collector/otel-collector.yaml"]
    # =============================================================================
    # TelemetryFlow Collector (TFO-Collector) - Custom TFO format (OTLP not implemented yet)
    # image: telemetryflow/telemetryflow-collector:${OTEL_VERSION:-latest}
    # command: ["--config=/etc/tfo-collector/tfo-collector.yaml"]
    # =============================================================================
    container_name: ${CONTAINER_OTEL:-{{.ProjectName | lower}}_otel}
    restart: unless-stopped
    volumes:
      # OTEL Collector Community Contributor config
      # - ./configs/otel/otel-collector.yaml:/etc/otelcol-contrib/config.yaml:ro
      # =============================================================================
      # TelemetryFlow Collector OCB config (Standard OTEL format)
      - ./configs/otel/otel-collector.yaml:/etc/tfo-collector/otel-collector.yaml:ro
      # =============================================================================
      # TelemetryFlow Collector config (Custom TFO format)
      # - ./configs/otel/tfo-collector.yaml:/etc/tfo-collector/tfo-collector.yaml:ro
    ports:
      - "${PORT_OTEL_GRPC:-4317}:4317"     # OTLP gRPC
      - "${PORT_OTEL_HTTP:-4318}:4318"     # OTLP HTTP
      - "${PORT_OTEL_METRICS:-8889}:8889"  # Prometheus metrics
      - "${PORT_OTEL_HEALTH:-13133}:13133" # Health check
    networks:
      {{.ProjectName | lower}}_net:
        ipv4_address: ${CONTAINER_IP_OTEL:-172.152.152.30}
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
    image: jaegertracing/jaeger:${JAEGER_VERSION:-2.13.0}
    container_name: ${CONTAINER_JAEGER:-{{.ProjectName | lower}}_jaeger}
    restart: unless-stopped
    ports:
      - "${PORT_JAEGER_UI:-16686}:16686"  # Jaeger UI
    environment:
      - COLLECTOR_OTLP_ENABLED=true
    networks:
      {{.ProjectName | lower}}_net:
        ipv4_address: ${CONTAINER_IP_JAEGER:-172.152.152.40}
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:16686"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 10s
{{- end}}
