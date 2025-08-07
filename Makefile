.PHONY: build build-all clean test test-verbose test-pretty test-dots clean help

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

# Test
test: ## Run tests
	go test ./...

test-verbose: ## Run tests with verbose output
	go test -v ./...

test-pretty: ## Run tests with pretty formatting (requires gotestsum)
	@command -v gotestsum >/dev/null 2>&1 || { echo >&2 "gotestsum is required but not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	gotestsum --format testname ./...

test-dots: ## Run tests with pretty formatting (requires gotestsum)
	@command -v gotestsum >/dev/null 2>&1 || { echo >&2 "gotestsum is required but not installed. Run: go install gotest.tools/gotestsum@latest"; exit 1; }
	gotestsum --format dots ./...

# Clean
clean: ## Clean build artifacts
	rm -f pyinit
	rm -rf dist/

# Help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
