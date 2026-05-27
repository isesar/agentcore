.PHONY: run test lint clean migrate

# Default target
all: run

# Run the application
run:
	go run main.go

# Run database migrations
migrate:
	go run ./cmd/migrate/main.go

# Run tests
test:
	go test ./... -v

# Run linter
lint:
	golangci-lint run --fix

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
install-deps:
	go mod download

# Build the application
build:
	go build -o bin/agentcore main.go

# Run with coverage
test-coverage:
	go test ./... -coverprofile=coverage.txt && go tool cover -html=coverage.txt

# Generate documentation
docs:
	godoc -http=:6060