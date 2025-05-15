# User Service API

Сервис для управления пользователями, сегментацией, расчетом CLV и прогнозированием удержания клиентов.

## Архитектура

Сервис построен с использованием принципов чистой архитектуры (Clean Architecture) и разделен на следующие слои:

- **Domain** - бизнес-сущности и интерфейсы репозиториев
- **Application** - сервисы приложения, реализующие бизнес-логику
- **Infrastructure** - реализации репозиториев и внешних сервисов
- **Interfaces** - HTTP обработчики и маршрутизация

## API Endpoints

### Пользователи

#### GET /api/v1/users/:id
Получение информации о пользователе по ID.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+375261234567",
  "age": 30,
  "gender": "male",
  "city": "Minsk",
  "registration_date": "2023-01-01T12:00:00Z",
  "last_activity": "2023-06-01T15:30:00Z"
}
```

#### POST /api/v1/users
Создание нового пользователя.

**Тело запроса:**
```json
{
  "email": "user@example.com",
  "phone": "+375261234567",
  "age": 30,
  "gender": "male",
  "city": "Minsk"
}
```

**Ответ:**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "phone": "+375261234567",
  "age": 30,
  "gender": "male",
  "city": "Minsk",
  "registration_date": "2023-06-01T15:30:00Z",
  "last_activity": "2023-06-01T15:30:00Z"
}
```

#### PUT /api/v1/users/:id
Обновление информации о пользователе.

**Параметры:**
- `id` - UUID пользователя

**Тело запроса:**
```json
{
  "email": "updated@example.com",
  "city": "New York"
}
```

**Ответ:**
```json
{
  "id": "uuid",
  "email": "updated@example.com",
  "phone": "+375261234567",
  "age": 30,
  "gender": "male",
  "city": "New York"
}
```

### Сегменты

#### POST /api/v1/segments/rfm
Запуск RFM сегментации для всех пользователей.

**Ответ:**
```json
{
  "message": "RFM segmentation completed"
}
```

#### POST /api/v1/segments/behavior
Запуск поведенческой сегментации.

**Ответ:**
```json
{
  "message": "Behavior segmentation completed"
}
```

#### POST /api/v1/segments
Создание нового сегмента.

**Тело запроса:**
```json
{
  "name": "High Value Customers",
  "type": "rfm"
}
```

**Ответ:**
```json
{
  "message": "Segment created",
  "id": "uuid"
}
```

#### PUT /api/v1/segments
Обновление существующего сегмента.

**Тело запроса:**
```json
{
  "id": "uuid",
  "name": "Updated Segment Name",
  "type": "rfm"
}
```

**Ответ:**
```json
{
  "message": "Segment updated"
}
```

#### GET /api/v1/segments
Получение всех сегментов определенного типа.

**Параметры запроса:**
- `type` - тип сегмента (например, "rfm", "behavior")

**Ответ:**
```json
[
  {
    "id": "uuid1",
    "name": "High Value",
    "type": "rfm"
  },
  {
    "id": "uuid2",
    "name": "Medium Value",
    "type": "rfm"
  }
]
```

#### GET /api/v1/segments/:id
Получение информации о сегменте по ID.

**Параметры:**
- `id` - UUID сегмента

**Ответ:**
```json
{
  "id": "uuid",
  "name": "High Value",
  "type": "rfm"
}
```

#### PUT /api/v1/segments/assign/:id
Назначение пользователя в сегмент.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "message": "User assigned to segment"
}
```

#### GET /api/v1/segments/user/:id
Получение сегмента пользователя.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "id": "uuid",
  "name": "High Value",
  "type": "rfm"
}
```

### Удержание

#### GET /api/v1/retention/:id/churn
Прогноз вероятности оттока для пользователя.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "user_id": "uuid",
  "churn_probability": 0.75,
  "churn_probability_score": "high"
}
```

#### GET /api/v1/retention/:id/time
Прогноз времени до события (оттока).

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "user_id": "uuid",
  "time_to_event": 1296000000000000,
  "estimated_days": "within a month"
}
```

### CLV (Customer Lifetime Value)

#### GET /api/v1/clv/:id
Расчет CLV для пользователя.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "user_id": "uuid",
  "value": 1000.0,
  "currency": "USD",
  "calculated_at": "2023-06-01T15:30:00Z",
  "forecast": 1200.0,
  "confidence": 0.85,
  "scenario": "default"
}
```

#### POST /api/v1/clv/update
Пакетное обновление CLV для всех пользователей.

**Тело запроса:**
```json
{
  "batch_size": 100
}
```

**Ответ:**
```json
{
  "message": "Batch CLV update completed"
}
```

#### GET /api/v1/clv/:id/estimate
Оценка CLV по сценарию.

**Параметры:**
- `id` - UUID пользователя
- `scenario` - сценарий оценки (например, "optimistic", "pessimistic", "default")

**Ответ:**
```json
{
  "user_id": "uuid",
  "value": 1200.0,
  "currency": "USD",
  "calculated_at": "2023-06-01T15:30:00Z",
  "forecast": 1500.0,
  "confidence": 0.8,
  "scenario": "optimistic"
}
```

#### GET /api/v1/clv/:id/history
Получение исторических данных CLV.

**Параметры:**
- `id` - UUID пользователя

**Ответ:**
```json
{
  "user_id": "uuid",
  "history": [
    {
      "date": "2023-01-01",
      "value": 800.0,
      "scenario": "default"
    },
    {
      "date": "2023-06-01",
      "value": 900.0,
      "scenario": "default"
    },
    {
      "date": "2023-12-01",
      "value": 1000.0,
      "scenario": "default"
    }
  ]
}
```

## Тестирование

Для тестирования API используются юнит-тесты с применением паттерна адаптера для моков сервисов. Тесты находятся в директории `/test`.

### Запуск тестов

```bash
go test ./test/...
```