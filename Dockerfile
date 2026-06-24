# --- Этап 1: Сборка бинарника ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем исходный код
COPY main.go .

# Компилируем под Linux со статическими линками
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --- Этап 2: Финальный легковесный образ ---
FROM alpine:latest

WORKDIR /app

# Копируем только скомпилированный файл из прошлого этапа
COPY --from=builder /app/server .

# Открываем твой порт 4321
EXPOSE 4321

# Запускаем сервер
CMD ["./server"]