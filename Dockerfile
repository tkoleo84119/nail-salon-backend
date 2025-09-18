# Stage 1: Build stage
FROM golang:1.24.5-alpine AS builder

# Install necessary packages for building
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application with optimizations for production
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o main ./cmd/server/main.go

# Stage 2: Runtime stage
FROM gcr.io/distroless/static-debian11:nonroot

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy necessary configuration files
COPY --from=builder /app/internal/errors/errors.json ./internal/errors/

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Set environment for production
ENV GIN_MODE=release

# Use nonroot user for security
USER nonroot:nonroot

# Expose the port that the app runs on
EXPOSE 13200

# Run the binary
ENTRYPOINT ["/app/main"]