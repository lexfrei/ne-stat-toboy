############################
# STEP 1: Build executable
############################
FROM golang:1.24-alpine AS builder

# Install git, CA certificates, and direct dependencies in one layer
RUN apk add --no-cache git ca-certificates tzdata && \
    # Create appuser in the same layer
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "appuser"

# Set working directory
WORKDIR /app

# Copy go mod files first for caching dependencies
COPY go.mod go.sum* ./

# Download dependencies and install tools in one layer
RUN go mod download && \
    go install github.com/a-h/templ/cmd/templ@v0.3.857 && \
    go install github.com/tdewolff/minify/v2/cmd/minify@latest

# Copy source code
COPY . .

# Generate templates
RUN templ generate

# Minify static files directly with the minify tool
RUN minify -r -o ./web/static/ ./web/static/

# Build the executable with optimizations, skip UPX to reduce build time
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /app/ne-stat-toboy ./cmd/server


############################
# STEP 2: Build final image
############################
FROM alpine:3.21 AS final

# Install wget for health check in a single layer
RUN apk --no-cache add wget ca-certificates tzdata

# Create appuser directory structure
RUN mkdir -p /web/static && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "appuser" && \
    chown -R appuser:appuser /web

# Copy the executable and minified static files
COPY --from=builder --chown=appuser:appuser /app/ne-stat-toboy /ne-stat-toboy
COPY --from=builder --chown=appuser:appuser /app/web/static /web/static

# Use appuser
USER appuser:appuser

# Expose port
EXPOSE 8080

# Add health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/healthz || exit 1

# Run the binary
ENTRYPOINT ["/ne-stat-toboy"]