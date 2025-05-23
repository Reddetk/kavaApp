package services

import (
	"context"
	"sort"

	"analitics-service/internal/domain/entities"
	"analitics-service/internal/domain/repositories"
)

// ABCAnalysisService предоставляет функционал для ABC-анализа товаров
type ABCAnalysisService interface {
	// PerformABCAnalysis выполняет ABC-анализ товаров на основе переданных критериев
	PerformABCAnalysis(ctx context.Context, criteria entities.ABCAnalysisCriteria) (*entities.ABCAnalysisResult, error)

	// GetProductSegmentation возвращает сегментацию продуктов по категориям A, B, C
	GetProductSegmentation(ctx context.Context, productID string) (*entities.ProductSegmentation, error)

	// GetSegmentSummary возвращает сводную информацию по сегментам
	GetSegmentSummary(ctx context.Context) (*entities.ABCSegmentSummary, error)
}

// ABCAnalysisServiceImpl реализация сервиса ABC-анализа
type ABCAnalysisServiceImpl struct {
	productRepo      repositories.ProductRepository
	salesRepo        repositories.SalesRepository
	abcSegmentRepo   repositories.ABCSegmentRepository
	profitMarginRepo repositories.ProfitMarginRepository
}

// NewABCAnalysisService создает новый экземпляр сервиса ABC-анализа
func NewABCAnalysisService(
	productRepo repositories.ProductRepository,
	salesRepo repositories.SalesRepository,
	abcSegmentRepo repositories.ABCSegmentRepository,
	profitMarginRepo repositories.ProfitMarginRepository,
) ABCAnalysisService {
	return &ABCAnalysisServiceImpl{
		productRepo:      productRepo,
		salesRepo:        salesRepo,
		abcSegmentRepo:   abcSegmentRepo,
		profitMarginRepo: profitMarginRepo,
	}
}

// PerformABCAnalysis выполняет многокритериальный ABC-анализ товаров
func (s *ABCAnalysisServiceImpl) PerformABCAnalysis(ctx context.Context, criteria entities.ABCAnalysisCriteria) (*entities.ABCAnalysisResult, error) {
	// Получаем все продукты
	products, err := s.productRepo.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	// Получаем данные о продажах за указанный период
	sales, err := s.salesRepo.GetSalesByPeriod(ctx, criteria.StartDate, criteria.EndDate)
	if err != nil {
		return nil, err
	}

	// Получаем данные о прибыльности продуктов
	profitMargins, err := s.profitMarginRepo.GetProfitMargins(ctx)
	if err != nil {
		return nil, err
	}

	// Подготавливаем данные для анализа
	productsData := prepareProductsData(products, sales, profitMargins)

	// Выполняем анализ по каждому критерию
	revenueSegmentation := s.analyzeByRevenue(productsData, criteria.ThresholdsRevenue)
	quantitySegmentation := s.analyzeByQuantity(productsData, criteria.ThresholdsQuantity)
	profitSegmentation := s.analyzeByProfit(productsData, criteria.ThresholdsProfit)

	// Объединяем результаты анализа по разным критериям
	finalSegmentation := s.combineSegmentations(
		revenueSegmentation,
		quantitySegmentation,
		profitSegmentation,
		criteria.Weights,
	)

	// Сохраняем результаты в репозиторий
	err = s.abcSegmentRepo.SaveSegmentation(ctx, finalSegmentation)
	if err != nil {
		return nil, err
	}

	// Формируем и возвращаем результат
	return &ABCAnalysisResult{
		AnalysisDate:         criteria.EndDate,
		ProductsSegmentation: finalSegmentation,
		Summary:              calculateSummary(finalSegmentation),
	}, nil
}

