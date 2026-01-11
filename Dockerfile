# Multi-stage build for gobuilder-mcp
FROM golang:1.21-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gobuilder-mcp main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S gobuilder && \
    adduser -u 1001 -S gobuilder -G gobuilder

# Set the working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/gobuilder-mcp .

# Change ownership to non-root user
RUN chown gobuilder:gobuilder /app/gobuilder-mcp

# Switch to non-root user
USER gobuilder

# Expose port (if needed for future web interface)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ./gobuilder-mcp || exit 1

# Set the entry point
ENTRYPOINT ["./gobuilder-mcp"]
