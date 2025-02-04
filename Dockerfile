ARG GOLANG_VERSION

FROM golang:${GOLANG_VERSION} AS builder

WORKDIR /app

# COPY go.mod go.sum ./

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" \
    -o helm-aggregator main.go

FROM alpine:3.20

WORKDIR /app

RUN apk --no-cache add ca-certificates \
    && update-ca-certificates \
    && addgroup -g 9999 app \
    && adduser -s /dev/false -u 9999 -D -G app app \
    && chown -R app:app /app \
    && rm -rf /var/cache/apk/*

COPY --from=builder /app/helm-aggregator .

EXPOSE 8080

USER app

CMD ["./helm-aggregator", "run"]
