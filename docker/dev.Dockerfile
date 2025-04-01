FROM golang:1.22.1-alpine3.19 AS builder

# Install Air for hot-reloading
RUN go install github.com/cosmtrek/air@v1.49.0

FROM golang:1.22.1-alpine3.19

# Create non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Install development dependencies
RUN apk add --no-cache git curl tzdata

# Copy air binary from builder and ensure it's executable
COPY --from=builder /go/bin/air /usr/local/bin/air
RUN chmod +x /usr/local/bin/air

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Set permissions
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Set health check
HEALTHCHECK --interval=10s --timeout=5s --start-period=10s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1

# Use Air for hot reloading in development
CMD ["/usr/local/bin/air", "-c", ".air.toml"]