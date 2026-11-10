YOUR_NAME = mrb4el
BINARY_NAME=terraform-provider-aeza
VERSION=0.1.0
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

TEST_DIR = test-terraform
PLUGIN_PATH = $(TEST_DIR)/.terraform/plugins/registry.terraform.io/$(YOUR_NAME)/aeza/$(VERSION)/$(GOOS)_$(GOARCH)

.PHONY: build clean test install

print-vars:
	@echo "BINARY_NAME: $(BINARY_NAME)"
	@echo "VERSION: $(VERSION)"
	@echo "GOOS: $(GOOS)"
	@echo "GOARCH: $(GOARCH)"
	@echo "YOUR_NAME: $(YOUR_NAME)"

build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) .

build-all-platforms:
	@echo "Building for all platforms..."
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/amd64/$(BINARY_NAME) .
	@GOOS=darwin GOARCH=amd64 go build -o bin/darwin/amd64/$(BINARY_NAME) .
	@GOOS=windows GOARCH=amd64 go build -o bin/windows/amd64/$(BINARY_NAME).exe .

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf bin/

test:
	@echo "Running tests..."
	@go test -v ./...

test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

install: build
	@echo "Installing provider to test directory..."
	@mkdir -p $(PLUGIN_PATH)
	@cp $(BINARY_NAME) $(PLUGIN_PATH)/$(BINARY_NAME)_v$(VERSION)
	@echo "Provider installed to: $(PLUGIN_PATH)"

init: 
	@rm -f $(TEST_DIR)/.terraform.lock.hcl
	@cd $(TEST_DIR) && terraform init

fmt:
	@echo "Formatting code..."
	@gofmt -w -s .

lint:
	@echo "Linting code..."
	@golangci-lint run

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build the provider"
	@echo "  build-all  - Build for all platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  install    - Install provider locally"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"