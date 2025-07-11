
# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o eoe cmd/eoe/main.go

# Final stage
FROM alpine:latest

# Copy the binary from builder
COPY --from=builder /app/eoe /usr/local/bin

# Command to run the executable
CMD ["eoe", "run"]
