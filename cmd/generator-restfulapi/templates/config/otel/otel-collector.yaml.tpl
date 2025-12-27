# OpenTelemetry Collector Configuration (Standard OTEL Format)
# TelemetryFlow Collector - Community Enterprise Observability Platform (CEOP)
# Copyright (c) 2024-2026 TelemetryFlow. All rights reserved.
#
# This config uses STANDARD OpenTelemetry Collector format.
# Use with OCB build: ./build/tfo-collector-ocb --config configs/otel-collector.yaml
# Or with community OTEL collector: otelcol --config configs/otel-collector.yaml
#
# For the TelemetryFlow Standalone build (custom format), use: configs/tfo-collector.yaml

# =============================================================================
# RECEIVERS - How telemetry data enters the collector
# =============================================================================
receivers:
  # OTLP receiver for OpenTelemetry Protocol
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

  # Prometheus scrape receiver (uncomment to enable)
  # prometheus:
  #   config:
  #     scrape_configs:
  #       - job_name: "node-exporter"
  #         scrape_interval: 15s
  #         scrape_timeout: 10s
  #         metrics_path: "/metrics"
  #         static_configs:
  #           - targets:
  #               - "localhost:9100"
  #             labels:
  #               env: "production"

# =============================================================================
# PROCESSORS - How telemetry data is processed
# =============================================================================
processors:
  # Batch processor for efficient data handling
  batch:
    timeout: 200ms
    send_batch_size: 8192
    send_batch_max_size: 0

  # Memory limiter to prevent OOM
  memory_limiter:
    check_interval: 1s
    limit_percentage: 80
    spike_limit_percentage: 25

  # Resource processor for adding attributes
  resource:
    attributes:
      - key: service.namespace
        value: {{.ServiceName}}
        action: upsert
      - key: deployment.environment
        value: production
        action: upsert

  # Attributes processor for data transformation (uncomment to enable)
  # attributes:
  #   actions:
  #     - key: environment
  #       action: insert
  #       value: production

# =============================================================================
# EXPORTERS - Where telemetry data is sent
# =============================================================================
exporters:
  # Debug exporter for development/troubleshooting
  debug:
    verbosity: detailed
    sampling_initial: 5
    sampling_thereafter: 200

  # Prometheus exporter for metrics scraping
  prometheus:
    endpoint: "0.0.0.0:8889"
    namespace: {{.ServiceName | replace "-" "_"}}
    const_labels:
      collector: tfo-collector
    send_timestamps: true
    metric_expiration: 5m
    resource_to_telemetry_conversion:
      enabled: true

  # OTLP exporter for forwarding to other collectors/backends (uncomment to enable)
  # otlp:
  #   endpoint: "your-backend:4317"
  #   tls:
  #     insecure: true
  #   headers:
  #     X-API-Key: "${env:TELEMETRYFLOW_API_KEY}"
  #   compression: gzip
  #   timeout: 30s
  #   retry_on_failure:
  #     enabled: true
  #     initial_interval: 5s
  #     max_interval: 30s
  #     max_elapsed_time: 300s
  #   sending_queue:
  #     enabled: true
  #     num_consumers: 10
  #     queue_size: 1000

  # OTLP exporter to Jaeger
  otlp/jaeger:
    endpoint: "jaeger:4317"  # Jaeger container hostname, NOT 0.0.0.0
    tls:
      insecure: true

  # File exporter for local storage (uncomment to enable)
  # file:
  #   path: "/var/lib/tfo-collector/output.json"
  #   format: json
  #   compression: ""
  #   flush_interval: 1s
  #   rotation:
  #     max_megabytes: 100
  #     max_days: 7
  #     max_backups: 3
  #     localtime: true

# =============================================================================
# EXTENSIONS - Additional collector capabilities
# =============================================================================
extensions:
  # Health check extension
  health_check:
    endpoint: "0.0.0.0:13133"

  # zPages extension for debugging
  zpages:
    endpoint: "0.0.0.0:55679"

  # pprof extension for profiling
  pprof:
    endpoint: "0.0.0.0:1777"

# =============================================================================
# SERVICE - Defines active components and pipelines
# =============================================================================
service:
  extensions: [health_check, zpages, pprof]

  pipelines:
    # Metrics pipeline
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug, prometheus]

    # Logs pipeline
    logs:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug]

    # Traces pipeline
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch, resource]
      exporters: [debug, otlp/jaeger]

  # Internal telemetry configuration
  telemetry:
    logs:
      level: info
      encoding: json

    metrics:
      level: detailed
