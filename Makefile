.PHONY: build build-all clean test test-verbose test-cover test-cover-func test-integration test-unit test-ci test-pretty test-dots test-watch fmt fmt-check vet lint lint-fix check help

# Get version from git tag, fallback to 'dev'
VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || echo "dev")
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')

# Go build flags
LDFLAGS := -ldflags "-X 'github.com/Pradyothsp/pyinit/internal/version.Version=$(VERSION)' \
	-X 'github.com/Pradyothsp/pyinit/internal/version.GitCommit=$(GIT_COMMIT)' \
	-X 'github.com/Pradyothsp/pyinit/internal/version.BuildDate=$(BUILD_DATE)'"

# Default build target
build: ## Build the binary
	go build $(LDFLAGS) -o pyinit ./cmd/pyinit

# Build for all platforms
build-all: ## Build for all platforms
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/pyinit-linux-amd64 ./cmd/pyinit
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/pyinit-linux-arm64 ./cmd/pyinit
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/pyinit-darwin-amd64 ./cmd/pyinit
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/pyinit-darwin-arm64 ./cmd/pyinit
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/pyinit-windows-amd64.exe ./cmd/pyinit

# Create dist directory
dist:
	mkdir -p dist

# Test targets
test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-cover: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-cover-func: ## Show test coverage by function
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

test-integration: ## Run integration tests only
	go test -v ./test/integration/...

test-unit: ## Run unit tests only (exclude integration)
	go test -v ./internal/... ./pkg/...

test-ci: ## Run tests suitable for CI (with race detection and coverage)
	go test -v -race -timeout=10m -coverprofile=coverage.out -covermode=atomic ./...

test-pretty: ## Run tests with pretty formatting (requires gotestsum)
	@command -v gotestsum >/dev/null 2>&1 || { echo >&2 "gotestsum is required but not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	gotestsum --format testname ./...

test-dots: ## Run tests with pretty formatting (requires gotestsum)
	@command -v gotestsum >/dev/null 2>&1 || { echo >&2 "gotestsum is required but not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	gotestsum --format dots ./...

test-watch: ## Watch for changes and run tests (requires gotestsum)
	@command -v gotestsum >/dev/null 2>&1 || { echo >&2 "gotestsum is required but not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	gotestsum --watch ./...

# Quality assurance targets
fmt: ## Format code
	go fmt ./...
	goimports -w .

fmt-check: ## Check if code is formatted
	@if [ -n "$$(gofmt -s -l .)" ]; then \
		echo "The following files are not formatted:"; \
		gofmt -s -l .; \
		echo "Please run 'make fmt' to format your code."; \
		exit 1; \
	fi

vet: ## Run go vet
	go vet ./...

lint: ## Run golangci-lint
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is required but not installed. See: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run ./...

lint-fix: ## Run golangci-lint with --fix
	@command -v golangci-lint >/dev/null 2>&1 || { echo >&2 "golangci-lint is required but not installed. See: https://golangci-lint.run/usage/install/"; exit 1; }
	golangci-lint run --fix ./...

check: fmt-check vet lint test-ci ## Run all quality checks (format, vet, lint, tests)

# Clean
clean: ## Clean build artifacts
	rm -f pyinit
	rm -rf dist/
	rm -f coverage.out coverage.html

# Help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
