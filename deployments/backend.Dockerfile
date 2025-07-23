# Stage 1: Build
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

RUN go build -o server ./cmd/server

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]
