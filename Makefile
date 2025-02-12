# Variables
BINARY_NAME=ecommerce-service
BUILD_DIR=bin

# Default target
all: build

# Build the Go binary
build:
	mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/${BINARY_NAME} main.go
	@echo "✅ Build complete!"

# Run the Go application
run: build
	./${BUILD_DIR}/${BINARY_NAME}

# Clean up generated files
clean:
	rm -rf ${BUILD_DIR}
	@echo "🧹 Cleaned up build files!"

# Install dependencies
deps:
	go mod tidy
	@echo "📦 Dependencies installed!"

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run || echo "⚠️  Linting issues found!"

.PHONY: all build clean run deps lint
