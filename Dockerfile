# ---------- Build stage ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build tools (optional but handy)
RUN apk add --no-cache ca-certificates tzdata

# Copy go.mod/go.sum first to leverage layer caching
COPY go.mod ./
RUN go mod download

COPY . .

# Build the binary
RUN go build -o dummy-endpoint .

# ---------- Runtime stage ----------
FROM alpine:3.22

# Create non-root user and group
RUN addgroup -S verumex && adduser -S verumex -G verumex

WORKDIR /app

# Copy app binary
COPY --from=builder /app/dummy-endpoint /app/dummy-endpoint

# Set file ownership
RUN chown -R verumex:verumex /app

# Switch to non-root user
USER verumex

# Expose HTTP and HTTPS ports
EXPOSE 8080 8443

# Run the binary
ENTRYPOINT ["/app/dummy-endpoint"]
