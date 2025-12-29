# =============================================================================
# TelemetryFlow Go SDK - Dockerfile
# =============================================================================
#
# TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
# Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# =============================================================================
# Multi-stage build for minimal image size
# =============================================================================

# -----------------------------------------------------------------------------
# Stage 1: Builder
# -----------------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

# Build arguments
ARG VERSION=1.1.1
ARG GIT_COMMIT=unknown
ARG GIT_BRANCH=unknown
ARG BUILD_TIME=unknown

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the generators with version information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.Version=${VERSION}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitCommit=${GIT_COMMIT}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitBranch=${GIT_BRANCH}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.BuildTime=${BUILD_TIME}'" \
    -o /telemetryflow-gen ./cmd/generator

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags "-s -w \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.Version=${VERSION}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitCommit=${GIT_COMMIT}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitBranch=${GIT_BRANCH}' \
        -X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.BuildTime=${BUILD_TIME}'" \
    -o /telemetryflow-restapi ./cmd/generator-restfulapi

# Verify binaries
RUN /telemetryflow-gen version && /telemetryflow-restapi version

# -----------------------------------------------------------------------------
# Stage 2: Runtime
# -----------------------------------------------------------------------------
FROM alpine:3.21

# =============================================================================
# TelemetryFlow Metadata Labels (OCI Image Spec)
# =============================================================================
LABEL org.opencontainers.image.title="TelemetryFlow Go SDK" \
      org.opencontainers.image.description="Go SDK and code generators for TelemetryFlow integration - Community Enterprise Observability Platform (CEOP)" \
      org.opencontainers.image.version="1.1.1" \
      org.opencontainers.image.vendor="TelemetryFlow" \
      org.opencontainers.image.authors="DevOpsCorner Indonesia <support@devopscorner.id>" \
      org.opencontainers.image.url="https://telemetryflow.id" \
      org.opencontainers.image.documentation="https://docs.telemetryflow.id" \
      org.opencontainers.image.source="https://github.com/telemetryflow/telemetryflow-go-sdk" \
      org.opencontainers.image.licenses="Apache-2.0" \
      org.opencontainers.image.base.name="alpine:3.21" \
      # TelemetryFlow specific labels
      io.telemetryflow.product="TelemetryFlow Go SDK" \
      io.telemetryflow.component="telemetryflow-sdk" \
      io.telemetryflow.platform="CEOP" \
      io.telemetryflow.maintainer="DevOpsCorner Indonesia"

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && rm -rf /var/cache/apk/*

# Create non-root user and group
RUN addgroup -g 10001 -S telemetryflow && \
    adduser -u 10001 -S telemetryflow -G telemetryflow -h /home/telemetryflow

# Create workspace directory
RUN mkdir -p /workspace && chown -R telemetryflow:telemetryflow /workspace

# Copy binaries from builder
COPY --from=builder /telemetryflow-gen /usr/local/bin/telemetryflow-gen
COPY --from=builder /telemetryflow-restapi /usr/local/bin/telemetryflow-restapi
RUN chmod +x /usr/local/bin/telemetryflow-gen /usr/local/bin/telemetryflow-restapi

# Switch to non-root user
USER telemetryflow

# Set working directory
WORKDIR /workspace

# =============================================================================
# Entrypoint & Command
# =============================================================================
ENTRYPOINT ["/usr/local/bin/telemetryflow-gen"]
CMD ["--help"]

# =============================================================================
# Build Information
# =============================================================================
# Build with:
#   docker build \
#     --build-arg VERSION=1.1.1 \
#     --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
#     --build-arg GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD) \
#     --build-arg BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ') \
#     -t telemetryflow/telemetryflow-sdk:1.1.1 .
#
# Run with:
#   # SDK Generator
#   docker run --rm -v $(pwd):/workspace telemetryflow/telemetryflow-sdk:1.1.1 \
#     init --project myapp --service my-service
#
#   # RESTful API Generator
#   docker run --rm -v $(pwd):/workspace --entrypoint /usr/local/bin/telemetryflow-restapi \
#     telemetryflow/telemetryflow-sdk:1.1.1 \
#     new --name my-api --module github.com/example/my-api
# =============================================================================
