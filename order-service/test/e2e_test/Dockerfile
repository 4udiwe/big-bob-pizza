# Step 1: modules cache
FROM golang:1.24-alpine3.21 AS modules
WORKDIR /modules
COPY go.mod go.sum ./
RUN go mod download

# Step 2: e2e tests
FROM golang:1.24-alpine3.21
RUN apk add --no-cache wget

WORKDIR /app

COPY --from=modules /go/pkg /go/pkg
COPY go.mod go.sum ./
COPY internal ./internal
COPY pkg ./pkg
COPY cmd ./cmd
COPY test ./test

CMD ["go", "test", "-v", "-tags=e2e", "./test/e2e_test/..."]
