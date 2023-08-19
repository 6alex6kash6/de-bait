FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go build ./cmd/main.go

FROM alpine:latest

COPY --from=builder /app /app

WORKDIR /app

CMD ["./main"]