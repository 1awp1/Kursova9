FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/app/
RUN go build -o main .  # Удалите CGO_ENABLED=0 GOOS=linux GOARCH=amd64

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache libc6-compat  # Убедитесь, что библиотеки libc установлены

COPY ./private_key.pem ./private_key.pem
COPY ./public_key.pem ./public_key.pem
COPY ./static ./static
COPY ./internal/templates ./internal/templates
COPY --from=builder /app/cmd/app/main .

ENTRYPOINT ["./main"]
