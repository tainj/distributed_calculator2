# Стадия 1: сборка
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Шаг 1: копируем ТОЛЬКО go.mod и go.sum
COPY go.mod go.sum ./

# Шаг 2: загружаем зависимости (этот слой будет закеширован, если mod/sum не менялись)
RUN go mod download

# Шаг 3: копируем ВЕСЬ код (уже после зависимостей)
COPY . .

# Шаг 4: собираем бинарники
RUN go build -o main cmd/main/main.go
RUN go build -o worker cmd/worker/main.go

# Стадия 2: минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарники и .env
COPY --from=builder /app/main ./main
COPY --from=builder /app/worker ./worker
COPY --from=builder /app/.env .env

RUN chmod +x ./main ./worker

EXPOSE 8080 50051 8081 8082 8083

# Команда будет задана в docker-compose