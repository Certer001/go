# --- Этап 1: Сборка бинарника внутри контейнера ---
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 1. Автоматически инициализируем модуль прямо в контейнере
RUN go mod init my-hello-server

# 2. Копируем все файлы с расширением .go, которые есть в папке
COPY *.go ./

# [ДОБАВЛЕНО] 3. Скачиваем все внешние библиотеки, которые импортированы в коде
RUN go mod tidy

# 4. Компилируем проект (теперь он соберется без ошибок)
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# --- Этап 2: Финальный легковесный образ ---
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 4321

CMD ["./server"]