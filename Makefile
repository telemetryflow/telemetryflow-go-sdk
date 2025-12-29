# TelemetryFlow Go SDK - Makefile
#
# TelemetryFlow Go SDK - Community Enterprise Observability Platform (CEOP)
# Copyright (c) 2024-2026 DevOpsCorner Indonesia. All rights reserved.
#
# Build and development commands for TelemetryFlow Go SDK

# Build configuration
PRODUCT_NAME := TelemetryFlow Go SDK
VERSION ?= 1.1.1
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GO_VERSION := $(shell go version | cut -d ' ' -f 3)

# Binary names
GENERATOR_NAME := telemetryflow-gen
RESTAPI_GENERATOR_NAME := telemetryflow-restapi

# Directories
BUILD_DIR := ./build
DIST_DIR := ./dist
GENERATOR_PATH := ./cmd/generator
RESTAPI_GENERATOR_PATH := ./cmd/generator-restfulapi
EXAMPLE_PATH := ./examples/basic

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Build flags
LDFLAGS := -s -w \
	-X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.Version=$(VERSION)' \
	-X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.GitBranch=$(GIT_BRANCH)' \
	-X 'github.com/telemetryflow/telemetryflow-go-sdk/internal/version.BuildTime=$(BUILD_TIME)'

# Platforms for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
BLUE := \033[0;34m
NC := \033[0m

.PHONY: all build build-all build-sdk build-generators build-gen build-restapi build-linux build-darwin clean test test-unit test-integration test-e2e test-all test-coverage test-short bench deps deps-update lint lint-fix fmt vet run-gen run-restapi run-example install uninstall help tidy verify version check ci release-check docs godoc docker-build docker-push generate-example generate-restapi-example

# Default target
all: build

# Help target
help:
	@echo "$(GREEN)$(PRODUCT_NAME) - Build System$(NC)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(NC)"
	@echo "  make                  - Build SDK and generators (default)"
	@echo "  make build            - Build SDK and generators"
	@echo "  make build-sdk        - Build SDK only"
	@echo "  make build-generators - Build all generators"
	@echo "  make build-gen        - Build telemetryflow-gen"
	@echo "  make build-restapi    - Build telemetryflow-restapi"
	@echo "  make build-all        - Build generators for all platforms"
	@echo "  make build-linux      - Build generators for Linux (amd64 and arm64)"
	@echo "  make build-darwin     - Build generators for macOS (amd64 and arm64)"
	@echo ""
	@echo "$(YELLOW)Development Commands:$(NC)"
	@echo "  make run-gen          - Run telemetryflow-gen"
	@echo "  make run-restapi      - Run telemetryflow-restapi"
	@echo "  make run-example      - Run basic example"
	@echo "  make generate-example - Generate SDK example project"
	@echo "  make generate-restapi-example - Generate RESTful API example"
	@echo ""
	@echo "$(YELLOW)Testing Commands:$(NC)"
	@echo "  make test             - Run unit and integration tests"
	@echo "  make test-e2e         - Run E2E tests only"
	@echo "  make test-integration - Run integration tests only"
	@echo "  make test-unit        - Run unit tests only"
	@echo "  make test-all         - Run all tests"
	@echo "  make test-coverage    - Run tests with coverage report"
	@echo "  make test-short       - Run short tests"
	@echo "  make bench            - Run benchmarks"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make lint             - Run linter"
	@echo "  make lint-fix         - Run linter with auto-fix"
	@echo "  make fmt              - Format code"
	@echo "  make vet              - Run go vet"
	@echo "  make check            - Run all checks (fmt, vet, lint, test)"
	@echo ""
	@echo "$(YELLOW)Dependencies:$(NC)"
	@echo "  make deps             - Download dependencies"
	@echo "  make deps-update      - Update dependencies"
	@echo "  make tidy             - Tidy go modules"
	@echo "  make verify           - Verify dependencies"
	@echo ""
	@echo "$(YELLOW)Other Commands:$(NC)"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make install          - Install generators globally"
	@echo "  make uninstall        - Uninstall generators"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make version          - Show version information"
	@echo "  make ci               - Run CI pipeline"
	@echo "  make docs             - Show documentation locations"
	@echo "  make godoc            - Start godoc server"
	@echo ""
	@echo "$(YELLOW)Configuration:$(NC)"
	@echo "  VERSION=$(VERSION)"
	@echo "  GIT_COMMIT=$(GIT_COMMIT)"
	@echo "  GIT_BRANCH=$(GIT_BRANCH)"
	@echo "  BUILD_TIME=$(BUILD_TIME)"
	@echo "  GO_VERSION=$(GO_VERSION)"

## Build commands
build: build-sdk build-generators
	@echo "$(GREEN)Build complete$(NC)"

build-sdk:
	@echo "$(GREEN)Building $(PRODUCT_NAME)...$(NC)"
	@$(GOBUILD) -v ./...
	@echo "$(GREEN)SDK build complete$(NC)"

build-generators: build-gen build-restapi
	@echo "$(GREEN)Generators build complete$(NC)"

build-gen:
	@echo "$(GREEN)Building $(GENERATOR_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(GENERATOR_NAME) $(GENERATOR_PATH)
	@echo "$(GREEN)Built: $(BUILD_DIR)/$(GENERATOR_NAME)$(NC)"

build-restapi:
	@echo "$(GREEN)Building $(RESTAPI_GENERATOR_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(RESTAPI_GENERATOR_NAME) $(RESTAPI_GENERATOR_PATH)
	@echo "$(GREEN)Built: $(BUILD_DIR)/$(RESTAPI_GENERATOR_NAME)$(NC)"

build-all:
	@echo "$(GREEN)Building generators for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} ; \
		echo "$(YELLOW)Building for $${GOOS}/$${GOARCH}...$(NC)" ; \
		output="$(DIST_DIR)/$(GENERATOR_NAME)-$${GOOS}-$${GOARCH}" ; \
		if [ "$${GOOS}" = "windows" ]; then output="$${output}.exe"; fi ; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} $(GOBUILD) -ldflags "$(LDFLAGS)" -o $${output} $(GENERATOR_PATH) ; \
		output="$(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-$${GOOS}-$${GOARCH}" ; \
		if [ "$${GOOS}" = "windows" ]; then output="$${output}.exe"; fi ; \
		GOOS=$${GOOS} GOARCH=$${GOARCH} $(GOBUILD) -ldflags "$(LDFLAGS)" -o $${output} $(RESTAPI_GENERATOR_PATH) ; \
	done
	@echo "$(GREEN)All platform builds complete in $(DIST_DIR)$(NC)"

build-linux:
	@echo "$(GREEN)Building generators for Linux...$(NC)"
	@mkdir -p $(DIST_DIR)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(GENERATOR_NAME)-linux-amd64 $(GENERATOR_PATH)
	@GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(GENERATOR_NAME)-linux-arm64 $(GENERATOR_PATH)
	@GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-linux-amd64 $(RESTAPI_GENERATOR_PATH)
	@GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-linux-arm64 $(RESTAPI_GENERATOR_PATH)
	@echo "$(GREEN)Linux builds complete$(NC)"

build-darwin:
	@echo "$(GREEN)Building generators for macOS...$(NC)"
	@mkdir -p $(DIST_DIR)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(GENERATOR_NAME)-darwin-amd64 $(GENERATOR_PATH)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(GENERATOR_NAME)-darwin-arm64 $(GENERATOR_PATH)
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-darwin-amd64 $(RESTAPI_GENERATOR_PATH)
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-darwin-arm64 $(RESTAPI_GENERATOR_PATH)
	@echo "$(GREEN)macOS builds complete$(NC)"

## Development commands
run-gen: build-gen
	@echo "$(GREEN)Running $(GENERATOR_NAME)...$(NC)"
	@$(BUILD_DIR)/$(GENERATOR_NAME) --help

run-restapi: build-restapi
	@echo "$(GREEN)Running $(RESTAPI_GENERATOR_NAME)...$(NC)"
	@$(BUILD_DIR)/$(RESTAPI_GENERATOR_NAME) --help

run-example:
	@echo "$(YELLOW)Make sure to set environment variables first$(NC)"
	@echo "$(GREEN)Running basic example...$(NC)"
	@$(GOCMD) run $(EXAMPLE_PATH)/main.go

generate-example: build-gen
	@echo "$(GREEN)Generating SDK example project...$(NC)"
	@mkdir -p _generated
	@$(BUILD_DIR)/$(GENERATOR_NAME) init \
		--project "sample-app" \
		--service "sample-service" \
		--output "_generated"
	@echo "$(GREEN)Example generated in _generated/$(NC)"

generate-restapi-example: build-restapi
	@echo "$(GREEN)Generating RESTful API example project...$(NC)"
	@mkdir -p _generated
	@$(BUILD_DIR)/$(RESTAPI_GENERATOR_NAME) new \
		--name "sample-api" \
		--module "github.com/example/sample-api" \
		--output "_generated"
	@echo "$(GREEN)RESTful API example generated in _generated/$(NC)"

## Test commands
test: test-unit test-integration
	@echo "$(GREEN)All tests completed$(NC)"

test-unit:
	@echo "$(GREEN)Running unit tests...$(NC)"
	@$(GOTEST) -v -race -short -coverprofile=coverage-unit.out ./tests/unit/...

test-integration:
	@echo "$(GREEN)Running integration tests...$(NC)"
	@$(GOTEST) -v -race -coverprofile=coverage-integration.out ./tests/integration/...

test-e2e:
	@echo "$(GREEN)Running E2E tests...$(NC)"
	@$(GOTEST) -v -timeout 10m ./tests/e2e/...

test-all: test-unit test-integration test-e2e
	@echo "$(GREEN)All tests completed$(NC)"

test-coverage: test
	@echo "$(GREEN)Generating coverage reports...$(NC)"
	@$(GOCMD) tool cover -html=coverage-unit.out -o coverage-unit.html 2>/dev/null || true
	@$(GOCMD) tool cover -html=coverage-integration.out -o coverage-integration.html 2>/dev/null || true
	@echo "$(GREEN)Coverage reports generated$(NC)"

test-short:
	@echo "$(GREEN)Running short tests...$(NC)"
	@$(GOTEST) -v -race -short ./...

bench:
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@$(GOTEST) -bench=. -benchmem ./...

## Dependencies
deps:
	@echo "$(GREEN)Downloading dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "$(GREEN)Dependencies downloaded$(NC)"

deps-update:
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@$(GOGET) -u ./...
	@$(GOMOD) tidy
	@echo "$(GREEN)Dependencies updated$(NC)"

tidy:
	@echo "$(GREEN)Tidying go modules...$(NC)"
	@$(GOMOD) tidy
	@echo "$(GREEN)Go modules tidied$(NC)"

verify:
	@echo "$(GREEN)Verifying dependencies...$(NC)"
	@$(GOMOD) verify
	@echo "$(GREEN)Dependencies verified$(NC)"

## Code quality
lint:
	@echo "$(GREEN)Running linter...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)"; \
	fi

lint-fix:
	@echo "$(GREEN)Running linter with auto-fix...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --fix ./...; \
	else \
		echo "$(YELLOW)golangci-lint not installed$(NC)"; \
	fi

fmt:
	@echo "$(GREEN)Formatting code...$(NC)"
	@$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

vet:
	@echo "$(GREEN)Running go vet...$(NC)"
	@$(GOCMD) vet ./...
	@echo "$(GREEN)Vet complete$(NC)"

check: fmt vet lint test
	@echo "$(GREEN)All checks passed$(NC)"

## Cleanup
clean:
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -rf _generated
	@rm -rf bin/
	@rm -f coverage*.out coverage*.html
	@echo "$(GREEN)Clean complete$(NC)"

## Installation
install: build-generators
	@echo "$(GREEN)Installing generators...$(NC)"
	@$(GOCMD) install -ldflags "$(LDFLAGS)" $(GENERATOR_PATH)
	@$(GOCMD) install -ldflags "$(LDFLAGS)" $(RESTAPI_GENERATOR_PATH)
	@echo "$(GREEN)Installed: $(GENERATOR_NAME) and $(RESTAPI_GENERATOR_NAME)$(NC)"

uninstall:
	@echo "$(GREEN)Removing installed generators...$(NC)"
	@rm -f $(GOPATH)/bin/$(GENERATOR_NAME)
	@rm -f $(GOPATH)/bin/$(RESTAPI_GENERATOR_NAME)
	@echo "$(GREEN)Uninstalled$(NC)"

## CI pipeline
ci: deps check
	@echo "$(GREEN)CI pipeline completed$(NC)"

release-check:
	@echo "$(GREEN)Checking release readiness...$(NC)"
	@echo "$(BLUE)1. Running tests...$(NC)"
	@$(MAKE) test
	@echo "$(BLUE)2. Running linter...$(NC)"
	@$(MAKE) lint
	@echo "$(BLUE)3. Building...$(NC)"
	@$(MAKE) build
	@echo "$(GREEN)Release checks passed$(NC)"

## Documentation
docs:
	@echo "$(GREEN)Documentation locations:$(NC)"
	@echo "  - Architecture: docs/ARCHITECTURE.md"
	@echo "  - Quick Start: docs/QUICKSTART.md"
	@echo "  - API Reference: docs/API_REFERENCE.md"
	@echo "  - Contributing: CONTRIBUTING.md"

godoc:
	@echo "$(GREEN)Starting godoc server...$(NC)"
	@if command -v godoc > /dev/null; then \
		echo "$(GREEN)Open http://localhost:6060$(NC)"; \
		godoc -http=:6060; \
	else \
		echo "$(YELLOW)godoc not installed. Install with: go install golang.org/x/tools/cmd/godoc@latest$(NC)"; \
	fi

## Docker
docker-build:
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t telemetryflow-go-sdk:latest .
	@echo "$(GREEN)Docker image built$(NC)"

docker-push: docker-build
	@echo "$(GREEN)Pushing Docker image...$(NC)"
	@docker push telemetryflow/telemetryflow-go-sdk:$(VERSION)
	@docker push telemetryflow/telemetryflow-go-sdk:latest

# ===========================================================================
# CI-Specific Targets
# ===========================================================================
# These targets are optimized for CI environments with race detection,
# coverage output, and atomic mode for accurate coverage in concurrent tests.

.PHONY: fmt-check deps-verify staticcheck test-unit-ci test-integration-ci test-e2e-ci security govulncheck coverage-merge coverage-report ci-build-gen ci-build-restapi ci-lint ci-test ci-build

## CI: Check formatting (fails if code needs formatting)
fmt-check:
	@echo "$(GREEN)Checking code formatting...$(NC)"
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "$(RED)The following files need formatting:$(NC)"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "$(GREEN)All files are properly formatted$(NC)"

## CI: Download and verify dependencies
deps-verify:
	@echo "$(GREEN)Downloading and verifying dependencies...$(NC)"
	@$(GOMOD) download
	@$(GOMOD) verify
	@echo "$(GREEN)Dependencies verified$(NC)"

## CI: Install and run staticcheck
staticcheck:
	@echo "$(GREEN)Running staticcheck...$(NC)"
	@if ! command -v staticcheck >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing staticcheck...$(NC)"; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
	fi
	@staticcheck ./...
	@echo "$(GREEN)Staticcheck complete$(NC)"

## CI: Run unit tests with race detection and coverage
test-unit-ci:
	@echo "$(GREEN)Running unit tests (CI mode)...$(NC)"
	@$(GOTEST) -v -race -timeout 5m -coverprofile=coverage-unit.out -covermode=atomic ./tests/unit/...
	@echo "$(GREEN)Unit tests complete$(NC)"

## CI: Run integration tests with race detection and coverage
test-integration-ci:
	@echo "$(GREEN)Running integration tests (CI mode)...$(NC)"
	@$(MAKE) build-generators
	@$(GOTEST) -v -race -timeout 10m -coverprofile=coverage-integration.out -covermode=atomic ./tests/integration/...
	@echo "$(GREEN)Integration tests complete$(NC)"

## CI: Run E2E tests
test-e2e-ci:
	@echo "$(GREEN)Running E2E tests (CI mode)...$(NC)"
	@$(MAKE) build-generators
	@$(GOTEST) -v -timeout 15m ./tests/e2e/...
	@echo "$(GREEN)E2E tests complete$(NC)"

## CI: Run security scan with gosec
security:
	@echo "$(GREEN)Running security scan...$(NC)"
	@if ! command -v gosec >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing gosec...$(NC)"; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
	fi
	@gosec -no-fail -fmt sarif -out gosec-results.sarif ./... || true
	@echo "$(GREEN)Security scan complete$(NC)"

## CI: Run govulncheck
govulncheck:
	@echo "$(GREEN)Running govulncheck...$(NC)"
	@if ! command -v govulncheck >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing govulncheck...$(NC)"; \
		go install golang.org/x/vuln/cmd/govulncheck@latest; \
	fi
	@govulncheck ./... || true
	@echo "$(GREEN)Vulnerability check complete$(NC)"

## CI: Merge coverage files
coverage-merge:
	@echo "$(GREEN)Merging coverage files...$(NC)"
	@if ! command -v gocovmerge >/dev/null 2>&1; then \
		echo "$(YELLOW)Installing gocovmerge...$(NC)"; \
		go install github.com/wadey/gocovmerge@latest; \
	fi
	@if [ -f coverage-integration.out ]; then \
		gocovmerge coverage-unit.out coverage-integration.out > coverage-merged.out; \
	else \
		cp coverage-unit.out coverage-merged.out; \
	fi
	@echo "$(GREEN)Coverage files merged$(NC)"

## CI: Generate coverage report
coverage-report: coverage-merge
	@echo "$(GREEN)Generating coverage report...$(NC)"
	@$(GOCMD) tool cover -func=coverage-merged.out | tee coverage-summary.txt
	@$(GOCMD) tool cover -html=coverage-merged.out -o coverage.html
	@echo "$(GREEN)Coverage report generated$(NC)"

## CI: Build telemetryflow-gen for specific platform
ci-build-gen:
	@echo "$(GREEN)Building $(GENERATOR_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(DIST_DIR)
	@OUTPUT="$(DIST_DIR)/$(GENERATOR_NAME)-$(GOOS)-$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then OUTPUT="$${OUTPUT}.exe"; fi; \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $${OUTPUT} $(GENERATOR_PATH)
	@echo "$(GREEN)Build complete: $(DIST_DIR)/$(GENERATOR_NAME)-$(GOOS)-$(GOARCH)$(NC)"

## CI: Build telemetryflow-restapi for specific platform
ci-build-restapi:
	@echo "$(GREEN)Building $(RESTAPI_GENERATOR_NAME) for $(GOOS)/$(GOARCH)...$(NC)"
	@mkdir -p $(DIST_DIR)
	@OUTPUT="$(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-$(GOOS)-$(GOARCH)"; \
	if [ "$(GOOS)" = "windows" ]; then OUTPUT="$${OUTPUT}.exe"; fi; \
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -ldflags "$(LDFLAGS)" -o $${OUTPUT} $(RESTAPI_GENERATOR_PATH)
	@echo "$(GREEN)Build complete: $(DIST_DIR)/$(RESTAPI_GENERATOR_NAME)-$(GOOS)-$(GOARCH)$(NC)"

## CI: Build both generators for specific platform
ci-build-generators: ci-build-gen ci-build-restapi
	@echo "$(GREEN)All generators built for $(GOOS)/$(GOARCH)$(NC)"

## CI: Run full lint suite
ci-lint: deps-verify fmt-check vet staticcheck
	@echo "$(GREEN)CI lint suite complete$(NC)"

## CI: Run all tests
ci-test: test-unit-ci test-integration-ci
	@echo "$(GREEN)CI test suite complete$(NC)"

## Version info
version:
	@echo "$(GREEN)$(PRODUCT_NAME)$(NC)"
	@echo "  Version:      $(VERSION)"
	@echo "  Git Commit:   $(GIT_COMMIT)"
	@echo "  Git Branch:   $(GIT_BRANCH)"
	@echo "  Build Time:   $(BUILD_TIME)"
	@echo "  Go Version:   $(GO_VERSION)"

.DEFAULT_GOAL := help
