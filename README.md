# Calendar Server (Go)

Простой, надежный и производительный HTTP-сервер для управления событиями календаря, написанный на Go.

## Основные возможности

- Полный набор CRUD операций для событий
- Фильтрация событий по дням, неделям и месяцам
- Структурированное логирование с цветным форматированием
- Graceful shutdown для корректного завершения работы
- Потокобезопасная архитектура
- Поддержка CORS для веб-клиентов
- Валидация входных данных
- Комплексное тестирование с проверкой race conditions

## Технические требования

- Go 1.24 или выше
- Git

## Быстрый старт

### Клонирование репозитория
```bash
git clone https://github.com/sj-shoff/calendar-server.git
cd calendar-server
```

### Запуск в режиме разработки
```bash
make run
```

### Сборка и запуск
```bash
make build
./bin/calendar-server
```

### Запуск с настройками
```bash
PORT=9090 ENVIRONMENT=production ./bin/calendar-server
```

## Тестирование и качество кода

### Запуск тестов
```bash
make test
```

### Проверка API через скрипт
```bash
make test-sh
```

### Статический анализ кода
```bash
make lint
make vet
```

### Установка линтера
```bash
make lint-install
```

## API Endpoints

### Создание события
```
POST /create_event
Content-Type: application/json

{
  "id": "event-1",
  "user_id": "user-123",
  "date": "2025-01-15",
  "title": "Встреча с командой"
}
```

### Обновление события
```
POST /update_event
Content-Type: application/json

{
  "id": "event-1",
  "user_id": "user-123",
  "date": "2025-01-15",
  "title": "Обновленная встреча"
}
```

### Удаление события
```
POST /delete_event
Content-Type: application/json

{
  "id": "event-1"
}
```

### Получение событий за день
```
GET /events_for_day?user_id=user-123&date=2025-01-15
```

### Получение событий за неделю
```
GET /events_for_week?user_id=user-123&date=2025-01-15
```

### Получение событий за месяц
```
GET /events_for_month?user_id=user-123&date=2025-01-15
```

## Формат ответов

### Успешный ответ
```json
{
  "result": "event created"
}
```

### Ошибка
```json
{
  "error": "event ID cannot be empty"
}
```

## Конфигурация

Сервер поддерживает настройку через флаги командной строки и переменные окружения:

- `PORT` - порт сервера (по умолчанию: 8888)
- `ENVIRONMENT` - среда выполнения (development/production)

Примеры использования:
```bash
# Через флаги
./calendar-server -port=9090 -env=production

# Через переменные окружения
PORT=9090 ENVIRONMENT=production ./calendar-server

# Комбинированный способ
PORT=9090 ./calendar-server -env=production
```

## Архитектура приложения

```
calendar-server/
├── cmd/calendar-server/          # Точка входа приложения
├── internal/                     # Внутренние пакеты приложения
│   ├── app/                      # Инициализация и запуск приложения
│   ├── config/                   # Управление конфигурацией
│   ├── domain/                   # Бизнес-сущности
│   ├── delivery/                 # Слой доставки
│   │   └── http-server/          # HTTP-сервер
│   │       ├── handler/          # Обработчики HTTP-запросов
│   │       ├── middleware/       # Промежуточное ПО
│   │       └── router/           # Маршрутизация
│   ├── usecase/                  # Бизнес-логика
│   │   └── event_usecase/        # Use cases для событий
│   └── repository/               # Слой данных
│       └── event_repository/     # Репозиторий событий
│           └── inmemory/         # In-memory реализация
├── pkg/                          # Вспомогательные пакеты
│   ├── errors/                   # Кастомные ошибки
│   └── logger/                   # Логирование
├── Makefile                      # Автоматизация
├── README.md                     # Документация
├── go.mod                        # Модуль Go
└── test_api.sh                   # Скрипт тестирования API
```

## Особенности реализации

### Логирование
- Структурированные логи с использованием Zap
- Цветное форматирование в режиме разработки
- Логирование HTTP-запросов через middleware

### Безопасность
- Валидация всех входных данных
- Обработка CORS для кросc-доменных запросов
- Защита от race conditions

### Надежность
- Graceful shutdown при получении сигналов OS
- Обработка таймаутов подключений
- Комплексная обработка ошибок

### Тестирование
- Unit-тесты для всех слоев приложения
- Проверка на race conditions
- Покрытие кода более 70%
