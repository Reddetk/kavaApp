// internal/domain/services/regression_service.go
package services

import (
	"errors"
	"fmt"
	"time"

	"analitics-service/internal/domain/enteties"

	"github.com/sajari/regression"
)

// Ошибки сервиса регрессионного анализа
var (
	ErrInsufficientData = errors.New("insufficient data for regression analysis")
	ErrInvalidParameter = errors.New("invalid parameter for regression analysis")
	ErrRegressionFailed = errors.New("regression analysis failed")
)

// RegressionService определяет интерфейс для сервиса регрессионного анализа
type RegressionService interface {
	// AnalyzeDiscountEffect анализирует влияние скидок на продажи
	AnalyzeDiscountEffect(productID string, period time.Duration) (*entities.DiscountEffect, error)

	// AnalyzeDiscountEffectByCategory анализирует влияние скидок на продажи по категории товаров
	AnalyzeDiscountEffectByCategory(category string, period time.Duration) (*entities.DiscountEffect, error)

	// GenerateDiscountRecommendations генерирует рекомендации по оптимальным скидкам
	GenerateDiscountRecommendations() ([]*entities.DiscountRecommendation, error)

	// AnalyzeABTestResults анализирует результаты A/B тестов для оптимизации скидок
	AnalyzeABTestResults(testIDs []string) (*entities.ABTestAnalysis, error)
}

// regressionServiceImpl реализация сервиса регрессионного анализа
type regressionServiceImpl struct {
	transactionRepo entities.TransactionRepository
	productRepo     entities.ProductRepository
	abTestRepo      entities.ABTestRepository
}

// NewRegressionService создает новый экземпляр сервиса регрессионного анализа
func NewRegressionService(
	transactionRepo entities.TransactionRepository,
	productRepo entities.ProductRepository,
	abTestRepo entities.ABTestRepository,
) RegressionService {
	return &regressionServiceImpl{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
		abTestRepo:      abTestRepo,
	}
}

// AnalyzeDiscountEffect анализирует влияние скидок на продажи
func (s *regressionServiceImpl) AnalyzeDiscountEffect(productID string, period time.Duration) (*entities.DiscountEffect, error) {
	// Получаем транзакции за указанный период
	endDate := time.Now()
	startDate := endDate.Add(-period)

	// Получаем все транзакции, содержащие данный товар
	transactions, err := s.transactionRepo.FindByProductAndDateRange(productID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %w", err)
	}

	if len(transactions) < 30 {
		return nil, ErrInsufficientData
	}

	// Подготавливаем данные для регрессионного анализа
	r := new(regression.Regression)
	r.SetObserved("Sales")
	r.SetVar(0, "Discount")    // Скидка как процент
	r.SetVar(1, "WeekDay")     // День недели (числовое представление)
	r.SetVar(2, "MonthPeriod") // Период месяца (1-3 недели)
	r.SetVar(3, "IsHoliday")   // Праздничный день (0/1)

	// Агрегируем данные по дням для анализа
	dailyData := s.aggregateTransactionsByDay(transactions, productID)

	// Добавляем данные в регрессионную модель
	for _, data := range dailyData {
		r.Train(
			regression.DataPoint(data.Sales,
				[]float64{
					data.AvgDiscount,
					float64(data.Date.Weekday()),
					float64((data.Date.Day()-1)/7 + 1),
					boolToFloat(data.IsHoliday),
				},
			),
		)
	}

	// Запускаем регрессионный анализ
	err = r.Run()
	if err != nil {
		return nil, fmt.Errorf("regression analysis failed: %w", err)
	}

	// Анализируем коэффициенты для определения влияния скидки на продажи
	// Коэффициент при переменной "Discount" - это и есть Lift Factor
	discountCoeff := r.Coeff(0)

	// Оцениваем оптимальный уровень скидки на основе модели
	// Максимизируем доход: Revenue = Price * (1 - Discount) * Sales, где Sales зависит от Discount
	optimalDiscount := s.findOptimalDiscount(r, productID)

	// Формируем результат анализа
	result := &entities.DiscountEffect{
		ProductID:         productID,
		LiftFactor:        discountCoeff,
		RSqaured:          r.R2,
		OptimalDiscount:   optimalDiscount,
		AnalysisTimestamp: time.Now(),
		Coefficients: map[string]float64{
			"Discount":    discountCoeff,
			"WeekDay":     r.Coeff(1),
			"MonthPeriod": r.Coeff(2),
			"IsHoliday":   r.Coeff(3),
			"Intercept":   r.Coeff(4),
		},
		DataPointsCount: len(dailyData),
		PeriodStart:     startDate,
		PeriodEnd:       endDate,
	}

	return result, nil
}

