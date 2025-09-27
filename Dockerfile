# Первый этап: сборка приложения
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/app/main.go

FROM alpine:latest

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]