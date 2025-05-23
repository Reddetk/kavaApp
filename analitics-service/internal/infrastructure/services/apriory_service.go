package services

import (
	"context"
	"fmt"
	"sort"

	"analitics-service/internal/domain/entities"
	"analitics-service/pkg/logger"

	apriori "github.com/eMAGTechLabs/go-apriori"
)

// AprioriService определяет интерфейс для анализа ассоциативных правил с использованием алгоритма Apriori
type AprioriService interface {
	// GenerateFrequentItemsets генерирует частые наборы товаров из транзакций
	// minSupport - минимальная поддержка (от 0 до 1)
	GenerateFrequentItemsets(ctx context.Context, transactions []entities.Transaction, minSupport float64) ([]entities.FrequentItemset, error)

	// GenerateAssociationRules генерирует ассоциативные правила из частых наборов
	// minConfidence - минимальная достоверность (от 0 до 1)
	GenerateAssociationRules(ctx context.Context, frequentItemsets []entities.FrequentItemset, minConfidence float64) ([]entities.AssociationRule, error)

	// AnalyzeTransactions выполняет полный анализ транзакций, возвращая ассоциативные правила
	// объединяет два предыдущих метода в одну операцию
	AnalyzeTransactions(ctx context.Context, transactions []entities.Transaction, minSupport, minConfidence float64) ([]entities.AssociationRule, error)

	// GetProductRecommendations возвращает рекомендации товаров на основе корзины товаров пользователя
	GetProductRecommendations(ctx context.Context, currentBasket []entities.Product, rules []entities.AssociationRule, limit int) ([]entities.ProductRecommendation, error)
}

// aprioriService реализует интерфейс AprioriService
type aprioriService struct {
	logger logger.Logger
}

// NewAprioriService создает новый экземпляр сервиса Apriori
func NewAprioriService(logger logger.Logger) *aprioriService {
	return &aprioriService{
		logger: logger,
	}
}

// GenerateFrequentItemsets генерирует частые наборы товаров из транзакций
func (s *aprioriService) GenerateFrequentItemsets(ctx context.Context, transactions []entities.Transaction, minSupport float64) ([]entities.FrequentItemset, error) {
	s.logger.Info(ctx, "Генерация частых наборов товаров", "транзакций", len(transactions), "minSupport", minSupport)

	// Преобразуем транзакции в формат, требуемый библиотекой go-apriori
	itemMatrix := prepareTransactionsData(transactions)

	// Создаем новый экземпляр обработчика Apriori
	ap := apriori.NewApriori(itemMatrix)

	// Генерируем частые наборы с указанной минимальной поддержкой
	aprioriResults := ap.Calculate(minSupport)

	// Преобразуем результаты библиотеки в наши доменные сущности
	result := make([]entities.FrequentItemset, 0, len(aprioriResults))
	for _, apResult := range aprioriResults {
		itemset := entities.FrequentItemset{
			Items:   convertAprioriItems(apResult.Items),
			Support: apResult.Support,
			Count:   int(apResult.Support * float64(len(transactions))),
		}
		result = append(result, itemset)
	}

	// Сортируем результаты по убыванию поддержки
	sort.Slice(result, func(i, j int) bool {
		return result[i].Support > result[j].Support
	})

	s.logger.Info(ctx, "Сгенерированы частые наборы товаров", "количество", len(result))
	return result, nil
}

