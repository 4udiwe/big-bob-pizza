# Step 1: Modules caching
FROM golang:1.24-alpine3.22 AS modules
COPY analytics-service/go.mod analytics-service/go.sum /modules/analytics-service/
COPY order-service/go.mod order-service/go.sum /modules/order-service/
WORKDIR /modules/analytics-service
RUN go mod download

# Step 2: Builder
FROM golang:1.24-alpine3.22 AS builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app/analytics-service

# Сборка бинарника в отдельную директорию
RUN mkdir -p /app/bin
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/analytics-service ./cmd/main.go

# Step 3: Final
FROM alpine:3.22
RUN apk --no-cache add ca-certificates wget

COPY --from=builder /app/analytics-service/config /config
COPY --from=builder /app/bin/analytics-service /app/analytics-service
COPY --from=builder /app/analytics-service/internal/database/migrations /app/database/migrations

WORKDIR /app
EXPOSE 8083

CMD ["/app/analytics-service"]

