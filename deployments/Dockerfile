FROM golang:alpine

WORKDIR /go/src/

COPY . .

RUN go build -o main cmd/api/main.go

CMD ["./main"]