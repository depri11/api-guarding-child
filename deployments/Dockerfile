FROM golang:alpine AS builder

WORKDIR /go/src/

COPY . ./

RUN go build -o /tmp/app cmd/api/main.go

WORKDIR /

COPY --from=builder /tmp/app ./

CMD ["./app"]
