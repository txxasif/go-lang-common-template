# Development stage
FROM golang:1.23.0-alpine

# Install build dependencies
RUN apk add --no-cache git curl

# Install air for hot reloading
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Run air for hot reloading
CMD ["air"]