# Big Bob Pizza - Microservices Architecture

Система управления заказами пиццерии на основе микросервисной архитектуры.

## Архитектура

Система состоит из следующих микросервисов:

- **order-service** (порт 8080) - управление заказами
- **payment-service** (порт 8081) - обработка платежей

## Запуск через Docker Compose

### Быстрый старт

```bash
docker-compose up -d
```

Это запустит все необходимые сервисы:
- PostgreSQL для order-service (порт 5432)
- PostgreSQL для payment-service (порт 5433)
- Redis для order-service (порт 6379)
- Kafka + Zookeeper (порты 9092, 2181)
- Kafka UI (порт 8082)
- order-service (порт 8080)
- payment-service (порт 8081)

### Переменные окружения

Вы можете настроить сервисы через переменные окружения или создать `.env` файл:

```env
# Order Service Database
ORDER_DB_USER=order
ORDER_DB_PASSWORD=orderpass
ORDER_DB_NAME=orderdb
ORDER_DB_PORT=5432
ORDER_APP_PORT=8080
ORDER_SERVER_PORT=8080

# Payment Service Database
PAYMENT_DB_USER=payment
PAYMENT_DB_PASSWORD=paymentpass
PAYMENT_DB_NAME=paymentdb
PAYMENT_DB_PORT=5433
PAYMENT_APP_PORT=8081
PAYMENT_SERVER_PORT=8081

# Redis
REDIS_PORT=6379

# Kafka UI
KAFKA_UI_PORT=8082
```


## API Endpoints

### Order Service (http://localhost:8080)

- `POST /orders` - создание заказа

Подробнее [order-service/](docs/EVENT_CATALOG.md)

### Payment Service (http://localhost:8081)

- `POST /payments` - проведение оплаты заказа
- `GET /health` - health check

## Kafka Topics

- `order.events` - события жизненного цикла заказа
- `payment.events` - события оплаты

Подробнее см. [docs/EVENT_CATALOG.md](docs/EVENT_CATALOG.md)
