# Menu Service

Сервис управления меню и акциями для системы Big Bob Pizza.

## Описание

Menu-service отвечает за управление блюдами и акциями. Это изолированный сервис, который не взаимодействует с другими микросервисами и не публикует события.

## Команды (Commands)

- **Добавить блюдо** (`POST /dishes`) - создает новое блюдо в меню
- **Изменить блюдо** (`PUT /dishes/{dish_id}`) - обновляет информацию о блюде
- **Удалить блюдо** (`DELETE /dishes/{dish_id}`) - удаляет блюдо из меню
- **Активировать акцию** (`POST /promotions/{promotion_id}/activate`) - активирует акцию (только для админа)

## Архитектура

Сервис построен по слоистой архитектуре:

```
menu_service/
├── domain/          # Доменные сущности (Dish, Promotion)
├── repositories/    # Репозитории для работы с БД
├── services/        # Бизнес-логика (MenuService)
├── handlers/        # HTTP handlers (FastAPI endpoints)
└── database/        # Модели БД и подключение
```

## Технологии

- **FastAPI** - веб-фреймворк
- **SQLAlchemy** - ORM
- **asyncpg** - асинхронный драйвер PostgreSQL
- **Alembic** - миграции БД

## Запуск

```bash
docker-compose up -d --build
```

## Документация

Документация доступна по http://localhost:8000/docs или в файле [openapi](openapi.json)

## Установка и запуск

### Требования

- Python 3.11+
- PostgreSQL 14+ (или Docker)

### Установка зависимостей

```bash
pip install -r requirements.txt
```

### Настройка

Создайте файл `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=menu_db
```

### Запуск PostgreSQL

**Вариант 1: Docker (рекомендуется)**
```bash
docker-compose up -d
```

**Вариант 2: Локальный PostgreSQL**
```bash
# Создайте базу данных
psql -U postgres
CREATE DATABASE menu_db;
\q
```

### Запуск сервиса

```bash
python main.py
```

Или через uvicorn:

```bash
uvicorn main:app --host 0.0.0.0 --port 8000
```

## API Endpoints

### Блюда

- `GET /dishes` - получить все блюда
- `GET /dishes/{dish_id}` - получить блюдо по ID
- `POST /dishes` - создать блюдо
- `PUT /dishes/{dish_id}` - обновить блюдо
- `DELETE /dishes/{dish_id}` - удалить блюдо

### Акции

- `GET /promotions/{promotion_id}` - получить акцию по ID
- `POST /promotions/{promotion_id}/activate` - активировать акцию

### Health Check

- `GET /health` - проверка здоровья сервиса

## Структура БД

- `dishes` - блюда
- `promotions` - акции

## Проверка работоспособности

### Быстрая проверка

1. **Health Check:**
   ```bash
   curl http://localhost:8000/health
   ```

2. **Swagger UI:**
   Откройте в браузере: http://localhost:8000/docs

3. **Автоматическое тестирование:**
   ```bash
   pip install httpx
   python test_service.py
   ```

Подробная инструкция в файле [CHECK_SERVICE.md](CHECK_SERVICE.md)

## Разработка

### Структура проекта

Проект следует принципам Clean Architecture и DDD:

- **Domain Layer** - доменные сущности без зависимостей
- **Repository Layer** - абстракции для работы с данными
- **Service Layer** - бизнес-логика и координация
- **Handler Layer** - HTTP endpoints

## Лицензия

MIT