// GenerateAssociationRules генерирует ассоциативные правила из частых наборов
func (s *aprioriService) GenerateAssociationRules(ctx context.Context, frequentItemsets []entities.FrequentItemset, minConfidence float64) ([]entities.AssociationRule, error) {
	s.logger.Info(ctx, "Генерация ассоциативных правил", "наборов", len(frequentItemsets), "minConfidence", minConfidence)

	// Преобразуем наши доменные сущности в формат, требуемый для библиотеки
	aprioriItemsets := make([]apriori.Itemset, 0, len(frequentItemsets))
	for _, itemset := range frequentItemsets {
		aprioriItemset := apriori.Itemset{
			Items:   convertToAprioriItems(itemset.Items),
			Support: itemset.Support,
		}
		aprioriItemsets = append(aprioriItemsets, aprioriItemset)
	}

	// Генерируем правила ассоциаций
	aprioriRules := apriori.GenerateAssociationRules(aprioriItemsets, minConfidence)

	// Преобразуем результаты в наши доменные сущности
	rules := make([]entities.AssociationRule, 0, len(aprioriRules))
	for _, rule := range aprioriRules {
		domainRule := entities.AssociationRule{
			Antecedent: convertAprioriItems(rule.From),
			Consequent: convertAprioriItems(rule.To),
			Support:    rule.Support,
			Confidence: rule.Confidence,
			Lift:       calculateLift(rule.From, rule.To, frequentItemsets),
		}
		rules = append(rules, domainRule)
	}

	// Сортируем правила по убыванию уверенности
	sort.Slice(rules, func(i, j int) bool {
		if rules[i].Confidence == rules[j].Confidence {
			return rules[i].Support > rules[j].Support
		}
		return rules[i].Confidence > rules[j].Confidence
	})

	s.logger.Info(ctx, "Сгенерированы ассоциативные правила", "количество", len(rules))
	return rules, nil
}

// AnalyzeTransactions выполняет полный анализ транзакций, возвращая ассоциативные правила
func (s *aprioriService) AnalyzeTransactions(ctx context.Context, transactions []entities.Transaction, minSupport, minConfidence float64) ([]entities.AssociationRule, error) {
	// Генерируем частые наборы
	itemsets, err := s.GenerateFrequentItemsets(ctx, transactions, minSupport)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации частых наборов: %w", err)
	}

	// Если нет частых наборов, возвращаем пустой результат
	if len(itemsets) == 0 {
		s.logger.Warn(ctx, "Не найдено частых наборов товаров с указанной поддержкой", "minSupport", minSupport)
		return []entities.AssociationRule{}, nil
	}

	// Генерируем ассоциативные правила
	rules, err := s.GenerateAssociationRules(ctx, itemsets, minConfidence)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации ассоциативных правил: %w", err)
	}

	return rules, nil
}

// GetProductRecommendations возвращает рекомендации товаров на основе корзины товаров пользователя
func (s *aprioriService) GetProductRecommendations(ctx context.Context, currentBasket []entities.Product, rules []entities.AssociationRule, limit int) ([]entities.ProductRecommendation, error) {
	s.logger.Info(ctx, "Получение рекомендаций товаров", "корзина", len(currentBasket), "правила", len(rules))

	if len(currentBasket) == 0 {
		return []entities.ProductRecommendation{}, nil
	}

	// Преобразуем корзину в множество ID товаров для быстрого поиска
	basketIDs := make(map[string]bool)
	for _, product := range currentBasket {
		basketIDs[product.ID] = true
	}

	// Создаем карту для объединения рекомендаций для одинаковых товаров
	recommendationsMap := make(map[string]entities.ProductRecommendation)

	// Оцениваем каждое ассоциативное правило
	for _, rule := range rules {
		// Проверяем, совпадает ли antecedent с товарами в корзине
		// Все элементы из antecedent должны быть в корзине
		matchedAll := true
		for _, antItem := range rule.Antecedent {
			if !basketIDs[antItem.ID] {
				matchedAll = false
				break
			}
		}

		if !matchedAll {
			continue
		}

		// Для каждого товара из consequent создаем или обновляем рекомендацию
		for _, conseqItem := range rule.Consequent {
			// Не рекомендуем товары, которые уже есть в корзине
			if basketIDs[conseqItem.ID] {
				continue
			}

			// Если товар уже есть в рекомендациях, обновляем его score
			if rec, exists := recommendationsMap[conseqItem.ID]; exists {
				// Используем максимальное значение confidence как основной показатель
				if rule.Confidence > rec.Score {
					rec.Score = rule.Confidence
					rec.Lift = rule.Lift
					rec.Support = rule.Support
					recommendationsMap[conseqItem.ID] = rec
				}
			} else {
				// Создаем новую рекомендацию
				product, err := s.getProductDetails(ctx, conseqItem.ID)
				if err != nil {
					s.logger.Error(ctx, "Ошибка получения деталей товара", "error", err, "productID", conseqItem.ID)
					continue
				}

				recommendationsMap[conseqItem.ID] = entities.ProductRecommendation{
					Product: product,
					Score:   rule.Confidence,
					Lift:    rule.Lift,
					Support: rule.Support,
				}
			}
		}
	}

	// Преобразуем карту в слайс и сортируем по score
	recommendations := make([]entities.ProductRecommendation, 0, len(recommendationsMap))
	for _, rec := range recommendationsMap {
		recommendations = append(recommendations, rec)
	}

	// Сортируем по убыванию score
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Ограничиваем количество рекомендаций
	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	s.logger.Info(ctx, "Получены рекомендации товаров", "количество", len(recommendations))
	return recommendations, nil
}