// boolToFloat конвертирует булево значение в float64
func boolToFloat(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// AnalyzeDiscountEffectByCategory анализирует влияние скидок на продажи по категории товаров
func (s *regressionServiceImpl) AnalyzeDiscountEffectByCategory(category string, period time.Duration) (*entities.DiscountEffect, error) {
	// Получаем транзакции за указанный период для категории
	endDate := time.Now()
	startDate := endDate.Add(-period)

	// Получаем все транзакции, содержащие товары из этой категории
	transactions, err := s.transactionRepo.FindByCategoryAndDateRange(category, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve transactions: %w", err)
	}

	if len(transactions) < 30 {
		return nil, ErrInsufficientData
	}

	// Подготавливаем данные для регрессионного анализа
	r := new(regression.Regression)
	r.SetObserved("CategorySales")
	r.SetVar(0, "AvgDiscount")  // Средняя скидка по категории
	r.SetVar(1, "WeekDay")      // День недели
	r.SetVar(2, "MonthPeriod")  // Период месяца
	r.SetVar(3, "IsHoliday")    // Праздничный день
	r.SetVar(4, "ProductCount") // Количество уникальных товаров в категории

	// Агрегируем данные по дням для анализа
	dailyData := s.aggregateCategoryTransactionsByDay(transactions, category)

	// Добавляем данные в регрессионную модель
	for _, data := range dailyData {
		r.Train(
			regression.DataPoint(data.Sales,
				[]float64{
					data.AvgDiscount,
					float64(data.Date.Weekday()),
					float64((data.Date.Day()-1)/7 + 1),
					boolToFloat(data.IsHoliday),
					float64(data.ProductCount),
				},
			),
		)
	}

	// Запускаем регрессионный анализ
	err = r.Run()
	if err != nil {
		return nil, fmt.Errorf("regression analysis failed: %w", err)
	}

	// Анализируем коэффициенты для определения влияния скидки на продажи
	discountCoeff := r.Coeff(0)

	// Оцениваем оптимальный уровень скидки на основе модели
	optimalDiscount := s.findOptimalCategoryDiscount(r, category)

	// Формируем результат анализа
	result := &entities.DiscountEffect{
		Category:          category,
		LiftFactor:        discountCoeff,
		RSqaured:          r.R2,
		OptimalDiscount:   optimalDiscount,
		AnalysisTimestamp: time.Now(),
		Coefficients: map[string]float64{
			"AvgDiscount":  discountCoeff,
			"WeekDay":      r.Coeff(1),
			"MonthPeriod":  r.Coeff(2),
			"IsHoliday":    r.Coeff(3),
			"ProductCount": r.Coeff(4),
			"Intercept":    r.Coeff(5),
		},
		DataPointsCount: len(dailyData),
		PeriodStart:     startDate,
		PeriodEnd:       endDate,
	}

	return result, nil
}

// GenerateDiscountRecommendations генерирует рекомендации по оптимальным скидкам
func (s *regressionServiceImpl) GenerateDiscountRecommendations() ([]*entities.DiscountRecommendation, error) {
	// Получаем все категории товаров
	categories, err := s.productRepo.GetAllCategories()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve categories: %w", err)
	}

	recommendations := make([]*entities.DiscountRecommendation, 0, len(categories))

	// Для каждой категории проводим анализ и генерируем рекомендации
	for _, category := range categories {
		// Анализируем влияние скидок за последние 3 месяца
		effect, err := s.AnalyzeDiscountEffectByCategory(category, 90*24*time.Hour)
		// Если недостаточно данных, пропускаем категорию
		if errors.Is(err, ErrInsufficientData) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to analyze category %s: %w", category, err)
		}

		// Получаем прибыльность категории
		abcCategory, err := s.productRepo.GetCategoryABCClassification(category)
		if err != nil {
			return nil, fmt.Errorf("failed to get ABC classification for category %s: %w", category, err)
		}

		// Формируем рекомендацию на основе анализа и категории ABC
		recommendation := &entities.DiscountRecommendation{
			Category:        category,
			OptimalDiscount: effect.OptimalDiscount,
			LiftFactor:      effect.LiftFactor,
			ABCCategory:     abcCategory,
			Confidence:      effect.RSqaured, // Используем R² как меру уверенности в рекомендации
			GeneratedAt:     time.Now(),
		}

		// Корректируем рекомендацию в зависимости от категории ABC
		s.adjustRecommendationByABCCategory(recommendation)

		recommendations = append(recommendations, recommendation)
	}

	return recommendations, nil
}

