FROM golang:bluster

WORKDIR /go/src/

COPY . .

COPY ./env/dev .env

RUN CGO_ENABLED=1 go build -o main cmd/api/main.go

CMD ["./main"]