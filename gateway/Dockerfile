FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org
RUN apk add --no-cache git
COPY . .
RUN go build -o /bin/gateway ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /bin/gateway /bin/gateway
COPY config /config
EXPOSE 8080
ENTRYPOINT ["/bin/gateway"]
