# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies - ca-certificates for HTTPS while fetching modules
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files before copying source code for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o packman \
    cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies, for external HTTPS calls and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/packman .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Run the application
ENTRYPOINT ["./packman"]
