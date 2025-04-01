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

# ✅ Copy the JSON data files too
COPY --from=builder /app/nava_items.json .
COPY --from=builder /app/preferences.json .
COPY --from=builder /app/products.json .

EXPOSE 8080
CMD ["./app"]
