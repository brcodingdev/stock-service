FROM golang:alpine AS builder

WORKDIR /go/src/github.com/brcodingdev/stock-service

COPY . .

RUN apk add --no-cache git
RUN go mod tidy
RUN go build -o stockservice ./cmd


FROM alpine:3.16

RUN apk update \
    && apk upgrade

WORKDIR /app
COPY .env /app/.env

ENV RABBIT_HOST=host.docker.internal
COPY --from=builder /go/src/github.com/brcodingdev/stock-service/stockservice .

CMD ["./stockservice"]

EXPOSE 9013