// AnalyzeABTestResults анализирует результаты A/B тестов для оптимизации скидок
func (s *regressionServiceImpl) AnalyzeABTestResults(testIDs []string) (*entities.ABTestAnalysis, error) {
	tests := make([]*entities.ABTestResult, 0, len(testIDs))

	// Получаем результаты всех указанных тестов
	for _, testID := range testIDs {
		test, err := s.abTestRepo.GetByID(testID)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve test %s: %w", testID, err)
		}
		tests = append(tests, test)
	}

	if len(tests) < 3 {
		return nil, ErrInsufficientData
	}

	// Проводим регрессионный анализ для определения зависимости между размером скидки и Lift-фактором
	r := new(regression.Regression)
	r.SetObserved("Lift")
	r.SetVar(0, "DiscountPct")
	r.SetVar(1, "BasePrice")
	r.SetVar(2, "TestDuration") // Длительность теста в днях

	// Преобразуем данные тестов для регрессионного анализа
	for _, test := range tests {
		// Проверяем, что в тесте использовалась скидка
		if !test.TestGroup.CouponUsed || test.TestGroup.DiscountPct <= 0 {
			continue
		}

		// Вычисляем длительность теста в днях
		duration := test.EndDate.Sub(test.StartDate).Hours() / 24

		// Получаем среднюю базовую цену тестируемых товаров
		basePrice, err := s.getAverageBasePriceForTest(test)
		if err != nil {
			continue
		}

		r.Train(
			regression.DataPoint(test.Lift,
				[]float64{
					test.TestGroup.DiscountPct,
					basePrice,
					duration,
				},
			),
		)
	}

	// Запускаем регрессионный анализ
	err := r.Run()
	if err != nil {
		return nil, fmt.Errorf("regression analysis failed: %w", err)
	}

	// Оцениваем оптимальный уровень скидки на основе модели
	// Предполагаем среднюю цену и среднюю длительность теста
	avgBasePrice, _ := s.getAverageBasePrice()
	optimalDiscount := s.findOptimalDiscountFromABTests(r, avgBasePrice, 14.0) // 14 дней стандартная длительность теста

	// Формируем результат анализа
	result := &entities.ABTestAnalysis{
		TestsAnalyzed:     len(tests),
		DiscountCoeff:     r.Coeff(0),
		BasePriceCoeff:    r.Coeff(1),
		DurationCoeff:     r.Coeff(2),
		InterceptCoeff:    r.Coeff(3),
		RSqaured:          r.R2,
		OptimalDiscount:   optimalDiscount,
		AnalysisTimestamp: time.Now(),
		Recommendations: []string{
			fmt.Sprintf("Оптимальная скидка для будущих тестов: %.2f%%", optimalDiscount*100),
			fmt.Sprintf("Зависимость Lift от скидки: %.3f", r.Coeff(0)),
		},
	}

	// Добавляем дополнительные рекомендации на основе результатов анализа
	if r.Coeff(0) < 0 {
		result.Recommendations = append(result.Recommendations,
			"Обнаружен отрицательный эффект скидок на продажи, рекомендуется пересмотреть стратегию ценообразования")
	}

	if r.Coeff(1) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Товары с высокой базовой ценой показывают более высокий Lift-фактор при скидках")
	}

	if r.Coeff(2) > 0 {
		result.Recommendations = append(result.Recommendations,
			"Более длительные тесты показывают более высокий Lift-фактор, рекомендуется увеличить длительность будущих тестов")
	}

	return result, nil
}

