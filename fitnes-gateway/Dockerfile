# Используем официальный образ Golang в качестве базового образа
FROM golang:latest


RUN mkdir /app

COPY .dockerignore /app/.dockerignore
COPY ./go.mod /app/go.mod
COPY ./ /app

WORKDIR /app
# Собираем приложение внутри контейнера
RUN go build -o app ./cmd

# Запускаем приложение по умолчанию при запуске контейнера
CMD ["./app"]
