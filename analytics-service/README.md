# Analytics Service

Сервис аналитики в системе Big Bob Pizza.

## Описание

Analytics Service собирает и хранит данные о событиях заказов для последующего анализа. Сервис слушает события из Kafka и сохраняет их в PostgreSQL, а также экспортирует метрики в Prometheus.

## Собираемые события

Сервис обрабатывает следующие события из топика `order.events`:

- `order.created` - создание заказа
- `order.paid` - оплата заказа
- `order.cancelled` - отмена заказа
- `order.completed` - завершение заказа

## API Документация

Полная документация API доступна через Swagger UI:

http://localhost:8083/swagger/index.html

или в [docs/swagger.yaml](docs/swagger.yaml)

### Основные эндпоинты

- `GET /analytics/stats` - получить статистику по событиям за период
- `GET /analytics/revenue` - получить выручку за период
- `GET /analytics/orders/{orderId}/events` - получить все события для заказа
- `GET /health` - health check
- `GET /metrics` - Prometheus метрики

### Примеры использования

#### Получить статистику за последний месяц

```bash
curl "http://localhost:8083/analytics/stats?startDate=2024-01-01T00:00:00Z&endDate=2024-01-31T23:59:59Z"
```

#### Получить выручку за период

```bash
curl "http://localhost:8083/analytics/revenue?startDate=2024-01-01T00:00:00Z&endDate=2024-01-31T23:59:59Z"
```

#### Получить события заказа

```bash
curl "http://localhost:8083/analytics/orders/{orderId}/events"
```

## Prometheus метрики

Сервис экспортирует следующие метрики:

- `order_events_total{event_type}` - общее количество событий по типам
- `order_amount{event_type}` - гистограмма сумм заказов по типам событий

Метрики доступны по адресу: `http://localhost:8083/metrics`

## Архитектура

### Хранение данных

- **PostgreSQL** - для хранения детальных данных о событиях
- **Prometheus** - для хранения временных рядов метрик

### Обработка событий

1. Consumer получает события из Kafka топика `order.events`
2. События сохраняются в PostgreSQL (с дедупликацией по `event_id`)
3. Обновляются Prometheus метрики

### База данных

Таблица `order_events` хранит:
- Информацию о событии (тип, время)
- Связанные сущности (заказ, пользователь, платеж)
- Дополнительные данные (сумма, причина отмены)

## Конфигурация

Основные параметры конфигурации в `config/config.yaml`:

```yaml
app:
  name: "big-bob-pizza-analytics-service"
  version: "1.0.0"

http:
  port: "8083"

postgres:
  connect_timeout: 5s

kafka:
  brokers:
    - "localhost:9092"
  topics:
    order_events: "order.events"
  consumer:
    group_id: "analytics-service"

prometheus:
  enabled: true
  path: "/metrics"
```

## Запуск

### Через Docker Compose

```bash
docker-compose up -d analytics-service
```

### Локально

```bash
go run cmd/main.go
```

## Переменные окружения

- `POSTGRES_URL` - строка подключения к PostgreSQL
- `KAFKA_BROKERS` - адреса брокеров Kafka (через запятую)
- `SERVER_PORT` - порт HTTP сервера
- `CONFIG_PATH` - путь к конфигурационному файлу
- `LOG_LEVEL` - уровень логирования (debug, info, warn, error)

