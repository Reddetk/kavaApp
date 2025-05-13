package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPClientConfig содержит настройки для HTTP-клиента
type HTTPClientConfig struct {
	Timeout               time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration
	KeepAlive             time.Duration
	DisableCompression    bool
	DisableKeepAlives     bool
}

// DefaultHTTPClientConfig возвращает конфигурацию HTTP-клиента по умолчанию
func DefaultHTTPClientConfig() HTTPClientConfig {
	return HTTPClientConfig{
		Timeout:               30 * time.Second,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		MaxConnsPerHost:       100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		KeepAlive:             30 * time.Second,
		DisableCompression:    false,
		DisableKeepAlives:     false,
	}
}

// NewHTTPClient создает новый HTTP-клиент с указанной конфигурацией
func NewHTTPClient(config HTTPClientConfig) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		MaxConnsPerHost:       config.MaxConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableCompression:    config.DisableCompression,
		DisableKeepAlives:     config.DisableKeepAlives,
	}
	
	return &http.Client{
		Timeout:   config.Timeout,
		Transport: transport,
	}
}

// HTTPResponse представляет ответ HTTP-запроса
type HTTPResponse struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Error      error
}

// HTTPRequest выполняет HTTP-запрос с указанными параметрами
func HTTPRequest(ctx context.Context, client *http.Client, method, url string, headers map[string]string, body []byte) HTTPResponse {
	var response HTTPResponse
	
	// Создаем запрос
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	if err != nil {
		response.Error = fmt.Errorf("error creating request: %w", err)
		return response
	}
	
	// Добавляем заголовки
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		response.Error = fmt.Errorf("error executing request: %w", err)
		return response
	}
	defer resp.Body.Close()
	
	// Читаем тело ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Error = fmt.Errorf("error reading response body: %w", err)
		return response
	}
	
	// Заполняем структуру ответа
	response.StatusCode = resp.StatusCode
	response.Headers = resp.Header
	response.Body = responseBody
	
	return response
}

// GetJSON выполняет GET-запрос и декодирует JSON-ответ в указанную структуру
func GetJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, target interface{}) error {
	response := HTTPRequest(ctx, client, http.MethodGet, url, headers, nil)
	if response.Error != nil {
		return response.Error
	}
	
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, string(response.Body))
	}
	
	return json.Unmarshal(response.Body, target)
}

// PostJSON выполняет POST-запрос с JSON-телом и декодирует JSON-ответ в указанную структуру
func PostJSON(ctx context.Context, client *http.Client, url string, headers map[string]string, requestBody interface{}, target interface{}) error {
	// Кодируем тело запроса в JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %w", err)
	}
	
	// Добавляем заголовок Content-Type, если он не указан
	if headers == nil {
		headers = make(map[string]string)
	}
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}
	
	response := HTTPRequest(ctx, client, http.MethodPost, url, headers, jsonBody)
	if response.Error != nil {
		return response.Error
	}
	
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d, body: %s", response.StatusCode, string(response.Body))
	}
	
	if target != nil {
		return json.Unmarshal(response.Body, target)
	}
	
	return nil
}

// BuildURL строит URL с query-параметрами
func BuildURL(baseURL string, queryParams map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("error parsing URL: %w", err)
	}
	
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// ExtractCookie извлекает значение cookie по имени
func ExtractCookie(cookies []*http.Cookie, name string) (string, bool) {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie.Value, true
		}
	}
	return "", false
}

// ParseJSONResponse парсит JSON-ответ в указанную структуру
func ParseJSONResponse(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

// IsSuccessStatusCode проверяет, является ли статус-код успешным (2xx)
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// IsRedirectStatusCode проверяет, является ли статус-код перенаправлением (3xx)
func IsRedirectStatusCode(statusCode int) bool {
	return statusCode >= 300 && statusCode < 400
}

// IsClientErrorStatusCode проверяет, является ли статус-код ошибкой клиента (4xx)
func IsClientErrorStatusCode(statusCode int) bool {
	return statusCode >= 400 && statusCode < 500
}

// IsServerErrorStatusCode проверяет, является ли статус-код ошибкой сервера (5xx)
func IsServerErrorStatusCode(statusCode int) bool {
	return statusCode >= 500 && statusCode < 600
}

// RetryableStatusCode проверяет, можно ли повторить запрос при получении данного статус-кода
func RetryableStatusCode(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable ||
		statusCode == http.StatusGatewayTimeout
}

// HTTPError представляет ошибку HTTP-запроса
type HTTPError struct {
	StatusCode int
	Message    string
	Body       string
}

// Error реализует интерфейс error
func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP error %d: %s", e.StatusCode, e.Message)
}

// NewHTTPError создает новую ошибку HTTP
func NewHTTPError(statusCode int, message string, body string) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	}
}