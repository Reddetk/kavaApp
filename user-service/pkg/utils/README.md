# Utils Package

## Overview
Пакет `utils` содержит набор вспомогательных функций общего назначения, которые используются в различных частях приложения. Эти утилиты упрощают выполнение типовых задач и обеспечивают единообразный подход к решению распространенных проблем.

## Структура пакета

### string_utils.go
Утилиты для работы со строками:

#### Генерация случайных строк
- `GenerateRandomString(length int) (string, error)`: Генерирует криптографически безопасную случайную строку заданной длины
- `GenerateRandomAlphanumeric(length int) (string, error)`: Генерирует случайную буквенно-цифровую строку заданной длины

#### Проверка форматов
- `IsValidEmail(email string) bool`: Проверяет корректность email-адреса
- `IsValidPhone(phone string) bool`: Проверяет корректность телефонного номера

#### Форматирование строк
- `FormatPhone(phone string) string`: Форматирует телефонный номер в стандартный формат +X(XXX)XXX-XX-XX
- `Truncate(s string, maxLength int) string`: Обрезает строку до указанной длины и добавляет многоточие
- `Capitalize(s string) string`: Делает первую букву строки заглавной
- `MaskSensitiveData(data string, visibleChars int) string`: Маскирует конфиденциальные данные, оставляя видимыми только указанное количество символов в начале и конце

#### Преобразование форматов
- `CamelCaseToSnakeCase(s string) string`: Преобразует строку из camelCase в snake_case
- `SnakeCaseToCamelCase(s string) string`: Преобразует строку из snake_case в camelCase
- `RemoveWhitespace(s string) string`: Удаляет все пробельные символы из строки

### time_utils.go
Утилиты для работы с временем:

#### Форматирование времени
- `FormatDate(t time.Time) string`: Форматирует время в строку даты (YYYY-MM-DD)
- `FormatTime(t time.Time) string`: Форматирует время в строку (YYYY-MM-DD HH:MM:SS)
- `FormatISODate(t time.Time) string`: Форматирует время в строку в формате ISO 8601

#### Парсинг времени
- `ParseDate(dateStr string) (time.Time, error)`: Парсит строку в формате YYYY-MM-DD
- `ParseTime(timeStr string) (time.Time, error)`: Парсит строку в формате YYYY-MM-DD HH:MM:SS
- `ParseISODate(isoDateStr string) (time.Time, error)`: Парсит строку в формате ISO 8601

#### Расчет временных интервалов
- `DaysBetween(start, end time.Time) int`: Возвращает количество дней между датами
- `MonthsBetween(start, end time.Time) int`: Возвращает количество месяцев между датами
- `YearsBetween(start, end time.Time) int`: Возвращает количество лет между датами

#### Работа с границами периодов
- `StartOfDay(t time.Time) time.Time`: Возвращает время начала дня (00:00:00)
- `EndOfDay(t time.Time) time.Time`: Возвращает время конца дня (23:59:59)
- `StartOfMonth(t time.Time) time.Time`: Возвращает время начала месяца
- `EndOfMonth(t time.Time) time.Time`: Возвращает время конца месяца

#### Форматирование длительности и относительное время
- `FormatDuration(d time.Duration) string`: Форматирует duration в человекочитаемую строку
- `TimeAgo(t time.Time) string`: Возвращает строку, описывающую, сколько времени прошло

#### Работа с рабочими днями
- `IsWeekend(t time.Time) bool`: Проверяет, является ли дата выходным днем
- `AddBusinessDays(t time.Time, days int) time.Time`: Добавляет указанное количество рабочих дней

### math_utils.go
Математические утилиты:

#### Округление чисел
- `Round(value float64, precision int) float64`: Округляет число до указанной точности
- `RoundUp(value float64, precision int) float64`: Округляет число вверх
- `RoundDown(value float64, precision int) float64`: Округляет число вниз

#### Статистические расчеты
- `Mean(values []float64) float64`: Вычисляет среднее арифметическое
- `Median(values []float64) float64`: Вычисляет медиану
- `Mode(values []float64) float64`: Вычисляет моду (наиболее частое значение)
- `StandardDeviation(values []float64) float64`: Вычисляет стандартное отклонение
- `Variance(values []float64) float64`: Вычисляет дисперсию

#### Поиск экстремумов и процентили
- `Min(values []float64) float64`: Возвращает минимальное значение
- `Max(values []float64) float64`: Возвращает максимальное значение
- `Percentile(values []float64, percentile float64) float64`: Вычисляет процентиль
- `IQR(values []float64) float64`: Вычисляет межквартильный размах

#### Нормализация данных
- `Normalize(value, min, max float64) float64`: Нормализует значение к диапазону [0, 1]
- `NormalizeArray(values []float64) []float64`: Нормализует массив значений

#### Z-оценки и выбросы
- `ZScore(value, mean, stdDev float64) float64`: Вычисляет z-оценку
- `ZScoreArray(values []float64) []float64`: Вычисляет z-оценки для массива
- `IsOutlier(value float64, values []float64) bool`: Проверяет, является ли значение выбросом
- `RemoveOutliers(values []float64) []float64`: Удаляет выбросы из массива

