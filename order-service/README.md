# Order Service

Сервис управления заказами в системе Big Bob Pizza.

## Описание

Order Service отвечает за создание и управление жизненным циклом заказов. Сервис публикует события в Kafka и слушает события от других сервисов (payment, kitchen, delivery) для синхронизации состояния заказов.

## API Документация

Полная документация API доступна через Swagger UI:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI Spec**: см. `docs/swagger.json` (генерируется командой `swag init`)

### Основные эндпоинты

- `POST /orders` - Создать новый заказ
- `GET /orders` - Получить все заказы (админ, с пагинацией)
- `GET /orders/{id}` - Получить заказ по ID
- `GET /orders/user/{userId}` - Получить заказы пользователя (с пагинацией)
- `GET /orders/user/{userId}/active` - Получить активные заказы пользователя
- `GET /health` - Health check

## Особенности работы

### Кэширование

Сервис использует Redis для кэширования активных заказов:
- Активные заказы хранятся в Redis для быстрого доступа
- При создании заказ добавляется в кэш
- При изменении статуса заказ обновляется в кэше
- При завершении заказ удаляется из кэша активных

### События Kafka

Сервис публикует следующие события в топик `order.events`:
- `order.created` - при создании заказа
- `order.paid` - после успешной оплаты
- `order.prepeared` - когда заказ приготовлен
- `order.delivering` - когда заказ передан курьеру
- `order.completed` - когда заказ доставлен
- `order.cancelled` - при отмене заказа

Сервис слушает следующие события:
- `payment.success` из топика `payment.events`
- `payment.failed` из топика `payment.events`
- `kitchen.accepted`, `kitchen.ready`, `kitchen.handedToCourier` из топика `kitchen.events`
- `delivery.completed` из топика `delivery.events`

### Outbox Pattern

Все события публикуются через outbox pattern для гарантированной доставки:
- События записываются в таблицу `outbox` в той же транзакции, что и изменения заказа
- Outbox Worker периодически публикует события в Kafka
- При ошибке публикации события помечаются как failed и перевыставляются позже

### Статусы заказа

- `created` - заказ создан, ожидает оплаты
- `paid` - заказ оплачен
- `preparing` - заказ готовится на кухне
- `prepared` - заказ приготовлен
- `delivering` - заказ доставляется
- `completed` - заказ доставлен и завершен
- `cancelled` - заказ отменен

## Пагинация

Все эндпоинты, возвращающие списки, поддерживают пагинацию через query параметры:
- `limit` - количество записей на странице (по умолчанию 20, максимум 100)
- `offset` - смещение для пагинации (по умолчанию 0)

Пример: `GET /orders?limit=10&offset=20`

## Запуск

```bash
# Локально
go run cmd/main.go

# Через Docker Compose (из корня проекта)
docker-compose up order-service
```

## Конфигурация

Конфигурация находится в `config/config.yaml`. Основные параметры:
- `http.port` - порт HTTP сервера (по умолчанию 8080)
- `postgres.url` - строка подключения к PostgreSQL
- `redis.addr` - адрес Redis сервера
- `kafka.brokers` - список брокеров Kafka
- `outbox.*` - настройки outbox worker

## База данных

Миграции находятся в `internal/database/migrations/` и выполняются автоматически при запуске сервиса.

Основные таблицы:
- `orders` - заказы
- `order_item` - позиции заказов
- `order_status` - справочник статусов
- `outbox` - события для публикации
- `outbox_status` - справочник статусов outbox

