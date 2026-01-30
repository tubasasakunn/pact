.PHONY: test test-unit test-integration test-e2e test-api test-coverage test-pkg test-run test-watch lint bench build clean

# Default target
all: lint test build

# Build
build:
	go build -o bin/pact ./cmd/pact

# Clean
clean:
	rm -rf bin/ coverage.out coverage.html

# All tests
test: lint test-unit test-integration test-e2e

# Unit tests
test-unit:
	go test -v -race ./internal/...

# Integration tests
test-integration:
	go test -v -race ./test/integration/...

# E2E tests
test-e2e:
	go test -v -race ./test/e2e/...

# Public API tests
test-api:
	go test -v -race ./pkg/...

# Coverage
test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Specific package test
test-pkg:
	go test -v -race ./$(PKG)/...

# Run specific test
test-run:
	go test -v -race -run $(RUN) ./...

# Watch mode (requires watchexec)
test-watch:
	watchexec -e go "make test-unit"

# Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

# Benchmark
bench:
	go test -bench=. -benchmem ./internal/infrastructure/parser/...
	go test -bench=. -benchmem ./internal/application/transformer/...

# Format
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Generate (if needed)
generate:
	go generate ./...

# Install development tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/watchexec/watchexec@latest

# Help
help:
	@echo "Available targets:"
	@echo "  build          - Build the CLI binary"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run all tests (unit, integration, e2e)"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-e2e       - Run E2E tests"
	@echo "  test-api       - Run public API tests"
	@echo "  test-coverage  - Generate coverage report"
	@echo "  test-pkg PKG=x - Test specific package"
	@echo "  test-run RUN=x - Run specific test"
	@echo "  test-watch     - Watch mode for unit tests"
	@echo "  lint           - Run linter"
	@echo "  bench          - Run benchmarks"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy go.mod"
	@echo "  tools          - Install development tools"
