# Build stage
FROM golang:1.22-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations for size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/ingestion \
    .

# Final stage - minimal image
FROM scratch

# Copy CA certificates for HTTPS connections
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy the binary
COPY --from=builder /app/ingestion /ingestion

# Run as non-root user (using numeric UID)
USER 65534

# Expose no ports (this is a background service)

# Set the entrypoint
ENTRYPOINT ["/ingestion"]
