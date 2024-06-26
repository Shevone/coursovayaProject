# Веб-приложение для автоматизации работы спортклуба
### 1. Описание проекта
Данный проект представляет собой веб-приложение для автоматизации работы спортклуба. Приложение реализовано с использованием архитектуры микросервисов и использует gRPC для коммуникации между сервисами.




### 2. Технологии
#### Backend:
- Go (Golang)
- PostgreSQL
- gRPC
- Protobuf
- Docker
- Http
#### Frontend:
- HTML
- CSS
- JavaScript


### Функционал
#### Пользователи:
- Регистрация и авторизация
- Управление профилем
- Просмотр расписания тренировок
- Запись на тренировки
- Отмена записи
- Просмотр своих записей
#### Тренеры:
- Создание и редактирование тренировок
- Управление списком тренировок

#### Администратор:
- Управление пользователями (добавление, редактирование, удаление)
- Управление тренерами (добавление, редактирование, удаление)
- Управление расписанием
- Настройка системы
- Инструкция по запуску
- Инициализация базы данных

### 3. Локальный запуск и тестирование
Сборка и запуск приложения с помощью Docker Compose
Откройте командную строку/терминал в корневой директории проекта.

_Выполните команду_

``docker-compose up --build``

После успешного запуска приложение будет доступно по адресу http://localhost:8080.
#### Запуск клиентского приложения
Откройте fitnes-client/index.html в браузере.
