FROM golang:alpine

RUN apk add build-base

WORKDIR /go/src/

COPY . .

COPY ./env/dev .env

RUN go build -o main cmd/api/main.go

CMD ["./main"]