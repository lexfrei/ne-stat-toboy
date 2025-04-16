############################
# STEP 1: Build executable
############################
FROM golang:1.24-alpine AS builder

# Install git and CA certificates for fetching dependencies
RUN apk update && apk add --no-cache git ca-certificates upx && update-ca-certificates

# Create appuser
ENV USER=appuser
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Set working directory
WORKDIR /app

# Copy go mod files first for caching dependencies
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@v0.3.857

# Copy source code
COPY . .

# Generate templates
RUN templ generate

# Minify static files
RUN go run ./cmd/minify/main.go ./web/static

# Build the executable with optimizations
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o /app/ne-stat-toboy ./cmd/server && \
    upx --best --lzma /app/ne-stat-toboy


############################
# STEP 2: Build final image
############################
FROM scratch

# Import CA certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import user and group files
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the executable
COPY --from=builder /app/ne-stat-toboy /ne-stat-toboy

# Copy minified static files
COPY --from=builder /app/web/static /web/static

# Use appuser
USER appuser:appuser

# Expose port
EXPOSE 8080

# Run the binary
ENTRYPOINT ["/ne-stat-toboy"]