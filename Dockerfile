FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy the Go module files to the working directory
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the rest of the application source code to the working directory
COPY . .

# Build the Go application
RUN go build -o tax-calculator ./main

# Create the service image
FROM alpine:latest as service

WORKDIR /app

# Copy the built Go binary from the previous stage
COPY --from=builder /app/tax-calculator .
COPY --from=builder /app/config.yml .

# Expose the port the server will be listening on
EXPOSE 8080
CMD ["./tax-calculator"]

# Build the integration test image
FROM builder as integration-test
WORKDIR /app
CMD ["go", "test", "./integration", "-v", "-tags=integration"]