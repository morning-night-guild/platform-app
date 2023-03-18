FROM golang:1.19.7-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app/cmd/app/core

CMD ["air"]
