# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies for SQLite
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for the go-sqlite3 driver
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

# Final stage
FROM alpine:3.18

WORKDIR /app

# Add build arguments for user and group IDs
ARG UID=1000
ARG GID=1000

# Create a non-root user with specified UID and GID
RUN addgroup -g ${GID} -S appgroup && adduser -u ${UID} -G appgroup -s /bin/sh -D appuser

# Install runtime dependencies for SQLite
RUN apk add --no-cache curl sqlite-libs

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy migrations
COPY database/migrations ./database/migrations

# Change ownership of the app directory
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose the port the application runs on
EXPOSE 8080