### http_utils.go
Утилиты для работы с HTTP:

#### Создание и настройка HTTP-клиента
- `DefaultHTTPClientConfig() HTTPClientConfig`: Возвращает конфигурацию по умолчанию
- `NewHTTPClient(config HTTPClientConfig) *http.Client`: Создает новый HTTP-клиент

#### Выполнение HTTP-запросов
- `HTTPRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body []byte) HTTPResponse`: Выполняет HTTP-запрос
- `GetJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, target interface{}) error`: Выполняет GET-запрос и декодирует JSON
- `PostJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, requestBody, target interface{}) error`: Выполняет POST-запрос с JSON

#### Работа с URL и параметрами
- `BuildURL(baseURL string, queryParams map[string]string) (string, error)`: Строит URL с query-параметрами

#### Обработка ответов и заголовков
- `ParseJSONResponse(body []byte, target interface{}) error`: Парсит JSON-ответ
- `ExtractCookie(cookies []*http.Cookie, name string) (string, bool)`: Извлекает значение cookie

#### Проверка статус-кодов
- `IsSuccessStatusCode(statusCode int) bool`: Проверяет успешность статус-кода (2xx)
- `IsRedirectStatusCode(statusCode int) bool`: Проверяет, является ли код перенаправлением (3xx)
- `IsClientErrorStatusCode(statusCode int) bool`: Проверяет, является ли код ошибкой клиента (4xx)
- `IsServerErrorStatusCode(statusCode int) bool`: Проверяет, является ли код ошибкой сервера (5xx)
- `RetryableStatusCode(statusCode int) bool`: Проверяет, можно ли повторить запрос

#### Обработка ошибок
- `NewHTTPError(statusCode int, message string, body string) HTTPError`: Создает новую HTTP-ошибку

## Примеры использования

### Работа со строками
```go
// Генерация случайной строки
randomStr, err := utils.GenerateRandomString(10)
if err != nil {
    log.Fatalf("Failed to generate random string: %v", err)
}
fmt.Println("Random string:", randomStr)

// Проверка email
email := "user@example.com"
if utils.IsValidEmail(email) {
    fmt.Println("Email is valid")
} else {
    fmt.Println("Email is invalid")
}

// Форматирование телефонного номера
phone := "+375291234567"
formattedPhone := utils.FormatPhone(phone)
fmt.Println("Formatted phone:", formattedPhone) // +375(29)123-45-67
```

### Работа с временем
```go
// Форматирование даты
now := time.Now()
formattedDate := utils.FormatDate(now)
fmt.Println("Formatted date:", formattedDate)

// Расчет количества дней между датами
start := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
days := utils.DaysBetween(start, end)
fmt.Println("Days between:", days) // 14

// Получение начала и конца дня
startOfDay := utils.StartOfDay(now)
endOfDay := utils.EndOfDay(now)
fmt.Println("Start of day:", startOfDay)
fmt.Println("End of day:", endOfDay)
```

### Математические расчеты
```go
// Округление чисел
value := 3.14159
rounded := utils.Round(value, 2)
fmt.Println("Rounded:", rounded) // 3.14

// Статистические расчеты
values := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
mean := utils.Mean(values)
median := utils.Median(values)
stdDev := utils.StandardDeviation(values)
fmt.Printf("Mean: %.2f, Median: %.2f, StdDev: %.2f\n", mean, median, stdDev)

// Нормализация данных
normalized := utils.NormalizeArray(values)
fmt.Println("Normalized:", normalized)
```

### HTTP-запросы
```go
// Создание HTTP-клиента
client := utils.NewHTTPClient(utils.DefaultHTTPClientConfig())

// Выполнение GET-запроса
var response struct {
    Message string `json:"message"`
}
err := utils.GetJSON(context.Background(), client, "https://api.example.com/data", nil, &response)
if err != nil {
    log.Fatalf("Failed to fetch data: %v", err)
}
fmt.Println("Response:", response.Message)

// Построение URL с параметрами
url, err := utils.BuildURL("https://api.example.com/search", map[string]string{
    "query": "example",
    "limit": "10",
})
if err != nil {
    log.Fatalf("Failed to build URL: %v", err)
}
fmt.Println("URL:", url)
```

## Тестирование
Все функции в пакете `utils` покрыты модульными тестами, которые находятся в файле `utils_test.go`. Для запуска тестов используйте команду:

```bash
go test -v ./pkg/utils
```

## Принципы проектирования
1. **Повторное использование**: Функции спроектированы для использования в разных частях приложения
2. **Минимальные зависимости**: Пакет имеет минимальные внешние зависимости
3. **Тестируемость**: Код легко тестируется с помощью модульных тестов
4. **Документированность**: Все публичные функции имеют документацию
5. **Безопасность**: Функции учитывают аспекты безопасности (например, при генерации случайных строк)