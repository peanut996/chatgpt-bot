FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -ldflags '-w -s' -o main .

CMD ["./main"]
