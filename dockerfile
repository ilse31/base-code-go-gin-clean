# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./


# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install SSL certificates for HTTPS and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy environment file (if exists)
COPY .env* ./

# Create target directory for email templates
RUN mkdir -p /app/internal/email/templates

# Copy email templates
COPY internal/email/templates/ /app/internal/email/templates/

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
CMD ["/app/main"]