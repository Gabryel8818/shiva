FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o main main.go

CMD ["./main"]