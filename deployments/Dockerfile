FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o binary cmd/api/main.go 

CMD ["binary"]