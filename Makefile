# Makefile for lingo

GO_PACKAGE ?= "lingo"

.PHONY: help test coverage lint fmt dev-deps

help:
	@echo "Usage: make <target>"
	@echo "Targets:"
	@echo "  help        Show this help message"
	@echo "  dev-deps    Install development dependencies"
	@echo "  test-unit   Run unit tests and generate coverage report"
	@echo "  lint        Run golangci-lint"
	@echo "  fmt         Tries to automatically fix linting errors"

# Install development dependencies
dev-deps:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.4.0

# Test the code and generate coverage report
test-unit:
	mkdir -p reporting
	go test -p=1 -short -cover -coverpkg=$$($(GO_PACKAGE) | tr '\n' ',') -coverprofile=reporting/profile_raw.out -json $$($(GO_PACKAGE)) > reporting/tests.json || true
	grep -v '_mock.go' reporting/profile_raw.out > reporting/profile.out
	go tool cover -html=reporting/profile.out -o reporting/coverage.html
	go tool cover -func=reporting/profile.out -o reporting/coverage.txt
	cat reporting/coverage.txt

# Run golangci-lint
lint:
	golangci-lint run

# Run with fix to automatically fix issues
fmt:
	golangci-lint run --fix