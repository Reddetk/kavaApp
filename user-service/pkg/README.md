# Pkg Directory

## Overview
Директория `pkg` содержит вспомогательные пакеты, которые могут быть использованы как внутри сервиса, так и потенциально другими сервисами. Здесь находятся утилиты общего назначения, которые не связаны напрямую с бизнес-логикой приложения.

## Структура директорий

### logger/
Пакет для логирования:

#### logger.go
Реализация логгера:
- Настройка уровней логирования (DEBUG, INFO, WARN, ERROR, FATAL)
- Форматирование логов (текстовый и JSON форматы)
- Ротация логов
- Контекстное логирование с метаданными

**Основные функции**:
- `NewLogger(config LoggerConfig) (Logger, error)`: Создание нового логгера с заданной конфигурацией
- `WithField(key string, value interface{}) Logger`: Добавление поля к логгеру
- `WithFields(fields map[string]interface{}) Logger`: Добавление нескольких полей к логгеру
- `WithError(err error) Logger`: Добавление информации об ошибке к логгеру
- `Debug/Info/Warn/Error/Fatal(msg string)`: Методы для логирования сообщений разных уровней

### utils/
Пакет с различными утилитами общего назначения:

#### string_utils.go
Утилиты для работы со строками:
- Генерация случайных строк
- Проверка форматов строк (email, телефон)
- Форматирование строк
- Преобразование форматов (camelCase, snake_case)
- Маскирование конфиденциальных данных

#### time_utils.go
Утилиты для работы с временем:
- Форматирование и парсинг времени
- Расчет временных интервалов
- Работа с границами периодов (начало/конец дня, месяца)
- Форматирование длительности
- Работа с рабочими днями

#### math_utils.go
Математические утилиты:
- Округление чисел
- Статистические расчеты (среднее, медиана, стандартное отклонение)
- Поиск экстремумов и расчет процентилей
- Нормализация данных
- Обработка выбросов

#### http_utils.go
Утилиты для работы с HTTP:
- Создание и настройка HTTP-клиента
- Выполнение HTTP-запросов
- Работа с URL и параметрами
- Обработка ответов и заголовков
- Проверка статус-кодов

## Принципы проектирования

### Повторное использование
Пакеты в директории `pkg` спроектированы для использования в разных частях приложения и даже в других сервисах. Они предоставляют общую функциональность, которая может понадобиться в различных контекстах.

### Минимальные зависимости
Пакеты имеют минимальные внешние зависимости, что упрощает их интеграцию в различные проекты. Большинство функций используют только стандартную библиотеку Go.

### Тестируемость
Код в директории `pkg` легко тестируется с помощью модульных тестов. Для каждого пакета предусмотрены тесты, которые проверяют корректность работы функций.

### Документированность
Все публичные функции и типы имеют документацию, которая объясняет их назначение, параметры и возвращаемые значения. Это облегчает использование пакетов другими разработчиками.

## Использование

### Импорт пакетов
```go
import (
    "github.com/your-org/user-service/pkg/logger"
    "github.com/your-org/user-service/pkg/utils"
)
```

### Пример использования логгера
```go
// Создание логгера
log := logger.NewLogger(logger.LoggerConfig{
    Level:  "INFO",
    Format: "json",
})

// Логирование с контекстом
log.WithField("component", "main").Info("Application started")
log.WithFields(map[string]interface{}{
    "user_id": "123",
    "action": "login",
}).Info("User logged in")

// Логирование ошибок
if err != nil {
    log.WithError(err).Error("Failed to process request")
}
```

### Пример использования утилит
```go
// Строковые утилиты
email := "user@example.com"
if utils.IsValidEmail(email) {
    phone := "+375291234567"
    formattedPhone := utils.FormatPhone(phone)
    log.WithFields(map[string]interface{}{
        "email": email,
        "phone": formattedPhone,
    }).Info("Contact information validated")
}

// Временные утилиты
now := time.Now()
startOfMonth := utils.StartOfMonth(now)
endOfMonth := utils.EndOfMonth(now)
log.WithFields(map[string]interface{}{
    "start": utils.FormatDate(startOfMonth),
    "end": utils.FormatDate(endOfMonth),
}).Info("Report period")

// Математические утилиты
values := []float64{1.2, 3.4, 5.6, 7.8, 9.0}
mean := utils.Mean(values)
stdDev := utils.StandardDeviation(values)
log.WithFields(map[string]interface{}{
    "mean": mean,
    "stdDev": stdDev,
}).Info("Statistical analysis")

// HTTP-утилиты
client := utils.NewHTTPClient(utils.DefaultHTTPClientConfig())
var response struct {
    Status string `json:"status"`
}
err := utils.GetJSON(context.Background(), client, "https://api.example.com/status", nil, &response)
if err != nil {
    log.WithError(err).Error("Failed to check status")
} else {
    log.WithField("status", response.Status).Info("Service status")
}
```

## Расширение
При необходимости добавления новой функциональности в директорию `pkg`, следуйте этим рекомендациям:

1. Определите, к какому существующему пакету относится новая функциональность, или создайте новый пакет, если она не вписывается в существующие
2. Реализуйте функциональность с минимальными внешними зависимостями
3. Напишите тесты для новой функциональности
4. Добавьте документацию для новых функций и типов
5. Обновите README.md файл соответствующего пакета

## Тестирование
Для запуска тестов всех пакетов в директории `pkg`:

```bash
go test -v ./pkg/...
```

Для запуска тестов конкретного пакета:

```bash
go test -v ./pkg/utils
go test -v ./pkg/logger
```