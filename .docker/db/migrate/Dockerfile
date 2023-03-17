FROM golang:1.19.7-alpine

WORKDIR /app

RUN apk update && apk upgrade && \
    apk add build-base && \
    go install github.com/cosmtrek/air@latest

WORKDIR /app/cmd/db/migrate

CMD ["air"]
