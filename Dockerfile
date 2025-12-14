# Stage 1: Frontend Build (Node)
FROM node:20-bookworm AS frontend-builder
WORKDIR /app
COPY client/package.json client/package-lock.json ./
RUN npm install
COPY client/ .
RUN npm run build

# Stage 2: Backend Build (Go)
FROM golang:1.25-bookworm AS builder
# Install build dependencies
RUN apt update && apt install git make -y
WORKDIR /build

# Copy Go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY server/ ./server/

# Copy built frontend from Stage 1 to where the Go binary expects it
COPY --from=frontend-builder /app/dist ./server/static/

# Build the Go binary (CGO_ENABLED=1 for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o krampus-server ./server

# Stage 3: Create minimal runtime image (Debian Slim)
FROM debian:bookworm-slim

# Install runtime dependencies (GLIBC compatibility and libsqlite3)
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates sqlite3 && \
    rm -rf /var/lib/apt/lists/*

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
RUN groupadd -r krampus && useradd -r -u 1000 -g krampus krampus && \
    chown -R krampus:krampus /app
USER krampus

# Run the binary
CMD ["./krampus-server"]
