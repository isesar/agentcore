# Use the official Golang image as base
FROM golang:1.25.6-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Create non-root user
RUN adduser -D -u 1000 appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the configuration files
COPY .env .

# Change ownership of the binary
RUN chown appuser:appuser main
RUN chown appuser:appuser .env

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]