// Вспомогательные функции

// prepareTransactionsData преобразует транзакции в формат для библиотеки go-apriori
func prepareTransactionsData(transactions []entities.Transaction) [][]apriori.Item {
	itemMatrix := make([][]apriori.Item, 0, len(transactions))

	for _, transaction := range transactions {
		items := make([]apriori.Item, 0, len(transaction.Items))
		for _, item := range transaction.Items {
			items = append(items, apriori.Item(item.ID))
		}
		itemMatrix = append(itemMatrix, items)
	}

	return itemMatrix
}

// convertAprioriItems преобразует items из формата библиотеки в наш формат
func convertAprioriItems(apItems []apriori.Item) []entities.Item {
	items := make([]entities.Item, 0, len(apItems))
	for _, item := range apItems {
		items = append(items, entities.Item{
			ID: string(item),
		})
	}
	return items
}

// convertToAprioriItems преобразует наши items в формат библиотеки
func convertToAprioriItems(items []entities.Item) []apriori.Item {
	apItems := make([]apriori.Item, 0, len(items))
	for _, item := range items {
		apItems = append(apItems, apriori.Item(item.ID))
	}
	return apItems
}

// calculateLift вычисляет показатель Lift для правила
func calculateLift(antecedent []apriori.Item, consequent []apriori.Item, itemsets []entities.FrequentItemset) float64 {
	// Находим поддержку antecedent
	antSupport := findItemsetSupport(convertAprioriItems(antecedent), itemsets)

	// Находим поддержку consequent
	consSupport := findItemsetSupport(convertAprioriItems(consequent), itemsets)

	// Находим общую поддержку
	combined := append([]apriori.Item{}, antecedent...)
	combined = append(combined, consequent...)
	combSupport := findItemsetSupport(convertAprioriItems(combined), itemsets)

	// Вычисляем lift
	// Lift = P(A∪B) / (P(A) * P(B))
	if antSupport > 0 && consSupport > 0 {
		return combSupport / (antSupport * consSupport)
	}

	return 0
}

// findItemsetSupport находит поддержку itemset в списке частых наборов
func findItemsetSupport(items []entities.Item, itemsets []entities.FrequentItemset) float64 {
	// Создаем карту ID товаров для быстрого поиска
	itemIDs := make(map[string]bool)
	for _, item := range items {
		itemIDs[item.ID] = true
	}

	// Ищем точное совпадение
	for _, itemset := range itemsets {
		if len(itemset.Items) != len(items) {
			continue
		}

		matchAll := true
		for _, item := range itemset.Items {
			if !itemIDs[item.ID] {
				matchAll = false
				break
			}
		}

		if matchAll {
			return itemset.Support
		}
	}

	return 0
}

// getProductDetails получает детали товара по ID
// В реальной имплементации этот метод должен обращаться к репозиторию
func (s *aprioriService) getProductDetails(ctx context.Context, productID string) (entities.Product, error) {
	// Здесь должен быть вызов репозитория
	// return s.productRepository.GetByID(ctx, productID)

	// Временный заглушка
	return entities.Product{
		ID:    productID,
		Name:  "Product " + productID,
		Price: 0,
	}, nil
}
