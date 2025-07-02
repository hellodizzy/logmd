# Makefile for logmd
# Learn: Makefiles provide consistent build automation across different environments.
# See: https://www.gnu.org/software/make/manual/make.html

# Build configuration
BINARY_NAME=logmd
DIST_DIR=dist
MAIN_PACKAGE=.
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go configuration
GOFLAGS=-mod=readonly
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Default target
.PHONY: all
all: test lint build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(DIST_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "Built $(DIST_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
.PHONY: build-all
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PACKAGE)
	@echo "Built binaries for all platforms in $(DIST_DIR)/"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test $(GOFLAGS) -v -race -cover ./...

# Run tests with coverage report
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test $(GOFLAGS) -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	go vet ./...
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not found, installing..."; \
		go install honnef.co/go/tools/cmd/staticcheck@latest; \
		staticcheck ./...; \
	fi
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, installing..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
		golangci-lint run; \
	fi

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Tidy dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Generate Homebrew formula
.PHONY: brew
brew: build-all
	@echo "Generating Homebrew formula..."
	@DARWIN_AMD64_SHA=$$(shasum -a 256 $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 | cut -d' ' -f1); \
	DARWIN_ARM64_SHA=$$(shasum -a 256 $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 | cut -d' ' -f1); \
	LINUX_AMD64_SHA=$$(shasum -a 256 $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 | cut -d' ' -f1); \
	LINUX_ARM64_SHA=$$(shasum -a 256 $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 | cut -d' ' -f1); \
	sed -e "s/{{VERSION}}/$(VERSION)/g" \
		-e "s/{{DARWIN_AMD64_SHA}}/$$DARWIN_AMD64_SHA/g" \
		-e "s/{{DARWIN_ARM64_SHA}}/$$DARWIN_ARM64_SHA/g" \
		-e "s/{{LINUX_AMD64_SHA}}/$$LINUX_AMD64_SHA/g" \
		-e "s/{{LINUX_ARM64_SHA}}/$$LINUX_ARM64_SHA/g" \
		logmd.rb.template > logmd.rb || echo "Formula template not found, creating basic formula..."; \
	if [ ! -f logmd.rb ]; then \
		echo 'class Logmd < Formula' > logmd.rb; \
		echo '  desc "A minimal, local-first journal CLI"' >> logmd.rb; \
		echo '  homepage "https://github.com/your-username/logmd"' >> logmd.rb; \
		echo '  version "$(VERSION)"' >> logmd.rb; \
		echo '' >> logmd.rb; \
		echo '  if OS.mac? && Hardware::CPU.intel?' >> logmd.rb; \
		echo '    url "https://github.com/your-username/logmd/releases/download/$(VERSION)/$(BINARY_NAME)-darwin-amd64"' >> logmd.rb; \
		echo "    sha256 \"$$DARWIN_AMD64_SHA\"" >> logmd.rb; \
		echo '  elsif OS.mac? && Hardware::CPU.arm?' >> logmd.rb; \
		echo '    url "https://github.com/your-username/logmd/releases/download/$(VERSION)/$(BINARY_NAME)-darwin-arm64"' >> logmd.rb; \
		echo "    sha256 \"$$DARWIN_ARM64_SHA\"" >> logmd.rb; \
		echo '  elsif OS.linux? && Hardware::CPU.intel?' >> logmd.rb; \
		echo '    url "https://github.com/your-username/logmd/releases/download/$(VERSION)/$(BINARY_NAME)-linux-amd64"' >> logmd.rb; \
		echo "    sha256 \"$$LINUX_AMD64_SHA\"" >> logmd.rb; \
		echo '  elsif OS.linux? && Hardware::CPU.arm?' >> logmd.rb; \
		echo '    url "https://github.com/your-username/logmd/releases/download/$(VERSION)/$(BINARY_NAME)-linux-arm64"' >> logmd.rb; \
		echo "    sha256 \"$$LINUX_ARM64_SHA\"" >> logmd.rb; \
		echo '  end' >> logmd.rb; \
		echo '' >> logmd.rb; \
		echo '  def install' >> logmd.rb; \
		echo '    bin.install "$(BINARY_NAME)"' >> logmd.rb; \
		echo '  end' >> logmd.rb; \
		echo '' >> logmd.rb; \
		echo '  test do' >> logmd.rb; \
		echo '    system "#{bin}/$(BINARY_NAME)", "--version"' >> logmd.rb; \
		echo '  end' >> logmd.rb; \
		echo 'end' >> logmd.rb; \
	fi
	@echo "Homebrew formula generated: logmd.rb"

# Install locally
.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(DIST_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "Installed $(BINARY_NAME)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html
	rm -f logmd.rb
	@echo "Clean complete"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary for current platform"
	@echo "  build-all    - Build binaries for all platforms"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  lint         - Run linters (go vet, staticcheck, golangci-lint)"
	@echo "  fmt          - Format code"
	@echo "  tidy         - Tidy dependencies"
	@echo "  brew         - Generate Homebrew formula"
	@echo "  install      - Install binary locally"
	@echo "  clean        - Clean build artifacts"
	@echo "  help         - Show this help" 