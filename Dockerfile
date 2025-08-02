FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rl-textlib-learner ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/rl-textlib-learner .

# Create directories for logs and data
RUN mkdir -p /app/logs /app/data /app/models && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Resource limits
ENV GOMAXPROCS=2
ENV GOMEMLIMIT=512MiB

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
    CMD /app/rl-textlib-learner --health-check || exit 1

# Run the application
CMD ["./rl-textlib-learner"]