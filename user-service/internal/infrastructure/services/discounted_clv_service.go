// internal/infrastructure/services/discounted_clv_service.go
package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
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
	// Репозиторий пользователей
	userRepository repositories.UserRepository
	// Репозиторий заказов
	orderRepository repositories.TransactionRepository
}

// NewDiscountedCLVService creates a new instance of DiscountedCLVService
func NewDiscountedCLVService(userRepo repositories.UserRepository, orderRepo repositories.TransactionRepository) services.CLVService {
	return &DiscountedCLVService{
		discountRate:    0.1,
		forecastPeriod:  12,
		clvHistory:      make(map[uuid.UUID][]entities.CLVDataPoint),
		userRepository:  userRepo,
		orderRepository: orderRepo,
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
	ctx := context.Background()

	// Получаем данные пользователя из репозитория
	user, err := s.userRepository.Get(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user data: %w", err)
	}

	// Получаем историю транзакций пользователя
	transactions, err := s.orderRepository.GetByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user transactions: %w", err)
	}

	// Рассчитываем коэффициент удержания на основе истории транзакций
	// Преобразуем []*entities.Segment в []entities.Segment
	transactionsVal := make([]entities.Transaction, len(transactions))
	for i, s := range transactions {
		if s != nil {
			transactionsVal[i] = *s
		}
	}
	retentionRate := s.calculateRetentionRate(transactionsVal)

	// Рассчитываем средний чек на основе истории транзакций
	avgCheck := s.calculateAverageCheck(transactionsVal)

	// Если у пользователя нет истории транзакций, используем значения по умолчанию
	if len(transactions) == 0 {
		retentionRate = 0.5 // Значение по умолчанию для новых пользователей
		avgCheck = 50.0     // Значение по умолчанию для новых пользователей
	}

	// Применяем модификаторы на основе данных пользователя
	retentionRate = s.applyUserFactors(*user, retentionRate)

	// Вызываем метод расчета CLV с актуальными параметрами
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

// calculateRetentionRate вычисляет коэффициент удержания на основе истории заказов
func (s *DiscountedCLVService) calculateRetentionRate(orders []entities.Transaction) float64 {
	if len(orders) < 2 {
		return 0.5 // Значение по умолчанию, если недостаточно данных
	}

	// Сортируем заказы по дате (предполагается, что в entities.Order есть поле CreatedAt)
	// sort.Slice(orders, func(i, j int) bool {
	//     return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	// })

	// Вычисляем среднее время между заказами
	var totalGaps int
	var totalDays float64

	for i := 1; i < len(orders); i++ {
		// Разница в днях между последовательными заказами
		// daysBetween := orders[i].CreatedAt.Sub(orders[i-1].CreatedAt).Hours() / 24
		daysBetween := 30.0 // Временно используем фиксированное значение

		if daysBetween <= 90 { // Учитываем только заказы в пределах 90 дней
			totalGaps++
			totalDays += daysBetween
		}
	}

	if totalGaps == 0 {
		return 0.4 // Низкий коэффициент удержания, если большие промежутки между заказами
	}

	// Вычисляем средний промежуток между заказами в днях
	avgGap := totalDays / float64(totalGaps)

	// Преобразуем в коэффициент удержания (чем меньше промежуток, тем выше удержание)
	// Формула: 1 - (avgGap / 90), ограниченная диапазоном [0.3, 0.95]
	retentionRate := 1 - (avgGap / 90)

	// Ограничиваем значение
	if retentionRate < 0.3 {
		retentionRate = 0.3
	} else if retentionRate > 0.95 {
		retentionRate = 0.95
	}

	return retentionRate
}

// calculateAverageCheck вычисляет средний чек на основе истории заказов
func (s *DiscountedCLVService) calculateAverageCheck(orders []entities.Transaction) float64 {
	if len(orders) == 0 {
		return 50.0 // Значение по умолчанию, если нет заказов
	}

	var totalAmount float64
	for _, order := range orders {
		totalAmount += order.Amount
		//totalAmount += 100.0 // Временно используем фиксированное значение
	}

	return totalAmount / float64(len(orders))
}

// applyUserFactors корректирует коэффициент удержания на основе данных пользователя
func (s *DiscountedCLVService) applyUserFactors(user entities.User, baseRate float64) float64 {
	adjustedRate := baseRate

	// Пример: корректировка на основе возраста аккаунта
	// accountAge := time.Since(user.CreatedAt).Hours() / 24 / 365 // в годах
	accountAge := 1.0 // Временно используем фиксированное значение

	if accountAge < 1 {
		// Новые пользователи имеют более низкий коэффициент удержания
		adjustedRate *= 0.9
	} else if accountAge > 3 {
		// Давние пользователи имеют более высокий коэффициент удержания
		adjustedRate *= 1.1
	}

	// Пример: корректировка на основе активности в программе лояльности
	// if user.LoyaltyLevel > 0 {
	//     adjustedRate *= (1 + float64(user.LoyaltyLevel) * 0.05)
	// }

	// Ограничиваем итоговое значение
	if adjustedRate < 0.3 {
		adjustedRate = 0.3
	} else if adjustedRate > 0.95 {
		adjustedRate = 0.95
	}

	return adjustedRate
}
