# Stage 1: Builder
FROM golang:1.25-alpine AS builder

# Установка необходимых зависимостей
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Копируем модули и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/api/main.go

# Stage 2: Финальный образ
FROM scratch

# Копируем сертификаты из builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Копируем бинарный файл
COPY --from=builder /app/main /main

# Порт приложения
EXPOSE 8000

# Запуск приложения
ENTRYPOINT ["/main"]