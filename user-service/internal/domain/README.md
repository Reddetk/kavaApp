# Domain Directory

## Overview
Директория `domain` содержит ядро бизнес-логики приложения. Здесь определены основные бизнес-сущности, интерфейсы репозиториев для работы с данными и интерфейсы сервисов для реализации бизнес-правил. Этот слой не зависит от внешних фреймворков и библиотек.

## Структура директорий

### entities/
Содержит определения бизнес-сущностей:

#### user.go
Определяет сущность пользователя:
- Структура `User` с полями: ID, имя, email, дата регистрации, статус и т.д.
- Методы валидации пользователя
- Константы для статусов пользователя

#### transaction.go
Определяет сущность транзакции:
- Структура `Transaction` с полями: ID, ID пользователя, сумма, дата, тип и т.д.
- Методы для расчета агрегированных показателей
- Константы для типов транзакций

#### user-metrics.go
Определяет метрики поведения пользователя:
- Структура `UserMetrics` с полями: ID пользователя, частота покупок, средний чек, время с последней покупки и т.д.
- Методы для расчета RFM-показателей (Recency, Frequency, Monetary)
- Методы для анализа поведенческих паттернов

#### segment.go
Определяет сегменты пользователей:
- Структура `Segment` с полями: ID, название, описание, тип, критерии и т.д.
- Структура `UserSegment` для связи пользователя с сегментом
- Константы для типов сегментов

#### clv.go
Определяет структуры для работы с пожизненной ценностью клиента (CLV):
- Структура `CLV` с полями: ID пользователя, значение CLV, дата расчета, прогноз и т.д.
- Структура `CLVDataPoint` для хранения исторических значений CLV
- Методы для расчета и анализа CLV

### repositories/
Содержит интерфейсы для доступа к данным:

#### user_repository.go
Интерфейс для работы с пользователями:
```go
type UserRepository interface {
    Create(user *entities.User) (*entities.User, error)
    Get(id string) (*entities.User, error)
    Update(user *entities.User) (*entities.User, error)
    Delete(id string) error
    List(filter UserFilter) ([]*entities.User, error)
}
```

#### transaction_repository.go
Интерфейс для работы с транзакциями:
```go
type TransactionRepository interface {
    Create(transaction *entities.Transaction) (*entities.Transaction, error)
    Get(id string) (*entities.Transaction, error)
    GetByUser(userId string) ([]*entities.Transaction, error)
    Update(transaction *entities.Transaction) (*entities.Transaction, error)
    Delete(id string) error
}
```

#### user_metrics_repository.go
Интерфейс для работы с метриками пользователей:
```go
type UserMetricsRepository interface {
    Create(metrics *entities.UserMetrics) (*entities.UserMetrics, error)
    Get(userId string) (*entities.UserMetrics, error)
    Update(metrics *entities.UserMetrics) (*entities.UserMetrics, error)
    Delete(userId string) error
    ListBySegment(segmentId string) ([]*entities.UserMetrics, error)
}
```

#### segment_repository.go
Интерфейс для работы с сегментами:
```go
type SegmentRepository interface {
    Create(segment *entities.Segment) (*entities.Segment, error)
    Get(id string) (*entities.Segment, error)
    Update(segment *entities.Segment) (*entities.Segment, error)
    Delete(id string) error
    List(filter SegmentFilter) ([]*entities.Segment, error)
    AssignUserToSegment(userId string, segmentId string) error
    RemoveUserFromSegment(userId string, segmentId string) error
    GetUserSegments(userId string) ([]*entities.Segment, error)
}
```

#### clv_repository.go
Интерфейс для работы с данными CLV:
```go
type CLVRepository interface {
    Create(clv *entities.CLV) (*entities.CLV, error)
    Get(userId string) (*entities.CLV, error)
    Update(clv *entities.CLV) (*entities.CLV, error)
    Delete(userId string) error
    GetHistory(userId string) ([]*entities.CLVDataPoint, error)
    AddHistoryPoint(dataPoint *entities.CLVDataPoint) error
}
```

### services/
Содержит интерфейсы бизнес-логики:

#### segmentation_service.go
Интерфейс для сегментации пользователей:
```go
type SegmentationService interface {
    PerformRFMClustering(metrics []*entities.UserMetrics) (map[string]string, error)
    PerformBehavioralClustering(metrics []*entities.UserMetrics) (map[string]string, error)
}
```

#### survival_analysis_service.go
Интерфейс для анализа выживаемости:
```go
type SurvivalAnalysisService interface {
    CalculateChurnProbability(metrics *entities.UserMetrics) (float64, error)
    PredictTimeToEvent(metrics *entities.UserMetrics) (time.Duration, error)
    FitModel(metrics []*entities.UserMetrics) error
}
```

#### state_transition_service.go
Интерфейс для моделирования переходов между состояниями:
```go
type StateTransitionService interface {
    BuildTransitionMatrix(transactions []*entities.Transaction) error
    PredictNextState(userId string, currentState string) (string, float64, error)
    CalculateTransitionProbabilities(userId string) (map[string]map[string]float64, error)
}
```

#### clv_service.go
Интерфейс для расчета CLV:
```go
type CLVService interface {
    Calculate(metrics *entities.UserMetrics, transactions []*entities.Transaction) (*entities.CLV, error)
    EstimateWithScenario(metrics *entities.UserMetrics, transactions []*entities.Transaction, scenario string) (*entities.CLV, error)
}
```

## Принципы проектирования
1. **Независимость от инфраструктуры**: Доменный слой не зависит от конкретных технологий и фреймворков
2. **Богатая доменная модель**: Сущности содержат не только данные, но и методы для работы с ними
3. **Инверсия зависимостей**: Доменный слой определяет интерфейсы, которые реализуются внешними слоями
4. **Чистые интерфейсы**: Интерфейсы репозиториев и сервисов определяют только необходимые методы без деталей реализации