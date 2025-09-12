# L0 Project - Тестовое задание

Микросервис для отображения информации о заказах с асинхронной обработкой сообщений.

## Технологии

- Backend: Go, Chi Route, PostgreSQL, Apache Kafka
- Frontend: HTML, CSS, JavaScript
- Инфраструктура: Docker, Docker Compose

## Быстрый старт

### Требования
- Docker
- Docker Compose

### Запуск
```
git clone https://github.com/GkadyrG/L0.git
cd L0
docker-compose up -d
```
### Доступные сервисы
- Frontend: http://localhost:8081
- Backend API: http://localhost:8080
- Kafka UI: http://localhost:8087
- PostgreSQL: localhost:5432

## API
- GET /api/orders/{id} - Получить заказ
- POST /api/orders - Получить превью всех заказов

## Конфиги
- Конфиги хранятся в backend/.env

## Особенности реализации
- Внутренний кэш ускоряет получение данных заказов и снижает нагрузку на базу
- Асинхронная обработка сообщений через Kafka обеспечивает высокую производительность и масштабируемость
- Миграции базы данных реализованы через go-migrate
- Сервис готов к работе в Docker-среде, все зависимости поднимаются через Docker Compose