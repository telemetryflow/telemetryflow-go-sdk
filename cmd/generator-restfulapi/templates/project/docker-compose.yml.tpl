version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: {{.ProjectName | lower}}-api
    ports:
      - "{{.ServerPort}}:{{.ServerPort}}"
    environment:
      - SERVER_PORT={{.ServerPort}}
      - DB_DRIVER={{.DBDriver}}
      - DB_HOST=db
      - DB_PORT={{.DBPort}}
      - DB_NAME={{.DBName}}
      - DB_USER={{.DBUser}}
      - DB_PASSWORD=${DB_PASSWORD:-password}
{{- if .EnableTelemetry}}
      - TELEMETRYFLOW_API_KEY_ID=${TELEMETRYFLOW_API_KEY_ID}
      - TELEMETRYFLOW_API_KEY_SECRET=${TELEMETRYFLOW_API_KEY_SECRET}
      - TELEMETRYFLOW_ENDPOINT=${TELEMETRYFLOW_ENDPOINT:-api.telemetryflow.id:4317}
{{- end}}
    depends_on:
      db:
        condition: service_healthy
    networks:
      - {{.ProjectName | lower}}-network
    restart: unless-stopped

{{- if eq .DBDriver "postgres"}}
  db:
    image: postgres:16-alpine
    container_name: {{.ProjectName | lower}}-db
    environment:
      - POSTGRES_USER={{.DBUser}}
      - POSTGRES_PASSWORD=${DB_PASSWORD:-password}
      - POSTGRES_DB={{.DBName}}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "{{.DBPort}}:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U {{.DBUser}} -d {{.DBName}}"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName | lower}}-network
    restart: unless-stopped
{{- else if eq .DBDriver "mysql"}}
  db:
    image: mysql:8.0
    container_name: {{.ProjectName | lower}}-db
    environment:
      - MYSQL_ROOT_PASSWORD=${DB_ROOT_PASSWORD:-rootpassword}
      - MYSQL_USER={{.DBUser}}
      - MYSQL_PASSWORD=${DB_PASSWORD:-password}
      - MYSQL_DATABASE={{.DBName}}
    volumes:
      - mysql_data:/var/lib/mysql
    ports:
      - "{{.DBPort}}:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - {{.ProjectName | lower}}-network
    restart: unless-stopped
{{- end}}

networks:
  {{.ProjectName | lower}}-network:
    driver: bridge

volumes:
{{- if eq .DBDriver "postgres"}}
  postgres_data:
{{- else if eq .DBDriver "mysql"}}
  mysql_data:
{{- end}}
