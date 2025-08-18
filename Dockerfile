# Stage 1: Build
FROM golang:1.24.3 AS builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /myapp ./cmd/main.go

# Stage 2: Runtime
FROM alpine:latest
WORKDIR /app
COPY --from=builder /myapp /app/
RUN apk add --no-cache libc6-compat && \
    chmod +x /app/myapp
CMD ["/app/myapp"]