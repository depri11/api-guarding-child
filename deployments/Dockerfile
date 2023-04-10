FROM golang:alpine

WORKDIR /go/src/

COPY . .

COPY .env .

RUN go build -o main cmd/api/main.go

CMD ["./main"]