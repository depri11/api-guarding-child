FROM golang:alpine

WORKDIR /go/src/

COPY . .

COPY ./env/dev .env

RUN export CGO_ENABLED=1

RUN go build -o main cmd/api/main.go

CMD ["./main"]