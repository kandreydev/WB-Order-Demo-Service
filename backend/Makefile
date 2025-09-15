# Default target: build the application binary
.PHONY: all build clean run test

all: build

# Compile the Go application
build:
	@go build -o app ./...

# Remove the compiled binary
clean:
	@rm -f app

# Build and run the application
run: build
	@./app

# Run tests
test:
	@go test ./...

# Fast run
frun:
	@go run ./cmd/main.go