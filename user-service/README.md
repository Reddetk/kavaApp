# User Service Architecture

## Overview

User Service is a microservice designed to manage user profiles and update customer segments based on their behavior. It implements sophisticated analytics for user segmentation, retention prediction, and Customer Lifetime Value (CLV) calculation.

## Architectural Approach

The service follows Clean Architecture principles with these key layers:

- **Domain Layer**: Core business entities and business rules
- **Application Layer**: Use cases that orchestrate business logic
- **Infrastructure Layer**: Technical implementations of interfaces
- **Interfaces Layer**: External system communication (HTTP, Kafka)

## Project Structure

```
user-service/
├── cmd/
│   └── main.go                        # Application entry point
├── config/
│   ├── config.go                      # Configuration loader
│   └── config.yaml                    # Configuration file
├── internal/
│   ├── application/                   # Use Cases
│   │   ├── clv_service.go             # CLV calculation orchestration
│   │   ├── retention_service.go       # Retention prediction
│   │   ├── segmentation-service.go    # User segmentation
│   │   └── user_service.go            # User management
│   ├── domain/                        # Business entities and logic
│   │   ├── entities/
│   │   │   ├── clv.go                 # CLV data structures
│   │   │   ├── segment.go             # User segment definition
│   │   │   ├── transaction.go         # Transaction records
│   │   │   ├── user-metrics.go        # User behavioral metrics
│   │   │   └── user.go                # User profile
│   │   ├── repositories/              # Repository interfaces
│   │   │   ├── clv_repository.go
│   │   │   ├── segment_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   ├── user_metrics_repository.go
│   │   │   └── user_repository.go
│   │   └── services/                  # Domain service interfaces
│   │       ├── clv_service.go
│   │       ├── segmentation_service.go
│   │       ├── state_transition_service.go
│   │       └── survival_analysis_service.go
│   ├── infrastructure/
│   │   ├── postgres/                  # Repository implementations
│   │   │   ├── segment_repository.go
│   │   │   ├── transaction_repository.go
│   │   │   ├── user_metrics_repository.go
│   │   │   └── user_repository.go
│   │   └── services/                  # Service implementations
│   │       ├── cox_survival_analysis.go
│   │       ├── discounted_clv_service.go
│   │       ├── kmeans_segmentation.go
│   │       └── markov_transition_service.go
│   └── interfaces/                    # External interfaces
│       ├── http/
│       │   ├── handlers/
│       │   │   ├── clv_handler.go
│       │   │   ├── retention_handler.go
│       │   │   ├── segment_handler.go
│       │   │   └── user_handler.go
│       │   └── router.go
│       └── kafka/
│           └── consumer.go            # Kafka event consumer
├── pkg/                               # Helper packages
│   ├── logger/
│   │   └── logger.go
│   └── utils/
├── test/                              # Tests
│   ├── main_test.go                   # Main integration tests
│   ├── repositories_test.go           # Repository tests
│   ├── *_repository_helpers.go        # Test helpers
│   ├── *_repository_test.go           # Individual repository tests
│   └── README.md                      # Testing documentation
├── .env                               # Environment variables
├── .gitignore
├── Dockerfile                         # Docker build instructions
├── go.mod                             # Go dependency manifest
├── go.sum
└── README.md                          # Project documentation
```

## Core Components

### Domain Layer

The domain layer contains business entities and core business rules:

- **Entities**: 
  - `User`: Core user profile information
  - `Transaction`: User purchase records
  - `Segment`: User segment definitions
  - `UserMetrics`: Calculated user behavioral metrics
  - `CLVDataPoint`: Historical CLV data points

- **Repositories**: Interfaces for data access with methods for CRUD operations
  - `UserRepository`
  - `TransactionRepository`
  - `SegmentRepository`
  - `UserMetricsRepository`

- **Services**: Business logic interfaces
  - `SegmentationService`: User clustering and segment assignment
  - `StateTransitionService`: Markov chain modeling for user state transitions
  - `SurvivalAnalysisService`: Churn prediction
  - `CLVService`: Customer Lifetime Value calculation

### Application Layer

The application layer implements use cases by orchestrating domain services:

- **UserService**: User profile management and metrics calculation
- **SegmentationService**: User clustering and segment assignment workflows
- **RetentionService**: Churn prediction and survival analysis coordination
- **CLVService**: Customer Lifetime Value calculation and updates

