# Big Bob Pizza - Microservices Architecture

Система управления заказами пиццерии на основе микросервисной архитектуры.

## Архитектура

Система состоит из следующих микросервисов:

- **order-service** - управление заказами
- **payment-service** - обработка платежей
- **analytics-service** - сбор и анализ данных о заказах

Сервис gateway является единой точкой входа (http://localhost:8080)

## Запуск через Docker Compose

```bash
docker-compose up -d
```

## API Endpoints

### [Order Service](order-service/README.md)

- `POST /orders` - создание заказа


### [Payment Service](payment-service/README.md)

- `POST /payments` - проведение оплаты заказа
- `GET /payments` - получить список платежей
- `GET /payments/{id}` - получить платеж по ID
- `GET /payments/order/{orderId}` - получить платеж по ID заказа
- `GET /health` - health check

### [Analytics Service](analytics-service/README.md)

- `GET /analytics/stats` - получить статистику по событиям за период
- `GET /analytics/revenue` - получить выручку за период
- `GET /analytics/orders/{orderId}/events` - получить все события для заказа
- `GET /metrics` - Prometheus метрики
- `GET /health` - health check

 [Общая документация](docs/swagger.yaml)

## Мониторинг

Система включает полноценный стек мониторинга:

### Prometheus (http://localhost:9090)

Собирает метрики из analytics-service:
- `order_events_total{event_type}` - счетчик событий заказов
- `order_amount{event_type}` - гистограмма сумм заказов

### Grafana (http://localhost:3000)

Автоматически настроенная Grafana с предустановленными дашбордами:

#### Доступные дашборды:

- **Big Bob Pizza - Analytics Overview** - обзор всех метрик аналитики:
  - Скорость событий заказов
  - Распределение событий по типам
  - Суммы заказов и процентили
  - Статистика по созданным, оплаченным, завершенным и отмененным заказам
  - Общая выручка
  - Временная шкала событий

Все дашборды и источники данных автоматически настраиваются при запуске через provisioning (см. `grafana/provisioning/`).

## Kafka

- `order.events` - события жизненного цикла заказа
- `payment.events` - события оплаты
- `kitchen.events` - события приготовления заказа
- `delivery.events`	- события доставки

Подробнее см. [docs/EVENT_CATALOG.md](docs/EVENT_CATALOG.md)
