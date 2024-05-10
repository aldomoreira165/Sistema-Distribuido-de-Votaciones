FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o client .

EXPOSE 3000

CMD ["./client"]