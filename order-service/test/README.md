# Тестирование Order Service

Этот документ описывает структуру тестов для Order Service и разницу между типами тестов.

## Типы тестов

### 1. Юнит-тесты (Unit Tests)

**Расположение:** `internal/service/order/service_test.go`

**Что тестируют:**
- Бизнес-логику сервиса изолированно
- Используют моки для всех зависимостей (репозитории, транзакции)
- Быстрые и не требуют внешних зависимостей

**Пример:**
```go
func TestService_CreateOrder(t *testing.T) {
    // Используются моки для OrderRepo, ItemsRepo, OutboxRepo, CacheRepo, Transactor
    // Тестируется только логика сервиса
}
```

**Запуск:**
```bash
make unit-test
# или
go test ./internal/service/order/...
```

### 2. Интеграционные тесты (Integration Tests)

**Расположение:** `test/integration_test/`

**Что тестируют:**
- Взаимодействие сервиса с реальной базой данных (PostgreSQL) и Redis
- Проверяют корректность работы репозиториев и кэша
- Тестируют транзакции и миграции

**Отличие от юнит-тестов:**
- Используют **реальные** PostgreSQL и Redis (в Docker)
- Тестируют **интеграцию** между слоями (сервис → репозиторий → БД)
- Проверяют работу с реальными данными

**Пример:**
```go
func TestService_CreateOrder_Integration(t *testing.T) {
    // Используется реальный testPostgres и testRedis
    // Тестируется полный цикл: сервис → репозиторий → БД
}
```

**Запуск:**
```bash
make integration-test
# или
docker compose -f docker-compose.test.yaml up --abort-on-container-exit --exit-code-from integration
```

### 3. E2E тесты (End-to-End Tests)

**Расположение:** `test/e2e_test/`

**Что тестируют:**
- Полный flow через HTTP API
- Тестируют весь стек: HTTP → Handler → Service → Repository → БД
- Проверяют корректность API контрактов и валидации

**Отличие от интеграционных тестов:**
- Тестируют через **HTTP запросы** (как реальный клиент)
- Запускают **полный сервис** (HTTP сервер + все зависимости)
- Проверяют **внешний интерфейс** (API endpoints)
- Тестируют **валидацию** и **обработку ошибок** на уровне HTTP

**Пример:**
```go
func TestCreateOrder_E2E(t *testing.T) {
    // HTTP запрос к реальному серверу
    err := Do(
        Post(basePath+"/orders"),
        Send().Body().JSON(req),
        Expect().Status().Equal(http.StatusCreated),
    )
}
```

**Запуск:**
```bash
make e2e-test
# или
docker compose -f docker-compose.test.yaml up --abort-on-container-exit --exit-code-from e2e
```

## Сравнение типов тестов

| Характеристика | Юнит-тесты | Интеграционные тесты | E2E тесты |
|----------------|------------|----------------------|-----------|
| **Зависимости** | Моки | Реальная БД + Redis | Полный стек |
| **Скорость** | Очень быстро | Средне | Медленно |
| **Изоляция** | Полная | Частичная | Минимальная |
| **Что тестируют** | Логика сервиса | Интеграция с БД | Полный flow |
| **Требования** | Нет | Docker (PostgreSQL, Redis) | Docker (весь стек) |

## Примеры сценариев

### Юнит-тест: CreateOrder с ошибкой репозитория
```go
// Мокируем ошибку репозитория
orderRepo.EXPECT().Create(...).Return(entity.Order{}, errors.New("db error"))
// Проверяем, что сервис правильно обрабатывает ошибку
```

### Интеграционный тест: CreateOrder с реальной БД
```go
// Создаем заказ через сервис
created, err := svc.CreateOrder(ctx, order)
// Проверяем, что заказ действительно сохранен в БД
retrieved, err := svc.GetOrderByID(ctx, created.ID)
```

### E2E тест: CreateOrder через HTTP
```go
// Отправляем HTTP POST запрос
Do(Post("/orders"), Send().Body().JSON(req))
// Проверяем HTTP статус и тело ответа
Expect().Status().Equal(http.StatusCreated)
```

## Запуск всех тестов

```bash
# Только юнит-тесты
make unit-test

# Интеграционные тесты
make integration-test

# E2E тесты
make e2e-test

# Покрытие кода
make cover
```

## Структура файлов

```
order-service/
├── internal/
│   └── service/
│       └── order/
│           ├── service_test.go      # Юнит-тесты
│           └── mocks/               # Моки для юнит-тестов
├── test/
│   ├── integration_test/            # Интеграционные тесты
│   │   ├── integration_test.go
│   │   ├── service_test.go
│   │   └── Dockerfile
│   └── e2e_test/                    # E2E тесты
│       ├── e2e_test.go
│       ├── order_test.go
│       └── Dockerfile
├── docker-compose.test.yaml         # Docker Compose для тестов
└── Makefile                         # Команды для запуска тестов
```


