FROM golang:1.22.1-alpine3.19 AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with additional security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o main ./cmd/api/main.go

# Use a minimal distroless image for the final stage
FROM gcr.io/distroless/static-debian12:nonroot

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env.example as .env for default configuration
COPY .env.example .env

# Expose the application port
EXPOSE 8080

# Set health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD ["/app/main", "health"] || exit 1

# Run the application
ENTRYPOINT ["/app/main"]