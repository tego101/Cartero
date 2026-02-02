# Build stage
FROM golang:1.23-alpine AS builder

# Install git, ca-certificates, and gcc for building
RUN apk add --no-cache git ca-certificates gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -a -o cartero .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and sqlite for database
RUN apk --no-cache add ca-certificates sqlite

# Create a non-root user
RUN addgroup -g 1001 -S cartero && \
    adduser -u 1001 -S cartero -G cartero

# Set the working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/cartero .

# Copy the assets directory
COPY --from=builder /app/assets ./assets

# Copy the views directory for templates
COPY --from=builder /app/views ./views

# Copy the layout directory for templates
COPY --from=builder /app/layout ./layout

# Create data directory for persistent database
RUN mkdir -p /app/data && chown -R cartero:cartero /app

# Switch to non-root user
USER cartero

# Expose ports for SMTP and Web
EXPOSE 1025 10122

# Set environment variables with defaults
ENV SMTP_HOST=0.0.0.0
ENV SMTP_PORT=1025
ENV WEB_HOST=0.0.0.0
ENV WEB_PORT=10122

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:10122/inbox || exit 1

# Run the application
CMD ["./cartero"]