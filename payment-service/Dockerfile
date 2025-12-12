# Step 1: Modules caching
FROM golang:1.24-alpine3.22 AS modules
COPY payment-service/go.mod payment-service/go.sum /modules/payment-service/
COPY order-service/go.mod order-service/go.sum /modules/order-service/
WORKDIR /modules/payment-service
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine3.22 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app/payment-service

# Сборка бинарника в отдельную директорию
RUN mkdir -p /app/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/payment-service ./cmd/main.go

# Step 3: Final
FROM alpine:3.22
RUN apk --no-cache add ca-certificates wget

COPY --from=builder /app/payment-service/config /config
COPY --from=builder /app/bin/payment-service /app/payment-service
COPY --from=builder /app/payment-service/internal/database/migrations /app/database/migrations

WORKDIR /app
EXPOSE 8081

CMD ["/app/payment-service"]
