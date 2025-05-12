# Тесты для репозиториев

В этой директории находятся тесты для репозиториев из пакета `internal/infrastructure/postgres`.

## Структура тестов

### Файлы с тестами
- `transaction_repository_test.go` - тесты для TransactionRepository
- `user_repository_test.go` - тесты для UserRepository
- `segment_repository_test.go` - тесты для SegmentRepository
- `user_metrics_repository_test.go` - тесты для UserMetricsRepository
- `repositories_test.go` - общие тесты для всех репозиториев
- `main_test.go` - тесты для функций из main.go

### Вспомогательные файлы
- `transaction_repository_helpers.go` - вспомогательные функции для тестирования TransactionRepository
- `user_repository_helpers.go` - вспомогательные функции для тестирования UserRepository
- `segment_repository_helpers.go` - вспомогательные функции для тестирования SegmentRepository
- `user_metrics_repository_helpers.go` - вспомогательные функции для тестирования UserMetricsRepository

## Архитектура тестов

Для каждого репозитория создан набор вспомогательных функций:
1. `Setup[Repository]Test` - создает мок базы данных и репозиторий для тестирования
2. `Test[Method]Helper` - тестирует конкретный метод репозитория

Эти вспомогательные функции используются в двух контекстах:
1. В отдельных тестах для каждого метода репозитория (файлы `*_test.go`)
2. В общем тесте для всех репозиториев (`repositories_test.go`)

## Решение проблемы совместимости типов

Для решения проблем совместимости типов были применены следующие подходы:

1. Использование интерфейса вместо конкретной реализации:
   - Используем `repositories.TransactionRepository` вместо `postgres.TransactionRepository`
   - Используем `repositories.UserRepository` вместо `postgres.UserRepository`
   - Используем `repositories.SegmentRepository` вместо `postgres.SegmentRepository`
   - Используем `repositories.UserMetricsRepository` вместо `postgres.UserMetricsRepository`

2. Создание общих вспомогательных функций в отдельных файлах для каждого репозитория

## Запуск тестов

Для запуска всех тестов выполните:

```bash
cd user-service
go test ./test/...
```

Для запуска тестов конкретного репозитория:

```bash
cd user-service
go test ./test -run TestUserRepository
```

Для запуска конкретного теста:

```bash
cd user-service
go test ./test -run TestUserRepository_Get_Standalone
```

## Зависимости

Для запуска тестов необходимы следующие зависимости:

```bash
go get github.com/DATA-DOG/go-sqlmock
go get github.com/stretchr/testify/assert
```

## Примечания

Тесты используют библиотеку `sqlmock` для имитации взаимодействия с базой данных, что позволяет тестировать репозитории без реальной базы данных.