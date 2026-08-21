# Stage 1: Build binary
FROM golang:1.26-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk caching module
COPY go.mod ./
# COPY go.sum ./ (Uncomment jika go.sum sudah ada)

# Download dependencies
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi (binary diletakkan di root dengan nama 'main')
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Stage 2: Runtime image yang ringan
FROM alpine:latest

WORKDIR /root/

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi Golang
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]