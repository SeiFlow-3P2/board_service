FROM golang:1.24-bullseye AS builder
RUN apt-get update && apt-get install -y gcc libc6-dev ca-certificates
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app/main.go

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y ca-certificates tzdata curl && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env.example .env

CMD ["./main"]