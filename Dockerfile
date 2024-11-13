# Указываем базовый образ для сборки
FROM golang:1.19 AS builder

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Сборка CLI-приложения
RUN go build -o loadtester main.go

# Финальный этап: минимальный образ для запуска
FROM alpine:latest

# Копируем скомпилированный бинарный файл из builder
COPY --from=builder /app/loadtester /usr/local/bin/loadtester

# Задание дефолтной команды при запуске
ENTRYPOINT ["loadtester"]