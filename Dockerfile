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

# Final stage - use alpine for debugging
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Copy binary
COPY --from=builder /build/go-minitrackr /go-minitrackr

# Create data directory
RUN mkdir -p /data
VOLUME ["/data"]

# Set environment
ENV PORT=8822
ENV DB_PATH=/data/go-minitrackr.db
ENV GOMEMLIMIT=25MiB

# Expose port
EXPOSE 8822

# Run
ENTRYPOINT ["/go-minitrackr"]
