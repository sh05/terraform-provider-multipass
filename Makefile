default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build the provider
.PHONY: build
build:
	go build -o opentofu-provider-multipass

# Install the provider locally for development
.PHONY: install-local
install-local: build
	@OS_ARCH=$$(go env GOOS)_$$(go env GOARCH) && \
	mkdir -p ~/.opentofu/plugins/registry.opentofu.org/sh05/multipass/0.1.0/$$OS_ARCH && \
	cp opentofu-provider-multipass ~/.opentofu/plugins/registry.opentofu.org/sh05/multipass/0.1.0/$$OS_ARCH/ && \
	echo "Provider installed for $$OS_ARCH"

# Run unit tests
.PHONY: test
test:
	go test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	go test -v -cover ./...

# Generate documentation
.PHONY: docs
docs:
	go generate

# Format Go code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint the code
.PHONY: lint
lint:
	golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	rm -f opentofu-provider-multipass
	go clean

# Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

# Initialize the provider for local development
.PHONY: init
init: deps fmt

# Full development cycle: clean, build, test, lint
.PHONY: dev
dev: clean deps fmt test lint build

# Release build
.PHONY: release
release:
	goreleaser release --rm-dist

# Install development tools
.PHONY: tools
tools:
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the provider binary"
	@echo "  install-local  - Install provider locally for development"
	@echo "  test           - Run unit tests"
	@echo "  testacc        - Run acceptance tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  docs           - Generate documentation"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Run linter"
	@echo "  clean          - Clean build artifacts"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  init           - Initialize for development"
	@echo "  dev            - Full development cycle"
	@echo "  tools          - Install development tools"
	@echo "  help           - Show this help message"