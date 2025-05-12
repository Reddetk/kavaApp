// internal/infrastructure/services/discounted_clv_service.go
package services

import (
	"errors"
	"fmt"
	"math"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

// DiscountedCLVService implements the CLVService interface
type DiscountedCLVService struct {
	// Ставка дисконтирования (например, 0.1 для 10%)
	discountRate float64
	// Период прогнозирования в месяцах
	forecastPeriod int
	// Хранилище исторических данных CLV
	clvHistory map[uuid.UUID][]entities.CLVDataPoint
}

// NewDiscountedCLVService creates a new instance of DiscountedCLVService
func NewDiscountedCLVService() services.CLVService {
	return &DiscountedCLVService{
		discountRate:   0.1,
		forecastPeriod: 12,
		clvHistory:     make(map[uuid.UUID][]entities.CLVDataPoint),
	}
}

// CalculateCLV calculates the Customer Lifetime Value using the discounted cash flow method
func (s *DiscountedCLVService) CalculateCLV(userID uuid.UUID, retentionRate float64, avgCheck float64) (float64, error) {
	if retentionRate <= 0 || retentionRate > 1 {
		return 0, errors.New("retention rate must be between 0 and 1")
	}

	if avgCheck <= 0 {
		return 0, errors.New("average check must be positive")
	}

	// Расчет CLV по формуле дисконтированного денежного потока
	var clv float64
	for month := 1; month <= s.forecastPeriod; month++ {
		// Вероятность того, что клиент все еще активен в данном месяце
		survivalProbability := math.Pow(retentionRate, float64(month))

		// Ожидаемый доход в данном месяце
		expectedRevenue := avgCheck * survivalProbability

		// Дисконтирование
		discountFactor := 1.0 / math.Pow(1+s.discountRate/12, float64(month))

		// Добавление к общему CLV
		clv += expectedRevenue * discountFactor
	}

	// Сохраняем историческую точку данных
	s.saveHistoricalDataPoint(userID, clv, "default")

	return clv, nil
}

// UpdateCLV updates the CLV for a specific user
func (s *DiscountedCLVService) UpdateCLV(userID uuid.UUID) (float64, error) {
	// В реальной реализации здесь будет логика получения данных пользователя
	// и вызов CalculateCLV с актуальными параметрами

	// Для демонстрации используем фиктивные значения
	retentionRate := 0.8 // 80% вероятность удержания
	avgCheck := 100.0    // Средний чек 100 единиц

	return s.CalculateCLV(userID, retentionRate, avgCheck)
}

// EstimateCLV estimates the CLV for a specific scenario
func (s *DiscountedCLVService) EstimateCLV(userID uuid.UUID, scenario string) (float64, error) {
	var retentionRate, avgCheck float64

	switch scenario {
	case "default":
		retentionRate = 0.8
		avgCheck = 100.0
	case "optimistic":
		retentionRate = 0.9
		avgCheck = 120.0
	case "pessimistic":
		retentionRate = 0.7
		avgCheck = 80.0
	default:
		return 0, fmt.Errorf("unknown scenario: %s", scenario)
	}

	clv, err := s.CalculateCLV(userID, retentionRate, avgCheck)
	if err != nil {
		return 0, err
	}

	// Сохраняем историческую точку данных для конкретного сценария
	s.saveHistoricalDataPoint(userID, clv, scenario)

	return clv, nil
}

// GetHistoricalCLV returns historical CLV data points for a user
func (s *DiscountedCLVService) GetHistoricalCLV(userID uuid.UUID, periodMonths int) ([]entities.CLVDataPoint, error) {
	history, exists := s.clvHistory[userID]
	if !exists {
		return []entities.CLVDataPoint{}, nil
	}

	// Фильтруем точки данных по периоду
	if periodMonths <= 0 {
		return history, nil
	}

	cutoffDate := time.Now().AddDate(0, -periodMonths, 0)
	var filteredHistory []entities.CLVDataPoint

	for _, point := range history {
		if point.Date.After(cutoffDate) {
			filteredHistory = append(filteredHistory, point)
		}
	}

	return filteredHistory, nil
}

// SetDiscountRate sets the discount rate
func (s *DiscountedCLVService) SetDiscountRate(rate float64) error {
	if rate < 0 || rate > 1 {
		return errors.New("discount rate must be between 0 and 1")
	}
	s.discountRate = rate
	return nil
}

// SetForecastPeriod sets the forecast period in months
func (s *DiscountedCLVService) SetForecastPeriod(months int) error {
	if months <= 0 {
		return errors.New("forecast period must be positive")
	}
	s.forecastPeriod = months
	return nil
}

// saveHistoricalDataPoint saves a historical CLV data point
func (s *DiscountedCLVService) saveHistoricalDataPoint(userID uuid.UUID, value float64, scenario string) {
	dataPoint := entities.CLVDataPoint{
		UserID:    userID,
		Date:      time.Now(),
		Value:     value,
		Scenario:  scenario,
		Algorithm: "DiscountedCashFlow",
	}

	if _, exists := s.clvHistory[userID]; !exists {
		s.clvHistory[userID] = []entities.CLVDataPoint{}
	}

	s.clvHistory[userID] = append(s.clvHistory[userID], dataPoint)
}
