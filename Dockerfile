# Указываем базовый образ для сборки
FROM golang:1.19 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Сборка CLI-приложения
RUN go build -o burden cmd/main.go

# Финальный этап: минимальный образ для запуска
FROM alpine:latest

# Копируем скомпилированный бинарный файл из builder
COPY --from=builder /app/burden /usr/local/bin/burden

# Задание дефолтной команды при запуске
ENTRYPOINT ["burden"]