// Вспомогательные методы

// aggregateTransactionsByDay агрегирует транзакции по дням
func (s *regressionServiceImpl) aggregateTransactionsByDay(transactions []*entities.Transaction, productID string) []*entities.DailyTransactionData {
	// Мапа для агрегации данных по дням
	dailyMap := make(map[string]*entities.DailyTransactionData)

	// Определяем праздничные дни (упрощенная реализация)
	holidays := s.getHolidayDates()

	// Агрегируем данные по дням
	for _, tx := range transactions {
		dateKey := tx.Date.Format("2006-01-02")
		daily, exists := dailyMap[dateKey]

		if !exists {
			// Инициализируем новую запись для дня
			isHoliday := s.isHoliday(tx.Date, holidays)
			daily = &entities.DailyTransactionData{
				Date:         tx.Date,
				Sales:        0,
				TotalPrice:   0,
				AvgDiscount:  0,
				DiscountedTx: 0,
				TotalTx:      0,
				IsHoliday:    isHoliday,
			}
			dailyMap[dateKey] = daily
		}

		// Находим интересующий нас товар в транзакции
		for _, item := range tx.Items {
			if item.ID == productID {
				daily.Sales += float64(item.Quantity)
				daily.TotalPrice += item.Price * float64(item.Quantity)

				if item.DiscountPct > 0 {
					daily.AvgDiscount += item.DiscountPct * float64(item.Quantity)
					daily.DiscountedTx++
				}

				daily.TotalTx++
			}
		}
	}

	// Вычисляем средние значения для каждого дня
	result := make([]*entities.DailyTransactionData, 0, len(dailyMap))
	for _, daily := range dailyMap {
		if daily.Sales > 0 && daily.TotalTx > 0 {
			// Вычисляем среднюю скидку
			if daily.DiscountedTx > 0 {
				daily.AvgDiscount /= daily.Sales
			}

			result = append(result, daily)
		}
	}

	return result
}

// aggregateCategoryTransactionsByDay агрегирует транзакции по категории и дням
func (s *regressionServiceImpl) aggregateCategoryTransactionsByDay(transactions []*entities.Transaction, category string) []*entities.DailyTransactionData {
	// Реализация аналогична aggregateTransactionsByDay, но для категории
	// ...
	// Упрощенная реализация
	dailyMap := make(map[string]*entities.DailyTransactionData)

	// Определяем праздничные дни
	holidays := s.getHolidayDates()

	// Агрегируем данные по дням
	for _, tx := range transactions {
		dateKey := tx.Date.Format("2006-01-02")
		daily, exists := dailyMap[dateKey]

		if !exists {
			isHoliday := s.isHoliday(tx.Date, holidays)
			daily = &entities.DailyTransactionData{
				Date:         tx.Date,
				Sales:        0,
				TotalPrice:   0,
				AvgDiscount:  0,
				DiscountedTx: 0,
				TotalTx:      0,
				IsHoliday:    isHoliday,
				ProductCount: 0,
				ProductIDs:   make(map[string]struct{}),
			}
			dailyMap[dateKey] = daily
		}

		// Находим товары из нужной категории
		for _, item := range tx.Items {
			if item.Category == category {
				daily.Sales += float64(item.Quantity)
				daily.TotalPrice += item.Price * float64(item.Quantity)

				if item.DiscountPct > 0 {
					daily.AvgDiscount += item.DiscountPct * float64(item.Quantity)
					daily.DiscountedTx++
				}

				daily.TotalTx++
				daily.ProductIDs[item.ID] = struct{}{}
			}
		}
	}

	// Вычисляем средние значения для каждого дня
	result := make([]*entities.DailyTransactionData, 0, len(dailyMap))
	for _, daily := range dailyMap {
		if daily.Sales > 0 && daily.TotalTx > 0 {
			// Вычисляем среднюю скидку
			if daily.DiscountedTx > 0 {
				daily.AvgDiscount /= daily.Sales
			}

			// Количество уникальных товаров
			daily.ProductCount = len(daily.ProductIDs)

			result = append(result, daily)
		}
	}

	return result
}

