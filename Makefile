BINARY_NAME=calendar-server
MAIN_PACKAGE=cmd/calendar-server/main.go

BUILD_DIR=bin

GOLANGCI_LINT_VERSION=v2.5.0

.PHONY: build test lint run vet lint-install test-sh


build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run:
	@echo "Starting $(BINARY_NAME)..."
	go run $(MAIN_PACKAGE)

test:
	@echo "Running tests..."
	go test -race -cover ./... 2>&1 | grep -E "^(ok|FAIL)" | grep -v "coverage: 0.0%"

test-sh:
	chmod +x test_api.sh
	./test_api.sh

lint:
	@echo "Running linters..."
	golangci-lint run ./...

lint-install:
	@echo "Installing golangci-lint..."
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)

vet:
	@echo "Running go vet..."
	go vet ./...
