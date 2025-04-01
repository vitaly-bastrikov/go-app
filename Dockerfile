# Stage 1: Build
FROM golang:1.23 AS builder

WORKDIR /app

# Copy Go module files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build the app
RUN go build -o app .

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy the binary
COPY --from=builder /app/app .

# Copy the data files
COPY --from=builder /app/data/nava_items.json ./data/
COPY --from=builder /app/data/preferences.json ./data/
COPY --from=builder /app/data/products.json ./data/

# Create data directory
RUN mkdir -p data

EXPOSE 8080
CMD ["./app"]
