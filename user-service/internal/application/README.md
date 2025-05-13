# Application Directory

## Overview
Директория `application` содержит реализацию бизнес-сценариев (use cases) приложения. Этот слой оркестрирует взаимодействие между различными доменными сервисами и репозиториями для выполнения конкретных задач бизнес-логики.

## Файлы

### user_service.go
Сервис для управления профилями пользователей:
- Создание, обновление и получение профилей пользователей
- Расчет и обновление метрик пользователей
- Координация обновления сегментов пользователей

**Основные методы**:
- `CreateUser(user *entities.User) (*entities.User, error)`
- `GetUser(id string) (*entities.User, error)`
- `UpdateUser(user *entities.User) (*entities.User, error)`
- `DeleteUser(id string) error`
- `CalculateUserMetrics(userId string) (*entities.UserMetrics, error)`

### segmentation-service.go
Сервис для сегментации пользователей:
- Выполнение RFM-сегментации (Recency, Frequency, Monetary)
- Выполнение поведенческой сегментации
- Назначение пользователей в сегменты

**Основные методы**:
- `PerformRFMSegmentation() error`
- `PerformBehavioralSegmentation() error`
- `AssignUserToSegment(userId string, segmentId string) error`
- `GetUserSegments(userId string) ([]*entities.Segment, error)`

### retention_service.go
Сервис для прогнозирования удержания пользователей:
- Расчет вероятности оттока
- Прогнозирование времени до следующей покупки
- Анализ выживаемости пользователей

**Основные методы**:
- `CalculateChurnProbability(userId string) (float64, error)`
- `PredictTimeToNextPurchase(userId string) (time.Duration, error)`
- `PerformSurvivalAnalysis(userIds []string) error`

### clv_service.go
Сервис для расчета пожизненной ценности клиента (CLV):
- Расчет CLV для отдельных пользователей
- Пакетное обновление CLV для всех пользователей
- Оценка CLV по различным сценариям

**Основные методы**:
- `CalculateUserCLV(userId string) (*entities.CLV, error)`
- `BatchUpdateCLV() error`
- `EstimateUserCLV(userId string, scenario string) (*entities.CLV, error)`
- `GetUserCLVHistory(userId string) ([]*entities.CLVDataPoint, error)`

## Взаимодействие с другими слоями

### Зависимости от доменного слоя
- Использует сущности из `domain/entities`
- Использует интерфейсы репозиториев из `domain/repositories`
- Использует интерфейсы сервисов из `domain/services`

### Использование в слое интерфейсов
- Используется HTTP-обработчиками в `interfaces/http/handlers`
- Используется Kafka-консьюмером в `interfaces/kafka/consumer.go`

## Принципы реализации
1. **Оркестрация**: Координирует работу нескольких доменных сервисов и репозиториев
2. **Транзакционность**: Обеспечивает атомарность операций, затрагивающих несколько сущностей
3. **Бизнес-правила**: Реализует бизнес-правила, не относящиеся к конкретной сущности
4. **Независимость от инфраструктуры**: Не зависит от конкретных реализаций репозиториев и сервисов