.PHONY: help install dev build clean test deps air-install css-build css-watch templ-generate run setup

# Variables
APP_NAME=parallax
BINARY_DIR=./bin
BINARY_PATH=$(BINARY_DIR)/$(APP_NAME)
TMP_DIR=./tmp

# Default target
help:
	@echo "Available commands:"
	@echo "  help          Show this help message"
	@echo "  install       Install all dependencies (Go, Tailwind CLI, Air)"
	@echo "  dev           Start development with Air (live reload)"
	@echo "  build         Build the application for production"
	@echo "  run           Run the application"
	@echo "  clean         Clean build artifacts"
	@echo "  test          Run tests"
	@echo "  deps          Download Go dependencies"
	@echo "  air-install   Install Air for live reload"
	@echo "  css-build     Build CSS with Tailwind"
	@echo "  css-watch     Watch and rebuild CSS"
	@echo "  templ-generate Generate Go code from Templ templates"
	@echo "  setup         Complete development setup (first time)"

# Install all dependencies
install: deps air-install tailwind-install
	@echo "All dependencies installed!"

# Install Air
air-install:
	@echo "Installing Air..."
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@echo "Air installed successfully!"

# Install Tailwind CLI
tailwind-install:
	@echo "Installing Tailwind CLI..."
	pnpm install tailwindcss @tailwindcss/cli
	@which tailwindcss > /dev/null || pnpm install -g tailwindcss @tailwindcss/cli
	@echo "Tailwind CLI installed successfully!"

# Download Go dependencies
deps:
	@echo "Downloading Go dependencies..."
	go mod download
	go mod tidy

# Generate Templ templates
templ-generate:
	@echo "Generating Templ templates..."
	@which templ > /dev/null || go install github.com/a-h/templ/cmd/templ@latest
	templ generate


# Build the application
build: clean
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BINARY_DIR)
	go build -ldflags="-s -w" -o $(BINARY_PATH) ./cmd/parallax
	@echo "Built $(BINARY_PATH)"

# Run the application
run: templ-generate css-build
	@echo "Running $(APP_NAME)..."
	go run ./cmd/parallax

# Development with Air (live reload)
dev:
	@echo "Starting development server with Air..."
	@mkdir -p $(TMP_DIR)
	air

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BINARY_DIR)
	@rm -rf $(TMP_DIR)
	@rm -rf web/static/dist
	@rm -f web/static/tailwind.css
	@rm -f build-errors.log
	@echo "Clean complete!"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Development setup (first time)
setup: install templ-generate css-build
	@echo "Development environment setup complete!"
	@echo "Run 'make dev' to start development with live reload"
