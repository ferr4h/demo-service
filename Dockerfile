# Use multi-stage build to keep the final image small
FROM golang:1.21-alpine AS builder

# Install git and ca-certificates (needed for some dependencies)
RUN apk add --no-cache git ca-certificates

# Set working directory in the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Note: Create .env file from env-example.txt in the project root

# Change ownership to the appuser
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port (default is 8080, can be changed via environment variable)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]