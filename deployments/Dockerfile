FROM golang:alpine

WORKDIR /go/src/

COPY . .

COPY ./env/dev .env

RUN export CGO_ENABLED=0

RUN go build -o main cmd/api/main.go

CMD ["./main"]