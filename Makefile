.PHONY: run build test clean

# Run the SSH server
run:
	@echo "Starting promptgo SSH server..."
	@go run cmd/server/main.go

# Build the binary
build:
	@echo "Building promptgo..."
	@go build -o bin/promptgo cmd/server/main.go
	@echo "Binary created at bin/promptgo"

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@echo "Clean complete"
