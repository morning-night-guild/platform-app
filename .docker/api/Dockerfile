FROM golang:1.19.6-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest

WORKDIR /app/cmd/appapi

CMD ["air"]
