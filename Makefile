.PHONY: help build test clean fmt vet lint mod-tidy mod-download install example-check

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all packages
	@echo "Building all packages..."
	@go build ./...

test: ## Run all tests
	@echo "Running tests..."
	@go test -v ./...

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	@go test -v -short ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v -cover ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-coverage-html: test-coverage ## Generate and open coverage report in browser
	@echo "Opening coverage report..."
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open coverage.html; \
	elif command -v open >/dev/null 2>&1; then \
		open coverage.html; \
	else \
		echo "Coverage report generated: coverage.html"; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@go clean ./...
	@rm -f coverage.out coverage.html
	@find . -name "*.test" -delete

fmt: ## Format all Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (if installed)
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

mod-tidy: ## Tidy go modules
	@echo "Tidying modules..."
	@go mod tidy

mod-download: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download

mod-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	@go mod verify

install: mod-download ## Install the package
	@echo "Installing package..."
	@go install ./...

example-check: ## Check if examples compile
	@echo "Checking examples..."
	@go build ./examples/...

check: fmt vet lint ## Run all checks (fmt, vet, lint)
	@echo "All checks passed!"

ci: mod-tidy build test vet ## Run CI checks (tidy, build, test, vet)
	@echo "CI checks completed!"

all: clean mod-tidy fmt vet build test ## Run everything: clean, tidy, fmt, vet, build, test
	@echo "All tasks completed!"
