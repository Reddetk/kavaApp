# Infrastructure Directory

## Overview
Директория `infrastructure` содержит реализации интерфейсов, определенных в доменном слое. Здесь находятся конкретные реализации репозиториев для работы с базой данных PostgreSQL и реализации сервисов для выполнения различных алгоритмов.

## Структура директорий

### postgres/
Содержит реализации репозиториев для работы с PostgreSQL:

#### user_repository.go
Реализация интерфейса `repositories.UserRepository`:
- Методы для создания, получения, обновления и удаления пользователей
- SQL-запросы для работы с таблицей пользователей
- Обработка ошибок базы данных и преобразование их в доменные ошибки

#### transaction_repository.go
Реализация интерфейса `repositories.TransactionRepository`:
- Методы для работы с транзакциями пользователей
- SQL-запросы для работы с таблицей транзакций
- Агрегация транзакций по различным критериям

#### user_metrics_repository.go
Реализация интерфейса `repositories.UserMetricsRepository`:
- Методы для работы с метриками пользователей
- SQL-запросы для работы с таблицей метрик
- Обновление метрик на основе новых данных

#### segment_repository.go
Реализация интерфейса `repositories.SegmentRepository`:
- Методы для работы с сегментами пользователей
- SQL-запросы для работы с таблицами сегментов и связей пользователей с сегментами
- Управление назначением пользователей в сегменты

### services/
Содержит реализации сервисов для выполнения различных алгоритмов:

#### kmeans_segmentation.go
Реализация интерфейса `services.SegmentationService` с использованием алгоритма K-means:
- Кластеризация пользователей на основе их метрик
- Настройка параметров алгоритма K-means
- Интерпретация результатов кластеризации

#### cox_survival_analysis.go
Реализация интерфейса `services.SurvivalAnalysisService` с использованием модели пропорциональных рисков Кокса:
- Анализ выживаемости для прогнозирования оттока пользователей
- Расчет вероятности оттока
- Прогнозирование времени до следующей покупки

#### markov_transition_service.go
Реализация интерфейса `services.StateTransitionService` с использованием цепей Маркова:
- Построение матрицы переходов между состояниями пользователей
- Прогнозирование следующего состояния пользователя
- Расчет вероятностей переходов между состояниями

#### discounted_clv_service.go
Реализация интерфейса `services.CLVService` с использованием метода дисконтированных денежных потоков:
- Расчет пожизненной ценности клиента (CLV)
- Учет дисконтирования будущих доходов
- Анализ различных сценариев (базовый, оптимистичный, пессимистичный)

## Технические детали

### Работа с базой данных
- Использование стандартного пакета `database/sql` для работы с PostgreSQL
- Управление соединениями с базой данных через пул соединений
- Транзакционная обработка операций, затрагивающих несколько таблиц

### Реализация алгоритмов
- Использование специализированных библиотек для машинного обучения и статистического анализа
- Кеширование результатов вычислений для повышения производительности
- Асинхронное выполнение длительных операций

## Взаимодействие с другими слоями
- Реализует интерфейсы, определенные в `domain/repositories` и `domain/services`
- Использует сущности из `domain/entities`
- Не зависит от слоев `application` и `interfaces`

## Принципы реализации
1. **Инверсия зависимостей**: Реализации зависят от абстракций, а не наоборот
2. **Изоляция от доменной логики**: Инфраструктурный код не содержит бизнес-логики
3. **Обработка ошибок**: Преобразование технических ошибок в доменные
4. **Тестируемость**: Код спроектирован так, чтобы его можно было легко тестировать