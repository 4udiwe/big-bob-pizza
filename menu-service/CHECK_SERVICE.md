# Проверка работоспособности menu-service

## Быстрая проверка

### 1. Health Check

Проверьте, что сервис запущен и отвечает:

```bash
curl http://localhost:8000/health
```

Ожидаемый ответ:
```json
{"status": "ok"}
```

### 2. Проверка через браузер

Откройте в браузере:
- Swagger UI: http://localhost:8000/docs
- ReDoc: http://localhost:8000/redoc

## Автоматическое тестирование

### Запуск тестового скрипта

```bash
# Установите httpx для тестов (если еще не установлен)
pip install httpx

# Запустите тесты
python test_service.py
```

Скрипт проверит:
- ✓ Health check endpoint
- ✓ Создание блюда
- ✓ Получение блюда по ID
- ✓ Получение всех блюд
- ✓ Обновление блюда
- ✓ Удаление блюда
- ✓ Обработку ошибок (404)

## Ручное тестирование через curl

### 1. Создать блюдо

```bash
curl -X POST "http://localhost:8000/dishes" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Пицца Маргарита",
    "description": "Классическая пицца с томатами и моцареллой",
    "price": 599.00,
    "category": "Пицца"
  }'
```

### 2. Получить все блюда

```bash
curl http://localhost:8000/dishes
```

### 3. Получить блюдо по ID

```bash
curl http://localhost:8000/dishes/{dish_id}
```

Замените `{dish_id}` на реальный UUID из ответа создания.

### 4. Обновить блюдо

```bash
curl -X PUT "http://localhost:8000/dishes/{dish_id}" \
  -H "Content-Type: application/json" \
  -d '{
    "price": 649.00,
    "is_available": false
  }'
```

### 5. Удалить блюдо

```bash
curl -X DELETE "http://localhost:8000/dishes/{dish_id}"
```

## Проверка базы данных

### Подключение к PostgreSQL

```bash
psql -h localhost -U postgres -d menu_db
```

### Проверка таблиц

```sql
-- Проверить список таблиц
\dt

-- Проверить блюда
SELECT * FROM dishes;

-- Проверить акции
SELECT * FROM promotions;
```

## Проверка логов

Сервис выводит логи в консоль. Убедитесь, что:
- Нет ошибок подключения к БД
- Нет ошибок при обработке запросов
- Логи показывают успешные операции

## Типичные проблемы

### 1. Сервис не запускается

**Проблема:** Ошибка подключения к БД

**Решение:**
- Проверьте, что PostgreSQL запущен
- Проверьте настройки в `.env` файле
- Убедитесь, что база данных `menu_db` создана

### 2. 500 Internal Server Error

**Проблема:** Ошибка при выполнении запроса

**Решение:**
- Проверьте логи сервиса
- Убедитесь, что миграции применены: `alembic upgrade head`
- Проверьте структуру БД

### 3. 404 Not Found

**Проблема:** Ресурс не найден

**Решение:**
- Проверьте правильность UUID
- Убедитесь, что ресурс существует в БД

### 4. Connection refused

**Проблема:** Сервис не отвечает

**Решение:**
- Убедитесь, что сервис запущен: `python main.py`
- Проверьте порт в настройках (по умолчанию 8000)
- Проверьте firewall настройки

## Проверка производительности

### Простой нагрузочный тест

```bash
# Установите Apache Bench (если нужно)
# Ubuntu/Debian: sudo apt-get install apache2-utils
# macOS: brew install httpd

# Тест health endpoint
ab -n 100 -c 10 http://localhost:8000/health
```

## Мониторинг

### Проверка метрик (если добавлены)

```bash
# Если есть endpoint для метрик
curl http://localhost:8000/metrics
```

## Дополнительные проверки

### Проверка валидации данных

Попробуйте отправить некорректные данные:

```bash
# Отрицательная цена (должна быть ошибка валидации)
curl -X POST "http://localhost:8000/dishes" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Тест",
    "price": -100,
    "category": "Тест"
  }'
```

Ожидается ошибка валидации (422 Unprocessable Entity).

