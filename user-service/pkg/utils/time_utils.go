package utils

import (
	"fmt"
	"math"
	"time"
)

// StandardDateFormat - стандартный формат даты YYYY-MM-DD
const StandardDateFormat = "2006-01-02"

// StandardTimeFormat - стандартный формат времени YYYY-MM-DD HH:MM:SS
const StandardTimeFormat = "2006-01-02 15:04:05"

// ISODateFormat - формат даты в стандарте ISO 8601
const ISODateFormat = "2006-01-02T15:04:05Z07:00"

// FormatDate форматирует время в строку даты в формате YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format(StandardDateFormat)
}

// FormatTime форматирует время в строку в формате YYYY-MM-DD HH:MM:SS
func FormatTime(t time.Time) string {
	return t.Format(StandardTimeFormat)
}

// FormatISODate форматирует время в строку в формате ISO 8601
func FormatISODate(t time.Time) string {
	return t.Format(ISODateFormat)
}

// ParseDate парсит строку в формате YYYY-MM-DD в time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(StandardDateFormat, dateStr)
}

// ParseTime парсит строку в формате YYYY-MM-DD HH:MM:SS в time.Time
func ParseTime(timeStr string) (time.Time, error) {
	return time.Parse(StandardTimeFormat, timeStr)
}

// ParseISODate парсит строку в формате ISO 8601 в time.Time
func ParseISODate(isoDateStr string) (time.Time, error) {
	return time.Parse(ISODateFormat, isoDateStr)
}

// DaysBetween возвращает количество дней между двумя датами
func DaysBetween(start, end time.Time) int {
	// Нормализуем даты, убирая время
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

	return int(math.Round(endDate.Sub(startDate).Hours() / 24))
}

// MonthsBetween возвращает приблизительное количество месяцев между двумя датами
func MonthsBetween(start, end time.Time) int {
	// Количество месяцев = (год_конца - год_начала) * 12 + (месяц_конца - месяц_начала)
	months := (end.Year()-start.Year())*12 + int(end.Month()) - int(start.Month())

	// Корректировка, если конечный день меньше начального
	if end.Day() < start.Day() {
		months--
	}

	return months
}

// YearsBetween возвращает приблизительное количество лет между двумя датами
func YearsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()

	// Корректировка, если конечный месяц и день меньше начальных
	if end.Month() < start.Month() || (end.Month() == start.Month() && end.Day() < start.Day()) {
		years--
	}

	return years
}

// StartOfDay возвращает время, соответствующее началу дня (00:00:00)
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay возвращает время, соответствующее концу дня (23:59:59.999999999)
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, int(time.Second-1), t.Location())
}

// StartOfMonth возвращает время, соответствующее началу месяца
func StartOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth возвращает время, соответствующее концу месяца
func EndOfMonth(t time.Time) time.Time {
	// Начало следующего месяца минус 1 наносекунда
	return StartOfMonth(t.AddDate(0, 1, 0)).Add(-time.Nanosecond)
}

// FormatDuration форматирует duration в человекочитаемую строку
func FormatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(math.Mod(d.Hours(), 24))
	minutes := int(math.Mod(d.Minutes(), 60))
	seconds := int(math.Mod(d.Seconds(), 60))

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// IsWeekend проверяет, является ли дата выходным днем (суббота или воскресенье)
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// AddBusinessDays добавляет указанное количество рабочих дней к дате
func AddBusinessDays(t time.Time, days int) time.Time {
	result := t

	for i := 0; i < days; {
		result = result.AddDate(0, 0, 1)
		if !IsWeekend(result) {
			i++
		}
	}

	return result
}

// FormatDays converts time.Duration to float64 days value
func FormatDays(timeEstimate time.Duration) float64 {
	return timeEstimate.Hours() / 24
}

// TimeAgo возвращает строку, описывающую, сколько времени прошло с указанной даты
func TimeAgo(t time.Time) string {
	now := time.Now()
	duration := now.Sub(t)

	seconds := int(duration.Seconds())
	minutes := int(duration.Minutes())
	hours := int(duration.Hours())
	days := int(hours / 24)
	months := int(days / 30)
	years := int(days / 365)

	if years > 0 {
		if years == 1 {
			return "1 год назад"
		}
		return fmt.Sprintf("%d лет назад", years)
	} else if months > 0 {
		if months == 1 {
			return "1 месяц назад"
		}
		return fmt.Sprintf("%d месяцев назад", months)
	} else if days > 0 {
		if days == 1 {
			return "1 день назад"
		}
		return fmt.Sprintf("%d дней назад", days)
	} else if hours > 0 {
		if hours == 1 {
			return "1 час назад"
		}
		return fmt.Sprintf("%d часов назад", hours)
	} else if minutes > 0 {
		if minutes == 1 {
			return "1 минуту назад"
		}
		return fmt.Sprintf("%d минут назад", minutes)
	} else {
		if seconds <= 10 {
			return "только что"
		}
		return fmt.Sprintf("%d секунд назад", seconds)
	}
}
