# Frontend Build Stage
FROM node:20-alpine AS webbuilder

WORKDIR /app/web

COPY web/package.json ./
RUN npm install

COPY web/ .
RUN npm run build

# Build Stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o obsidian ./cmd/obsidian

# Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates sqlite

# Copy binary from builder
COPY --from=builder /app/obsidian .
COPY --from=webbuilder /app/web/dist ./web/dist

# Create data directory
RUN mkdir -p data

# Environment variables
ENV HEADLESS=true
ENV WEB_STATIC_DIR=/app/web/dist

# Volume for persistence
VOLUME ["/app/data"]

# Run
CMD ["./obsidian"]
