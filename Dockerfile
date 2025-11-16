# =========================
# Build stage
# =========================
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary (static binary for Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o civix-backend ./cmd/api

# =========================
# Final stage
# =========================
FROM alpine:3.18

WORKDIR /app

# Install runtime dependencies (e.g., timezones)
RUN apk add --no-cache tzdata ca-certificates && update-ca-certificates

# Copy the compiled binary
COPY --from=builder /app/civix-backend .

# Optionally copy static assets if needed
# COPY --from=builder /app/static ./static
# COPY --from=builder /app/templates ./templates

# Expose the app port
EXPOSE 8080

# Use a non-root user (recommended for production)
RUN adduser -D appuser
USER appuser

# Run the app
CMD ["./civix-backend"]
