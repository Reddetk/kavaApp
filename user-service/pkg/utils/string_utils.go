package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// EmailRegex - регулярное выражение для проверки email
var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// PhoneRegex - регулярное выражение для проверки телефонного номера
// Поддерживает форматы: +7(123)456-78-90, 8(123)456-78-90, +7 123 456 78 90 и т.д.
var PhoneRegex = regexp.MustCompile(`^(\+?\d{1,3}[ -]?)?\(?\d{3}\)?[ -]?\d{3}[ -]?\d{2}[ -]?\d{2}$`)

// GenerateRandomString генерирует случайную строку заданной длины
func GenerateRandomString(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}

// GenerateRandomAlphanumeric генерирует случайную буквенно-цифровую строку заданной длины
func GenerateRandomAlphanumeric(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	for i := range buffer {
		buffer[i] = charset[int(buffer[i])%len(charset)]
	}

	return string(buffer), nil
}

// IsValidEmail проверяет, является ли строка корректным email-адресом
func IsValidEmail(email string) bool {
	// Allow local domain names (without TLD) for development/testing environments
	var EmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+(\.[a-zA-Z]{2,}|$)`)
	return EmailRegex.MatchString(email)
}

// IsValidPhone проверяет, является ли строка корректным телефонным номером
func IsValidPhone(phone string) bool {
	return PhoneRegex.MatchString(phone)
}

// FormatPhone форматирует телефонный номер в стандартный формат +375(XXX)XXX-XX-XX
func FormatPhone(phone string) string {
	// Удаляем все нецифровые символы
	digits := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, phone)

	// Проверяем длину и форматируем
	length := len(digits)
	if length < 9 || length > 12 {
		return phone // Возвращаем исходный номер, если он некорректной длины
	}

	// Обрабатываем разные форматы номеров
	var countryCode, areaCode, prefix, lineNumber string

	// Определяем код страны и остальные части номера
	if length == 12 && strings.HasPrefix(digits, "375") {
		// Международный формат с кодом Беларуси: 375XXXXXXXXX
		countryCode = "375"
		areaCode = digits[3:6]
		prefix = digits[6:9]
		lineNumber = digits[9:12]
	} else if length == 11 && digits[0] == '8' {
		// Формат с 8: 89XXXXXXXXX (для России)
		countryCode = "8"
		areaCode = digits[1:4]
		prefix = digits[4:7]
		lineNumber = digits[7:11]
	} else if length == 11 && strings.HasPrefix(digits, "375") {
		// Сокращенный международный формат: 375XXXXXXXX
		countryCode = "375"
		areaCode = digits[3:5]
		prefix = digits[5:8]
		lineNumber = digits[8:11]
	} else if length == 10 {
		// Формат с оператором: 9XXXXXXXXX (для России)
		countryCode = "7" // Используем 7 по умолчанию для России
		areaCode = digits[0:3]
		prefix = digits[3:6]
		lineNumber = digits[6:10]
	} else if length == 9 {
		// Формат без кода страны: XXXXXXXXX (для Беларуси)
		countryCode = "375"
		areaCode = digits[0:2]
		prefix = digits[2:5]
		lineNumber = digits[5:9]
	} else {
		// Для других случаев пытаемся определить формат
		if strings.HasPrefix(digits, "375") {
			countryCode = "375"
			restDigits := digits[3:]
			if len(restDigits) >= 9 {
				areaCode = restDigits[0:3]
				prefix = restDigits[3:6]
				lineNumber = restDigits[6:9]
			} else {
				return fmt.Sprintf("+%s(%s)%s-%s-%s", countryCode, areaCode, prefix, lineNumber[:2], lineNumber[2:])
			}
		} else {
			return fmt.Sprintf("+%s(%s)%s-%s-%s", countryCode, areaCode, prefix, lineNumber[:2], lineNumber[2:])
		}
	}

	return fmt.Sprintf("+%s(%s)%s-%s-%s", countryCode, areaCode, prefix, lineNumber[:2], lineNumber[2:])
}

// Truncate обрезает строку до указанной длины и добавляет многоточие, если строка была обрезана
func Truncate(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Capitalize делает первую букву строки заглавной
func Capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// CamelCaseToSnakeCase преобразует строку из camelCase в snake_case
func CamelCaseToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// SnakeCaseToCamelCase преобразует строку из snake_case в camelCase
func SnakeCaseToCamelCase(s string) string {
	var result strings.Builder
	nextUpper := false

	for _, r := range s {
		if r == '_' {
			nextUpper = true
		} else {
			if nextUpper {
				result.WriteRune(unicode.ToUpper(r))
				nextUpper = false
			} else {
				result.WriteRune(r)
			}
		}
	}

	return result.String()
}

// RemoveWhitespace удаляет все пробельные символы из строки
func RemoveWhitespace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

// MaskSensitiveData маскирует конфиденциальные данные (например, номер карты)
func MaskSensitiveData(data string, visibleChars int) string {
	length := len(data)
	if length <= visibleChars*2 {
		return data
	}

	prefix := data[:visibleChars]
	suffix := data[length-visibleChars:]
	mask := strings.Repeat("*", length-visibleChars*2)

	return prefix + mask + suffix
}
