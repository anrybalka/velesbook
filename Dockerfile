# Используем официальный образ Golang
FROM golang:1.23.6 AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go.mod и go.sum, чтобы позже скачать зависимости
COPY go.mod go.sum ./

# Устанавливаем и обновляем пакеты сертификатов для Debian
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

# Скачиваем зависимости
RUN go mod download

# Копируем остальные исходники приложения
COPY . .

# Собираем бинарный файл
RUN go build -o app .

# Используем более новый образ с более свежими версиями библиотек
FROM debian:bookworm-slim  

# Копируем скомпилированный бинарник из предыдущего этапа сборки
COPY --from=builder /app/app /app/app

# Указываем рабочую директорию для выполнения
WORKDIR /app

# Открываем порт, если необходимо (например, 8080)
EXPOSE 8080

# Команда для запуска приложения
CMD ["/app/app"]