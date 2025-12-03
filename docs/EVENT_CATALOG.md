# Каталог событий и топиков для взаимодействия микросервисов

Этот файл является **единым источником правды** по взаимодействию сервисов через Kafka.
Он описывает:

* Kafka топики
* События и их payload
* Кто публикует, кто потребляет

# Event Envelope (единый формат всех сообщений Kafka)

Все события в публикуются в формате JSON-обёртки:
```json
{
  "eventId": "UUID",
  "eventType": "string",
  "occurredAt": "RFC3339",
  "data": { ... payload ... }
}
```

- event_id — уникальный UUID события
- event_type — тип события (например: "order.created")
- occurred_at — точное время возникновения события в домене
- data — конкретный payload (структура описана ниже для каждого события)


# Topics

Ниже перечислены **все топики Kafka**, используемые в системе.

## Список топиков

| Топик                                             | Описание                        | Публикует сервис     |
| ------------------------------------------------- | ------------------------------- | -------------------- |
| [order.events](#topic-orderevents)                | События жизненного цикла заказа | order                |
| [payment.events](#topic-paymentevents)            | Статус оплаты                   | payment              |
| [kitchen.events](#topic-kitchenevents)            | Этапы приготовления             | kitchen              |
| [delivery.events](#topic-deliveryevents)          | Этапы доставки                  | delivery             |


# TOPIC: **order.events**

## Event: `order.created`

- **Описание:** Создан новый заказ пользователем
- **Публикует:** order-service
- **Слушают:** payment, analytics

### Payload:

```json
{
  "orderId": "UUID",
  "userId": "UUID",
  "items": [
    { "dishId": "UUID", "quantity": 1 }
  ],
  "totalPrice": 123.45
}
```

---

## Event: `order.paid`

- **Описание:** Заказ оплачен, после оплаты заказ принимается на готовку кухней.
- **Публикует:** order-service
- **Слушают:** kitchen, analytics

### Payload:

```json
{
  "orderId": "UUID",
  "paymentId": "UUID"
}
```

---

## Event: `order.cancelled`

- **Описание:** Заказ отменён
- **Публикует:** order-service
- **Слушают:** analytics

### Payload:

```json
{
  "orderId": "UUID",
  "reason": "string"
}
```

---

## Event: `order.prepeared`

- **Описание:** Заказ приготовлен кухней
- **Публикует:** order-service
- **Слушают:** notification

```json
{
  "orderId": "UUID"
}
```

---

## Event: `order.delivering`

- **Описание:** Заказ передан курьеру
- **Публикует:** order-service
- **Слушают:** notification

```json
{
  "orderId": "UUID"
}
```

---

## Event: `order.completed`

- **Описание:** Заказ доставлен и завершён
- **Публикует:** order-service
- **Слушают:** notification, analytics

```json
{
  "orderId": "UUID"
}
```

---

# TOPIC: payment.events

## Event: `payment.success`

- Публикует: payment-service
- Слушают: order

```json
{
  "paymentId": "UUID",
  "orderId": "UUID",
  "amount": 123.45
}
```

---

## Event: `payment.failed`

- Публикует: payment-service
- Слушают: order

```json
{
  "paymentId": "UUID",
  "orderId": "UUID",
  "reason": "string"
}
```

---

# TOPIC: kitchen.events

## Event: `kitchen.accepted`

- Описание: Заказ принят в готовку.
- Публикует: kitchen-service
- Слушают: order

```json
{
  "orderId": "UUID"
}
```

---

## Event: `kitchen.ready`

- Описание: Заказ полностью готов
- Публикует: kitchen-service
- Слушают: delivery-service, order-service

```json
{
  "orderId": "UUID"
}
```

---

## Event: `kitchen.handedToCourier`

- Описание: Заказ передан курьеру
- Публикует: kitchen-service
- Слушают: order-service

```json
{
  "orderId": "UUID"
}
```

---

# TOPIC: delivery.events

## Event: `delivery.assigned`

Публикует: delivery-service
Слушают: notification

```json
{
  "orderId": "UUID",
  "courierId": "UUID"
}
```

---

## Event: `delivery.completed`

Публикует: delivery-service
Слушают: order-service

```json
{
  "orderId": "UUID",
  "deliveredAt": "RFC3339"
}
```

---