### Infrastructure Layer

The infrastructure layer provides concrete implementations:

- **Repository Implementations**: PostgreSQL-based data access
  - `postgres.UserRepository`
  - `postgres.TransactionRepository`
  - `postgres.SegmentRepository`
  - `postgres.UserMetricsRepository`

- **Algorithm Implementations**:
  - `KMeansSegmentation`: K-means clustering for user segmentation
  - `CoxSurvivalAnalysis`: Cox proportional hazards model for survival analysis
  - `MarkovTransitionService`: Markov chains for user behavior modeling
  - `DiscountedCLVService`: Discounted cash flow method for CLV calculation

### Interfaces Layer

The interfaces layer handles external communications:

- **HTTP Handlers**: REST API endpoints
  - `UserHandler`: User profile management endpoints
  - `SegmentHandler`: Segment management endpoints
  - `RetentionHandler`: Retention prediction endpoints
  - `CLVHandler`: CLV calculation endpoints

- **Kafka Consumer**: Event processing for transactions
  - Processes transaction events
  - Updates user metrics
  - Triggers re-segmentation when necessary

## Key Features

- **User Segmentation**: 
  - RFM (Recency, Frequency, Monetary) segmentation
  - Behavioral segmentation based on transaction patterns
  - K-means clustering for automatic segment discovery

- **Survival Analysis**: 
  - Predicting user churn probability
  - Estimating time to next purchase
  - Cox proportional hazards model for survival analysis

- **State Transition Analysis**: 
  - Markov chains for user behavior modeling
  - Transition probabilities between user states
  - User state prediction

- **CLV Calculation**: 
  - Discounted cash flow method
  - Multiple scenario analysis (default, optimistic, pessimistic)
  - Historical CLV tracking

## Event Processing

The service consumes transaction events via Kafka to:
- Update user metrics in real-time
- Trigger re-segmentation when necessary
- Recalculate retention probabilities and CLV
- Update user state in Markov model

## Testing

The service includes comprehensive tests:

- **Repository Tests**: Tests for all repository implementations
- **Helper Functions**: Reusable test helpers for common testing scenarios
- **Integration Tests**: Tests for the entire service workflow

Tests use `sqlmock` to simulate database interactions without requiring a real database connection.

## Setup & Deployment

### Prerequisites
- Go 1.16+
- PostgreSQL 12+
- Kafka

### Configuration
Configure database connection, Kafka brokers, and other settings in `config/config.yaml`:

```yaml
server:
  address: ":8080"
  read_timeout_seconds: 15
  write_timeout_seconds: 15
  idle_timeout_seconds: 60

database:
  dsn: "postgres://user:password@localhost:5432/userservice?sslmode=disable"
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime_minutes: 30

kafka:
  brokers:
    - "localhost:9092"
  topic: "transactions"
  group_id: "user-service"

segmentation:
  rfm_clustering:
    algorithm: "KMeans"
    clusters: 5
```

### Running the Service
```bash
# Build
go build -o user-service ./cmd/main.go

# Run
./user-service

# Using Docker
docker build -t user-service .
docker run -p 8080:8080 user-service
```

## Development Guidelines

1. **Clean Architecture**:
   - Keep business logic in domain and application layers
   - Infrastructure and interfaces should depend on inner layers, not vice versa
   - Use dependency injection for service composition

2. **Testing**:
   - Write tests for all repository implementations
   - Use mock objects for external dependencies
   - Test both success and error cases

3. **Error Handling**:
   - Use meaningful error messages
   - Log errors with appropriate context
   - Return appropriate HTTP status codes

4. **Logging**:
   - Use structured logging
   - Include relevant context in log messages
   - Use appropriate log levels

## Recent Updates

- Added Kafka consumer for real-time transaction processing
- Implemented CLV calculation with multiple scenario analysis
- Added comprehensive tests for all repositories
- Improved error handling and logging
- Added graceful shutdown for HTTP server and Kafka consumer
- Fixed issues with segmentation service and CLV calculation
- Updated configuration to support Kafka and database connection pooling

## Contribution

1. Fork this repository
2. Create a feature branch
3. Submit a pull request with detailed descriptions

## License

This project is licensed under the MIT License.
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