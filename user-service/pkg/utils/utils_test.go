package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Тесты для string_utils.go

func TestGenerateRandomString(t *testing.T) {
	length := 10
	s1, err := GenerateRandomString(length)
	if err != nil {
		t.Fatalf("GenerateRandomString returned error: %v", err)
	}

	if len(s1) != length {
		t.Errorf("GenerateRandomString returned string of length %d, expected %d", len(s1), length)
	}

	s2, err := GenerateRandomString(length)
	if err != nil {
		t.Fatalf("GenerateRandomString returned error: %v", err)
	}

	if s1 == s2 {
		t.Errorf("GenerateRandomString returned same string twice: %s", s1)
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"test.email@example.co.uk", true},
		{"test+tag@example.com", true},
		{"test@localhost", true},
		{"test", false},
		{"test@", false},
		{"@example.com", false},
		{"test@example", true}, // Это валидно по нашему регулярному выражению
		{"", false},
	}

	for _, test := range tests {
		result := IsValidEmail(test.email)
		if result != test.expected {
			t.Errorf("IsValidEmail(%q) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestFormatPhone(t *testing.T) {
	tests := []struct {
		phone    string
		expected string
	}{
		{"3759261234567", "+375(926)123-45-67"},
		{"+3759261234567", "+375(926)123-45-67"},
		{"89261234567", "+8(926)123-45-67"},
		{"9261234567", "+7(926)123-45-67"},
		{"926-123-45-67", "+7(926)123-45-67"},
		{"+375 (926) 123-45-67", "+375(926)123-45-67"},
		{"123", "123"}, // Слишком короткий номер, возвращаем как есть
	}

	for _, test := range tests {
		result := FormatPhone(test.phone)
		if result != test.expected {
			t.Errorf("FormatPhone(%q) = %q, expected %q", test.phone, result, test.expected)
		}
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input     string
		maxLength int
		expected  string
	}{
		{"Hello, World!", 5, "He..."},
		{"Hello", 10, "Hello"},
		{"", 5, ""},
	}

	for _, test := range tests {
		result := Truncate(test.input, test.maxLength)
		if result != test.expected {
			t.Errorf("Truncate(%q, %d) = %q, expected %q", test.input, test.maxLength, result, test.expected)
		}
	}
}

// Тесты для time_utils.go

func TestFormatDate(t *testing.T) {
	date := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	expected := "2023-05-15"

	result := FormatDate(date)
	if result != expected {
		t.Errorf("FormatDate() = %q, expected %q", result, expected)
	}
}

func TestDaysBetween(t *testing.T) {
	start := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	end := time.Date(2023, 5, 20, 15, 45, 0, 0, time.UTC)
	expected := 5

	result := DaysBetween(start, end)
	if result != expected {
		t.Errorf("DaysBetween() = %d, expected %d", result, expected)
	}
}

func TestStartOfDay(t *testing.T) {
	date := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

	result := StartOfDay(date)
	if !result.Equal(expected) {
		t.Errorf("StartOfDay() = %v, expected %v", result, expected)
	}
}

func TestEndOfDay(t *testing.T) {
	date := time.Date(2023, 5, 15, 10, 30, 0, 0, time.UTC)
	expected := time.Date(2023, 5, 15, 23, 59, 59, int(time.Second-1), time.UTC)

	result := EndOfDay(date)
	if !result.Equal(expected) {
		t.Errorf("EndOfDay() = %v, expected %v", result, expected)
	}
}

// Тесты для math_utils.go

func TestRound(t *testing.T) {
	tests := []struct {
		value     float64
		precision int
		expected  float64
	}{
		{3.14159, 2, 3.14},
		{3.14159, 3, 3.142},
		{3.14159, 0, 3},
		{-3.14159, 2, -3.14},
	}

	for _, test := range tests {
		result := Round(test.value, test.precision)
		if result != test.expected {
			t.Errorf("Round(%f, %d) = %f, expected %f", test.value, test.precision, result, test.expected)
		}
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 2, 3, 4, 5}, 3},
		{[]float64{-1, 0, 1}, 0},
		{[]float64{}, 0},
	}

	for _, test := range tests {
		result := Mean(test.values)
		if result != test.expected {
			t.Errorf("Mean(%v) = %f, expected %f", test.values, result, test.expected)
		}
	}
}

func TestMedian(t *testing.T) {
	tests := []struct {
		values   []float64
		expected float64
	}{
		{[]float64{1, 3, 5, 7, 9}, 5},
		{[]float64{1, 3, 5, 7}, 4},
		{[]float64{}, 0},
	}

	for _, test := range tests {
		result := Median(test.values)
		if result != test.expected {
			t.Errorf("Median(%v) = %f, expected %f", test.values, result, test.expected)
		}
	}
}

func TestStandardDeviation(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	expected := 2.0

	result := StandardDeviation(values)
	if !almostEqual(result, expected, 0.01) {
		t.Errorf("StandardDeviation(%v) = %f, expected %f", values, result, expected)
	}
}

// Вспомогательная функция для сравнения float64 с погрешностью
func almostEqual(a, b, tolerance float64) bool {
	return (a-b) < tolerance && (b-a) < tolerance
}

// Тесты для http_utils.go

func TestNewHTTPClient(t *testing.T) {
	config := DefaultHTTPClientConfig()
	client := NewHTTPClient(config)

	if client == nil {
		t.Error("NewHTTPClient returned nil")
	}

	if client.Timeout != config.Timeout {
		t.Errorf("NewHTTPClient returned client with timeout %v, expected %v", client.Timeout, config.Timeout)
	}
}

func TestHTTPRequest(t *testing.T) {
	// Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		// Проверяем заголовки
		if r.Header.Get("X-Test-Header") != "test-value" {
			t.Errorf("Expected header X-Test-Header: test-value, got %s", r.Header.Get("X-Test-Header"))
		}

		// Отправляем ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer server.Close()

	// Создаем клиент
	client := NewHTTPClient(DefaultHTTPClientConfig())

	// Выполняем запрос
	headers := map[string]string{
		"X-Test-Header": "test-value",
	}
	response := HTTPRequest(context.Background(), client, http.MethodGet, server.URL, headers, nil)

	// Проверяем результат
	if response.Error != nil {
		t.Fatalf("HTTPRequest returned error: %v", response.Error)
	}

	if response.StatusCode != http.StatusOK {
		t.Errorf("HTTPRequest returned status code %d, expected %d", response.StatusCode, http.StatusOK)
	}

	if string(response.Body) != `{"message":"success"}` {
		t.Errorf("HTTPRequest returned body %q, expected %q", string(response.Body), `{"message":"success"}`)
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		baseURL     string
		queryParams map[string]string
		expected    string
	}{
		{
			"https://example.com",
			map[string]string{"param1": "value1", "param2": "value2"},
			"https://example.com?param1=value1&param2=value2",
		},
		{
			"https://example.com?existing=param",
			map[string]string{"param1": "value1"},
			"https://example.com?existing=param&param1=value1",
		},
		{
			"https://example.com",
			map[string]string{},
			"https://example.com",
		},
	}

	for _, test := range tests {
		result, err := BuildURL(test.baseURL, test.queryParams)
		if err != nil {
			t.Fatalf("BuildURL returned error: %v", err)
		}

		if result != test.expected {
			t.Errorf("BuildURL(%q, %v) = %q, expected %q", test.baseURL, test.queryParams, result, test.expected)
		}
	}
}