// GetProductSegmentation возвращает информацию о сегментации конкретного продукта
func (s *ABCAnalysisServiceImpl) GetProductSegmentation(ctx context.Context, productID string) (*entities.ProductSegmentation, error) {
	// Получаем информацию о сегментации продукта из репозитория
	segmentation, err := s.abcSegmentRepo.GetProductSegmentation(ctx, productID)
	if err != nil {
		return nil, err
	}

	return segmentation, nil
}

// GetSegmentSummary возвращает сводную информацию по сегментам A, B, C
func (s *ABCAnalysisServiceImpl) GetSegmentSummary(ctx context.Context) (*entities.ABCSegmentSummary, error) {
	// Получаем полную информацию о сегментации
	segmentation, err := s.abcSegmentRepo.GetFullSegmentation(ctx)
	if err != nil {
		return nil, err
	}

	// Рассчитываем сводную информацию
	summary := calculateSummary(segmentation)

	return summary, nil
}

// prepareProductsData подготавливает данные о продуктах для анализа
func prepareProductsData(products []entities.Product, sales []entities.Sale, profitMargins map[string]float64) []ProductAnalysisData {
	// Агрегируем данные по каждому продукту
	productData := make(map[string]ProductAnalysisData)

	// Инициализируем данные о продуктах
	for _, product := range products {
		productData[product.ID] = ProductAnalysisData{
			Product:      product,
			Revenue:      0,
			Quantity:     0,
			Profit:       0,
			ProfitMargin: profitMargins[product.ID],
		}
	}

	// Агрегируем данные о продажах
	for _, sale := range sales {
		if data, exists := productData[sale.ProductID]; exists {
			data.Revenue += sale.Price * float64(sale.Quantity)
			data.Quantity += sale.Quantity
			data.Profit += (sale.Price * float64(sale.Quantity)) * (data.ProfitMargin / 100.0)
			productData[sale.ProductID] = data
		}
	}

	// Конвертируем map в массив для дальнейшей обработки
	result := make([]ProductAnalysisData, 0, len(productData))
	for _, data := range productData {
		result = append(result, data)
	}

	return result
}

// Define the missing types that were previously used without proper imports
type ProductAnalysisData struct {
	Product      entities.Product
	Revenue      float64
	Quantity     int
	Profit       float64
	ProfitMargin float64
}

// analyzeByRevenue выполняет ABC-анализ по выручке
func (s *ABCAnalysisServiceImpl) analyzeByRevenue(productsData []ProductAnalysisData, thresholds entities.Thresholds) map[string]string {
	// Сортируем продукты по выручке в порядке убывания
	sort.Slice(productsData, func(i, j int) bool {
		return productsData[i].Revenue > productsData[j].Revenue
	})

	// Рассчитываем общую выручку
	totalRevenue := 0.0
	for _, data := range productsData {
		totalRevenue += data.Revenue
	}

	// Определяем сегменты
	return determineSegments(productsData, totalRevenue, func(data ProductAnalysisData) float64 {
		return data.Revenue
	}, thresholds)
}

// analyzeByQuantity выполняет ABC-анализ по количеству продаж
func (s *ABCAnalysisServiceImpl) analyzeByQuantity(productsData []ProductAnalysisData, thresholds entities.Thresholds) map[string]string {
	// Сортируем продукты по количеству продаж в порядке убывания
	sort.Slice(productsData, func(i, j int) bool {
		return productsData[i].Quantity > productsData[j].Quantity
	})

	// Рассчитываем общее количество
	totalQuantity := 0
	for _, data := range productsData {
		totalQuantity += data.Quantity
	}

	// Определяем сегменты
	return determineSegments(productsData, float64(totalQuantity), func(data ProductAnalysisData) float64 {
		return float64(data.Quantity)
	}, thresholds)
}