// findOptimalDiscount находит оптимальный уровень скидки для максимизации дохода
func (s *regressionServiceImpl) findOptimalDiscount(r *regression.Regression, productID string) float64 {
	// Находим базовую цену товара
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return 0.15 // Если не удалось получить товар, возвращаем стандартное значение
	}

	// Получаем коэффициент при скидке
	discountCoeff := r.Coeff(0)

	// Простая модель: Sales = Intercept + DiscountCoeff * Discount
	// Revenue = Price * (1 - Discount) * Sales
	// Оптимальный уровень скидки - производная Revenue по Discount = 0

	// Для линейной модели Sales = a + b*Discount
	// Revenue = Price * (1 - Discount) * (a + b*Discount)
	// dRevenue/dDiscount = Price * (-a - b*Discount + b - b*Discount) = Price * (b - a - 2*b*Discount)
	// Приравниваем к нулю: b - a - 2*b*Discount = 0
	// Discount = (b - a) / (2*b)

	intercept := r.Coeff(4) // Значение свободного члена

	// Если коэффициент при скидке положительный (скидка увеличивает продажи)
	if discountCoeff > 0 {
		optimalDiscount := (discountCoeff - intercept) / (2 * discountCoeff)

		// Ограничиваем скидку разумными пределами
		if optimalDiscount < 0.05 {
			return 0.05 // Минимальная скидка 5%
		} else if optimalDiscount > 0.5 {
			return 0.5 // Максимальная скидка 50%
		}
		return optimalDiscount
	}

	// Если скидка имеет отрицательное влияние на продажи
	return 0.0 // Рекомендуем не делать скидку
}

// findOptimalCategoryDiscount находит оптимальный уровень скидки для категории
func (s *regressionServiceImpl) findOptimalCategoryDiscount(r *regression.Regression, category string) float64 {
	// Аналогично findOptimalDiscount, но для категории
	// ...
	// Упрощенная реализация

	// Получаем коэффициент при скидке и свободный член
	discountCoeff := r.Coeff(0)
	intercept := r.Coeff(5)

	if discountCoeff > 0 {
		optimalDiscount := (discountCoeff - intercept) / (2 * discountCoeff)

		// Ограничиваем скидку разумными пределами
		if optimalDiscount < 0.05 {
			return 0.05
		} else if optimalDiscount > 0.5 {
			return 0.5
		}
		return optimalDiscount
	}

	return 0.0
}

// adjustRecommendationByABCCategory корректирует рекомендацию на основе ABC категории
func (s *regressionServiceImpl) adjustRecommendationByABCCategory(rec *entities.DiscountRecommendation) {
	// Корректируем скидку на основе ABC-категории
	switch rec.ABCCategory {
	case "A":
		// Категория A - высокодоходные товары, минимальные скидки
		if rec.OptimalDiscount > 0.2 {
			rec.OptimalDiscount = 0.2
			rec.AdjustmentReason = "Скидка ограничена для высокодоходной категории A"
		}
	case "B":
		// Категория B - средние скидки
		if rec.OptimalDiscount > 0.3 {
			rec.OptimalDiscount = 0.3
			rec.AdjustmentReason = "Скидка скорректирована для категории B"
		}
	case "C":
		// Категория C - более высокие скидки для стимулирования продаж
		if rec.OptimalDiscount < 0.1 && rec.LiftFactor > 0 {
			rec.OptimalDiscount = 0.1
			rec.AdjustmentReason = "Скидка увеличена для стимулирования низкодоходной категории C"
		}
	}
}

