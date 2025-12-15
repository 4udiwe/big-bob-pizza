# Grafana Configuration

Эта директория содержит конфигурацию Grafana для автоматической настройки источников данных и дашбордов.

## Структура

```
grafana/
├── provisioning/
│   ├── datasources/
│   │   └── datasources.yml    # Автоматическая настройка Prometheus datasource
│   └── dashboards/
│       └── dashboards.yml     # Конфигурация загрузки дашбордов
└── dashboards/
    └── analytics-overview.json # Дашборд с метриками аналитики
```

## Datasources

В `provisioning/datasources/datasources.yml` настроен источник данных Prometheus, который автоматически подключается к Prometheus серверу.

## Dashboards

Дашборды автоматически загружаются из директории `dashboards/` при запуске Grafana.

### Analytics Overview

Дашборд "Big Bob Pizza - Analytics Overview" содержит следующие панели:

1. **Order Events Rate** - график скорости событий заказов (events/sec)
2. **Total Order Events** - общее количество событий
3. **Events by Type** - круговая диаграмма распределения событий по типам
4. **Order Amount Percentiles** - процентили сумм заказов
5. **Average Order Amount** - средняя сумма заказа
6. **Order Events Count by Type** - столбчатая диаграмма событий по типам
7. **Created Orders** - количество созданных заказов
8. **Paid Orders** - количество оплаченных заказов
9. **Completed Orders** - количество завершенных заказов
10. **Cancelled Orders** - количество отмененных заказов
11. **Total Revenue** - общая выручка
12. **Order Events Timeline** - временная шкала событий

## Использование

После запуска `docker-compose up`, Grafana автоматически:
- Настроит источник данных Prometheus
- Загрузит все дашборды из директории `dashboards/`

Для доступа к Grafana:
- URL: http://localhost:3000
- Логин: admin
- Пароль: admin (или значение из переменной окружения GRAFANA_PASSWORD)

