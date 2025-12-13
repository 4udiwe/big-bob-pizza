# Step 1: Modules caching
FROM golang:1.24-alpine3.21 AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Step 2: Tests
FROM golang:1.24-alpine3.21
RUN apk add --no-cache wget

WORKDIR /app

COPY --from=modules /go/pkg /go/pkg
COPY go.mod go.sum ./
COPY internal ./internal
COPY pkg ./pkg
COPY cmd ./cmd
COPY test ./test

CMD ["go", "test", "-v", "-tags=integration", "./test/integration_test/..."]
