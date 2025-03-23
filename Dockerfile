FROM golang:1.20-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api/main.go

# Use a minimal alpine image for the final stage
FROM alpine:3.17

# Set working directory
WORKDIR /app

# Install necessary packages
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main .

# Copy .env.example as .env for default configuration
COPY .env.example .env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]