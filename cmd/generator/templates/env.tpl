# =============================================================================
# TelemetryFlow Configuration
# =============================================================================
# Generated for: {{.ProjectName}}
# Compatible with TFO-Collector v1.1.2 (OCB-native)
# =============================================================================

# -----------------------------------------------------------------------------
# API Credentials
# -----------------------------------------------------------------------------
TELEMETRYFLOW_API_KEY_ID={{.APIKeyID}}
TELEMETRYFLOW_API_KEY_SECRET={{.APIKeySecret}}

# -----------------------------------------------------------------------------
# Endpoint Configuration
# -----------------------------------------------------------------------------
TELEMETRYFLOW_ENDPOINT={{.Endpoint}}
TELEMETRYFLOW_PROTOCOL=grpc
TELEMETRYFLOW_INSECURE=true

# -----------------------------------------------------------------------------
# Service Information
# -----------------------------------------------------------------------------
TELEMETRYFLOW_SERVICE_NAME={{.ServiceName}}
TELEMETRYFLOW_SERVICE_VERSION={{.ServiceVersion}}
TELEMETRYFLOW_SERVICE_NAMESPACE=telemetryflow
TELEMETRYFLOW_ENVIRONMENT={{.Environment}}

# -----------------------------------------------------------------------------
# TFO v2 API Settings (aligned with tfoexporter)
# -----------------------------------------------------------------------------
TELEMETRYFLOW_USE_V2_API=true
TELEMETRYFLOW_V2_ONLY=false

# -----------------------------------------------------------------------------
# Collector Identity (aligned with tfoidentityextension)
# -----------------------------------------------------------------------------
TELEMETRYFLOW_COLLECTOR_ID=
TELEMETRYFLOW_COLLECTOR_NAME={{.ProjectName}} SDK Client
TELEMETRYFLOW_DATACENTER=default
TELEMETRYFLOW_ENRICH_RESOURCES=true

# -----------------------------------------------------------------------------
# Signals Configuration
# -----------------------------------------------------------------------------
TELEMETRYFLOW_ENABLE_TRACES={{.EnableTraces}}
TELEMETRYFLOW_ENABLE_METRICS={{.EnableMetrics}}
TELEMETRYFLOW_ENABLE_LOGS={{.EnableLogs}}
TELEMETRYFLOW_ENABLE_EXEMPLARS=true

# -----------------------------------------------------------------------------
# Performance Settings
# -----------------------------------------------------------------------------
TELEMETRYFLOW_TIMEOUT=10
TELEMETRYFLOW_BATCH_TIMEOUT=5000
TELEMETRYFLOW_BATCH_MAX_SIZE=512
TELEMETRYFLOW_COMPRESSION=false
TELEMETRYFLOW_RATE_LIMIT=0

# -----------------------------------------------------------------------------
# Retry Settings
# -----------------------------------------------------------------------------
TELEMETRYFLOW_RETRY_ENABLED=true
TELEMETRYFLOW_MAX_RETRIES=3
TELEMETRYFLOW_RETRY_BACKOFF=500
