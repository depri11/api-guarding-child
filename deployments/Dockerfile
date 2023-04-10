FROM golang:alpine AS builder
WORKDIR /go/src/
COPY . ./
RUN go build -o main -ldflags=-X=main.version=${VERSION} cmd/api/main.go 
COPY --from=builder /go/src/main ./
CMD ["main"]