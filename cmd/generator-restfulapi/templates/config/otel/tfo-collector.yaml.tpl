# =============================================================================
# TelemetryFlow Collector Configuration
# =============================================================================
# {{.ProjectName}} - Community Enterprise Observability Platform (CEOP)
# Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
#
# This configuration uses STANDARD OpenTelemetry Collector format for
# compatibility with both build types:
#
#   - Standalone: ./tfo-collector start --config tfo-collector.yaml
#   - OCB:        ./tfo-collector-ocb --config tfo-collector.yaml
#
# Copy to: /etc/tfo-collector/tfo-collector.yaml or ~/.tfo-collector/
#
# =============================================================================
# Environment Variables
# =============================================================================
#   TELEMETRYFLOW_API_KEY_ID      - API Key ID for authentication (tfk_xxx)
#   TELEMETRYFLOW_API_KEY_SECRET  - API Key Secret for authentication (tfs_xxx)
#   TELEMETRYFLOW_ENDPOINT        - TelemetryFlow backend endpoint (host:port)
#   TELEMETRYFLOW_COLLECTOR_ID    - Unique collector identifier (optional)
#   TELEMETRYFLOW_COLLECTOR_NAME  - Human-readable collector name (optional)
#   TELEMETRYFLOW_ENVIRONMENT     - Deployment environment (optional)
#
# =============================================================================

# =============================================================================
# TelemetryFlow Extensions (Standalone-specific, ignored by OCB)
# =============================================================================
telemetryflow:
  api_key_id: "${TELEMETRYFLOW_API_KEY_ID}"
  api_key_secret: "${TELEMETRYFLOW_API_KEY_SECRET}"
  endpoint: "${TELEMETRYFLOW_ENDPOINT:-localhost:4317}"
  tls:
    enabled: true
    insecure_skip_verify: false

collector:
  id: "${TELEMETRYFLOW_COLLECTOR_ID}"
  hostname: ""
  name: "${TELEMETRYFLOW_COLLECTOR_NAME:-{{.ProjectName}} Collector}"
  description: "{{.ProjectName}} - TelemetryFlow Collector"
  version: ""
  tags:
    environment: "${TELEMETRYFLOW_ENVIRONMENT:-production}"
    datacenter: "default"
    service: "{{.ServiceName}}"

# =============================================================================
# RECEIVERS - How telemetry data enters the collector
# =============================================================================
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: "0.0.0.0:4317"
        max_recv_msg_size_mib: 4
        max_concurrent_streams: 100
        read_buffer_size: 524288
        write_buffer_size: 524288
        keepalive:
          server_parameters:
            max_connection_idle: 15s
            max_connection_age: 30s
            max_connection_age_grace: 5s
            time: 10s
            timeout: 5s
      http:
        endpoint: "0.0.0.0:4318"
        cors:
          allowed_origins:
            - "*"
          allowed_headers:
            - "*"
          max_age: 7200

# =============================================================================
# PROCESSORS - How telemetry data is processed
# =============================================================================
processors:
  batch:
    timeout: 200ms
    send_batch_size: 8192
    send_batch_max_size: 0

  memory_limiter:
    check_interval: 1s
    limit_percentage: 80
    spike_limit_percentage: 25

  resource:
    attributes:
      - key: service.namespace
        value: {{.ServiceName}}
        action: upsert
      - key: deployment.environment
        value: "${TELEMETRYFLOW_ENVIRONMENT:-production}"
        action: upsert

# =============================================================================
# CONNECTORS - Pipeline bridging for Exemplars and derived metrics
# =============================================================================
connectors:
  spanmetrics:
    histogram:
      explicit:
        buckets: [1ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s]
    dimensions:
      - name: http.method
        default: GET
      - name: http.status_code
      - name: http.route
      - name: rpc.method
      - name: rpc.service
    exemplars:
      enabled: true
    namespace: traces
    metrics_flush_interval: 15s

  servicegraph:
    latency_histogram_buckets: [1ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s]
    dimensions:
      - http.method
      - http.status_code
    store:
      ttl: 2s
      max_items: 1000
    cache_loop: 1s
    store_expiration_loop: 2s
    virtual_node_peer_attributes:
      - db.system
      - messaging.system
      - rpc.service

# =============================================================================
# EXPORTERS - Where telemetry data is sent
# =============================================================================
exporters:
  debug:
    verbosity: detailed
    sampling_initial: 5
    sampling_thereafter: 200

  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: {{.ServiceName | replace "-" "_"}}
    const_labels:
      collector: tfo-collector
    send_timestamps: true
    metric_expiration: 5m
    enable_open_metrics: true
    resource_to_telemetry_conversion:
      enabled: true

  # OTLP exporter to Jaeger (uncomment to enable)
  # otlp/jaeger:
  #   endpoint: "jaeger:4317"
  #   tls:
  #     insecure: true

  # OTLP exporter to TelemetryFlow backend (uncomment to enable)
  # otlp/tfo:
  #   endpoint: "${TELEMETRYFLOW_ENDPOINT:-localhost:4317}"
  #   tls:
  #     insecure: false
  #   headers:
  #     X-TelemetryFlow-Key-ID: "${TELEMETRYFLOW_API_KEY_ID}"
  #     X-TelemetryFlow-Key-Secret: "${TELEMETRYFLOW_API_KEY_SECRET}"
  #   compression: gzip
  #   timeout: 30s

# =============================================================================
# EXTENSIONS - Additional collector capabilities
# =============================================================================
extensions:
  health_check:
    endpoint: "0.0.0.0:13133"

  zpages:
    endpoint: "0.0.0.0:55679"

  pprof:
    endpoint: "0.0.0.0:1777"

# =============================================================================
# SERVICE - Defines active components and pipelines
# =============================================================================
service:
  extensions: [health_check, zpages, pprof]

  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug, spanmetrics, servicegraph]

    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug, prometheus]

    metrics/spanmetrics:
      receivers: [spanmetrics]
      processors: [memory_limiter, batch]
      exporters: [prometheus]

    metrics/servicegraph:
      receivers: [servicegraph]
      processors: [memory_limiter, batch]
      exporters: [prometheus]

    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug]

  telemetry:
    logs:
      level: info
      encoding: json

    metrics:
      level: detailed
      readers:
        - pull:
            exporter:
              prometheus:
                host: "0.0.0.0"
                port: 8888