// findOptimalDiscountFromABTests находит оптимальный уровень скидки на основе A/B тестов
func (s *regressionServiceImpl) findOptimalDiscountFromABTests(r *regression.Regression, basePrice, duration float64) float64 {
	// Функция для расчета прогнозируемого Lift при заданной скидке
	predictLift := func(discount float64) float64 {
		// Прогнозируем Lift на основе регрессионной модели
		return r.Coeff(3) + r.Coeff(0)*discount + r.Coeff(1)*basePrice + r.Coeff(2)*duration
	}

	// Функция для расчета изменения дохода
	calcRevenueChange := func(discount float64) float64 {
		lift := predictLift(discount)
		// Изменение дохода = Lift * (1 - скидка)
		return lift * (1 - discount)
	}

	// Находим оптимальную скидку методом перебора с шагом 0.01
	maxRevenue := 0.0
	optimalDiscount := 0.0

	for discount := 0.01; discount <= 0.5; discount += 0.01 {
		revenue := calcRevenueChange(discount)
		if revenue > maxRevenue {
			maxRevenue = revenue
			optimalDiscount = discount
		}
	}

	return optimalDiscount
}

// isHoliday проверяет, является ли дата праздничным днем
func (s *regressionServiceImpl) isHoliday(date time.Time, holidays map[string]struct{}) bool {
	dateKey := date.Format("01-02") // MM-DD формат
	_, isHoliday := holidays[dateKey]

	// Также считаем выходными субботу и воскресенье
	isWeekend := date.Weekday() == time.Saturday || date.Weekday() == time.Sunday

	return isHoliday || isWeekend
}

// getHolidayDates возвращает список праздничных дней
func (s *regressionServiceImpl) getHolidayDates() map[string]struct{} {
	// Упрощенный список основных праздников (MM-DD формат)
	holidays := map[string]struct{}{
		"01-01": {}, // Новый год
		"02-14": {}, // День святого Валентина
		"03-08": {}, // Международный женский день
		"05-01": {}, // Праздник весны и труда
		"05-09": {}, // День Победы
		"06-12": {}, // День России
		"11-04": {}, // День народного единства
		"12-25": {}, // Рождество
		"12-31": {}, // Канун Нового года
	}
	return holidays
}

// getAverageBasePriceForTest возвращает среднюю базовую цену товаров в тесте
func (s *regressionServiceImpl) getAverageBasePriceForTest(test *entities.ABTestResult) (float64, error) {
	if test == nil {
		return 0, ErrInvalidParameter
	}

	var totalPrice float64
	var totalProducts int

	// Получаем все товары, участвовавшие в тесте
	products, err := s.productRepo.GetProductsByTestID(test.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to get test products: %w", err)
	}

	if len(products) == 0 {
		return 0, ErrInsufficientData
	}

	// Суммируем базовые цены всех товаров
	for _, product := range products {
		totalPrice += product.BasePrice
		totalProducts++
	}

	// Вычисляем среднюю цену
	avgPrice := totalPrice / float64(totalProducts)

	return avgPrice, nil
}

// getAverageBasePrice возвращает среднюю базовую цену всех товаров
func (s *regressionServiceImpl) getAverageBasePrice() (float64, error) {
	// Получаем все активные товары
	products, err := s.productRepo.GetAllActive()
	if err != nil {
		return 0, fmt.Errorf("failed to get products: %w", err)
	}

	if len(products) == 0 {
		return 0, ErrInsufficientData
	}

	var totalPrice float64
	for _, product := range products {
		totalPrice += product.BasePrice
	}

	return totalPrice / float64(len(products)), nil
}
