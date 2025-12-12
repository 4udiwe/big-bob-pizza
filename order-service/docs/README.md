# Swagger Documentation

## Генерация документации

Для генерации Swagger документации используйте команду:

```bash
# Установите swag, если еще не установлен
go install github.com/swaggo/swag/cmd/swag@latest

# Сгенерируйте документацию
swag init -g cmd/main.go -o docs
```

После генерации документация будет доступна:
- Swagger UI: http://localhost:8080/swagger/index.html
- JSON spec: `docs/swagger.json`
- YAML spec: `docs/swagger.yaml`

## Обновление документации

При изменении API (добавлении/изменении эндпоинтов или моделей) необходимо:
1. Обновить Swagger комментарии в хендлерах
2. Запустить `swag init` для регенерации документации
3. Закоммитить обновленные файлы в `docs/`

