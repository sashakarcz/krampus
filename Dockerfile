# Multi-stage build for minimal image size

# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY server/ ./server/

# Copy pre-built frontend static files (built outside container)
COPY server/static/ ./server/static/

# Build the Go binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o krampus-server ./server

# Stage 2: Create minimal runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/krampus-server .

# Copy example env file
COPY .env.example .

# Create database directory
RUN mkdir -p /app/database

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/ping || exit 1

# Run as non-root user
RUN adduser -D -u 1000 krampus && \
    chown -R krampus:krampus /app
USER krampus

# Run the binary
CMD ["./krampus-server"]
