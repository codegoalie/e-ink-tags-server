FROM golang:1.21.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o etag-server ./...

FROM alpine:latest

WORKDIR /

COPY --from=builder /app/etag-server .

EXPOSE 8080

CMD ["./etag-server"]
