FROM golang:bluster

WORKDIR /go/src/

COPY . .

RUN go build -o main cmd/api/main.go

CMD ["./main"]