// analyzeByProfit выполняет ABC-анализ по прибыли
func (s *ABCAnalysisServiceImpl) analyzeByProfit(productsData []ProductAnalysisData, thresholds entities.Thresholds) map[string]string {
	// Сортируем продукты по прибыли в порядке убывания
	sort.Slice(productsData, func(i, j int) bool {
		return productsData[i].Profit > productsData[j].Profit
	})

	// Рассчитываем общую прибыль
	totalProfit := 0.0
	for _, data := range productsData {
		totalProfit += data.Profit
	}

	// Определяем сегменты
	return determineSegments(productsData, totalProfit, func(data ProductAnalysisData) float64 {
		return data.Profit
	}, thresholds)
}

// determineSegments определяет сегменты A, B, C на основе кумулятивного процента
func determineSegments(productsData []ProductAnalysisData, total float64, valueFunc func(ProductAnalysisData) float64, thresholds entities.Thresholds) map[string]string {
	segments := make(map[string]string)
	cumulativePercent := 0.0

	for _, data := range productsData {
		value := valueSelector(data)
		percent := (value / total) * 100
		cumulativePercent += percent

		// Определяем сегмент на основе кумулятивного процента
		switch {
		case cumulativePercent <= thresholds.AThreshold:
			segments[data.Product.ID] = "A"
		case cumulativePercent <= thresholds.BThreshold:
			segments[data.Product.ID] = "B"
		default:
			segments[data.Product.ID] = "C"
		}
	}

	return segments
}

// combineSegmentations объединяет результаты сегментаций по разным критериям
func (s *ABCAnalysisServiceImpl) combineSegmentations(
	revenueSegmentation map[string]string,
	quantitySegmentation map[string]string,
	profitSegmentation map[string]string,
	weights entities.CriteriaWeights,
) map[string]entities.ProductFullSegmentation {

	combinedSegmentation := make(map[string]ProductFullSegmentation)

	// Присваиваем числовые значения сегментам (A=3, B=2, C=1)
	segmentValues := map[string]int{
		"A": 3,
		"B": 2,
		"C": 1,
	}

	// Объединяем сегментации для каждого продукта
	for productID, revenueSegment := range revenueSegmentation {
		quantitySegment := quantitySegmentation[productID]
		profitSegment := profitSegmentation[productID]

		// Рассчитываем взвешенную оценку
		weightedScore := float64(segmentValues[revenueSegment])*weights.RevenueWeight +
			float64(segmentValues[quantitySegment])*weights.QuantityWeight +
			float64(segmentValues[profitSegment])*weights.ProfitWeight

		// Определяем финальный сегмент на основе взвешенной оценки
		var finalSegment string
		switch {
		case weightedScore >= 2.5:
			finalSegment = "A"
		case weightedScore >= 1.5:
			finalSegment = "B"
		default:
			finalSegment = "C"
		}

		// Сохраняем результат
		combinedSegmentation[productID] = ProductFullSegmentation{
			ProductID:       productID,
			RevenueSegment:  revenueSegment,
			QuantitySegment: quantitySegment,
			ProfitSegment:   profitSegment,
			FinalSegment:    finalSegment,
			Score:           weightedScore,
		}
	}

	return combinedSegmentation
}

// calculateSummary рассчитывает сводную информацию по сегментам
func calculateSummary(segmentation map[string]entities.ProductFullSegmentation) *entities.ABCSegmentSummary {
	summary := &ABCSegmentSummary{
		SegmentCounts: map[string]int{
			"A": 0,
			"B": 0,
			"C": 0,
		},
		SegmentPercentages: map[string]float64{
			"A": 0,
			"B": 0,
			"C": 0,
		},
	}

	// Подсчитываем количество продуктов в каждом сегменте
	totalProducts := 0
	for _, seg := range segmentation {
		summary.SegmentCounts[seg.FinalSegment]++
		totalProducts++
	}

	// Рассчитываем процентное соотношение
	if totalProducts > 0 {
		for segment, count := range summary.SegmentCounts {
			summary.SegmentPercentages[segment] = float64(count) / float64(totalProducts) * 100
		}
	}

	return summary
}
