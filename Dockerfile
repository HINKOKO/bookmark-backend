FROM golang:1.18-alpine

WORDIR /app

COPY . .

RUN go build -o bookmark_backend ./cmd/api

CMD ["./bookmark_backend"]
