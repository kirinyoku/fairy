.PHONY: all tidy fmt lint test build clean

# Default target
all: tidy fmt lint test

# Download and tidy up dependencies
tidy:
	@echo "==> Tidying up module dependencies..."
	go mod tidy

# Format the code
fmt:
	@echo "==> Formatting code..."
	go fmt ./...

# Run the linter
lint:
	@echo "==> Running golangci-lint..."
	golangci-lint run ./...

# Run tests
test:
	@echo "==> Running tests..."
	go test -v -race ./...

# Build the project (checks if it compiles)
build:
	@echo "==> Building packages..."
	go build ./...

# Clean build cache
clean:
	@echo "==> Cleaning build cache..."
	go clean -testcache
	go clean -modcache
