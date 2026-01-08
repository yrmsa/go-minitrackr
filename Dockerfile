# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o go-minitrackr \
    ./cmd/server

# Final stage
FROM scratch

# Copy binary
COPY --from=builder /build/go-minitrackr /go-minitrackr

# Copy CA certificates for HTTPS (if needed)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create data directory
VOLUME ["/data"]

# Set environment
ENV PORT=8822
ENV DB_PATH=/data/go-minitrackr.db
ENV GOMEMLIMIT=25MiB

# Expose port
EXPOSE 8822

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/go-minitrackr", "health"] || exit 1

# Run
ENTRYPOINT ["/go-minitrackr"]
