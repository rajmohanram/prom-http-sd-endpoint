# Build stage
FROM golang:1.21-alpine AS build

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Final stage
FROM alpine:3.18

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN adduser -D -H -h /app appuser
USER appuser

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/server /app/
# Create a directory for config files
RUN mkdir -p /app/config

# Copy sample configuration file
COPY --from=build /app/targets.yaml /app/config/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./server", "-config", "config/targets.yaml"]
