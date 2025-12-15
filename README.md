# Big Bob Pizza - Microservices Architecture

Система управления заказами пиццерии на основе микросервисной архитектуры.

## Архитектура

Система состоит из следующих микросервисов:

- **order-service** (порт 8080) - управление заказами
- **payment-service** (порт 8081) - обработка платежей
- **analytics-service** (порт 8083) - сбор и анализ данных о заказах

## Запуск через Docker Compose

### Быстрый старт

```bash
docker-compose up -d
```

Это запустит все необходимые сервисы:
- PostgreSQL для order-service (порт 5432)
- PostgreSQL для payment-service (порт 5433)
- PostgreSQL для analytics-service (порт 5434)
- Redis для order-service (порт 6379)
- Kafka + Zookeeper (порты 9092, 2181)
- Kafka UI (порт 8082)
- Prometheus (порт 9090)
- Grafana (порт 3000)
- order-service (порт 8080)
- payment-service (порт 8081)
- analytics-service (порт 8083)

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

# Analytics Service Database
ANALYTICS_DB_USER=analytics
ANALYTICS_DB_PASSWORD=analyticspass
ANALYTICS_DB_NAME=analyticsdb
ANALYTICS_DB_PORT=5434
ANALYTICS_APP_PORT=8083
ANALYTICS_SERVER_PORT=8083

# Prometheus
PROMETHEUS_PORT=9090

# Grafana
GRAFANA_PORT=3000
GRAFANA_PASSWORD=admin
```


## API Endpoints

### Order Service (http://localhost:8080)

- `POST /orders` - создание заказа

Подробнее [order-service/](docs/EVENT_CATALOG.md)

### Payment Service (http://localhost:8081)

- `POST /payments` - проведение оплаты заказа
- `GET /payments` - получить список платежей
- `GET /payments/{id}` - получить платеж по ID
- `GET /payments/order/{orderId}` - получить платеж по ID заказа
- `GET /health` - health check

### Analytics Service (http://localhost:8083)

- `GET /analytics/stats` - получить статистику по событиям за период
- `GET /analytics/revenue` - получить выручку за период
- `GET /analytics/orders/{orderId}/events` - получить все события для заказа
- `GET /metrics` - Prometheus метрики
- `GET /health` - health check

## Мониторинг

Система включает полноценный стек мониторинга:

### Prometheus (http://localhost:9090)

Собирает метрики из analytics-service:
- `order_events_total{event_type}` - счетчик событий заказов
- `order_amount{event_type}` - гистограмма сумм заказов

### Grafana (http://localhost:3000)

Автоматически настроенная Grafana с предустановленными дашбордами:
- **Логин**: admin
- **Пароль**: admin (или значение из переменной окружения `GRAFANA_PASSWORD`)

#### Доступные дашборды:

- **Big Bob Pizza - Analytics Overview** - обзор всех метрик аналитики:
  - Скорость событий заказов
  - Распределение событий по типам
  - Суммы заказов и процентили
  - Статистика по созданным, оплаченным, завершенным и отмененным заказам
  - Общая выручка
  - Временная шкала событий

Все дашборды и источники данных автоматически настраиваются при запуске через provisioning (см. `grafana/provisioning/`).

## Kafka Topics

- `order.events` - события жизненного цикла заказа
- `payment.events` - события оплаты

Подробнее см. [docs/EVENT_CATALOG.md](docs/EVENT_CATALOG.md)
