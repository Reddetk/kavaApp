package services

import (
	"errors"
	"math"
	"math/rand"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

// KMeansSegmentation представляет сервис для сегментации пользователей на основе алгоритма K-means
type KMeansSegmentation struct {
	numClusters int
	maxIter     int
	epsilon     float64
}

// NewKMeansSegmentation создает новый экземпляр сервиса сегментации K-means
func NewKMeansSegmentation(numClusters int) services.SegmentationService {
	return &KMeansSegmentation{
		numClusters: numClusters,
		maxIter:     100,  // Максимальное количество итераций
		epsilon:     1e-6, // Порог сходимости
	}
}

// PerformRFMClustering выполняет кластеризацию пользователей на основе RFM-метрик
func (s *KMeansSegmentation) PerformRFMClustering(users []entities.UserMetrics) ([]entities.Segment, error) {
	if len(users) == 0 {
		return nil, errors.New("no user metrics provided")
	}

	// Подготовка данных и нормализация
	points, ranges := s.prepareAndNormalizeRFMData(users)

	// Выполнение кластеризации методом K-means++
	centroids, assignments := s.kmeansClusteringPlusPlus(points)

	// Денормализация центроидов
	denormalizedCentroids := s.denormalizeCentroids(centroids, ranges)

	// Анализ кластеров для определения их характеристик
	clusterCharacteristics := s.analyzeRFMClusters(denormalizedCentroids)

	// Создание объектов сегментов на основе результатов кластеризации
	segments := make([]entities.Segment, s.numClusters)
	for i := 0; i < s.numClusters; i++ {
		// Подсчитаем количество пользователей в сегменте
		userCount := 0
		for _, a := range assignments {
			if a == i {
				userCount++
			}
		}

		segments[i] = entities.Segment{
			ID:        uuid.New(),
			Name:      s.generateRFMSegmentName(clusterCharacteristics[i]),
			Type:      "RFM",
			Algorithm: "KMeans",
			CentroidData: map[string]interface{}{
				"values":          denormalizedCentroids[i],
				"metrics":         []string{"recency", "frequency", "monetary"},
				"characteristics": clusterCharacteristics[i],
				"user_count":      userCount,
				"percentage":      float64(userCount) / float64(len(users)) * 100,
			},
			CreatedAt: time.Now(),
		}
	}

	return segments, nil
}

// PerformBehaviorClustering выполняет кластеризацию пользователей на основе поведенческих метрик
func (s *KMeansSegmentation) PerformBehaviorClustering(transactions []entities.Transaction) ([]entities.Segment, error) {
	if len(transactions) == 0 {
		return nil, errors.New("no transactions provided")
	}

	// Агрегируем транзакции по пользователям и рассчитываем поведенческие метрики
	userBehaviors, userIDs := s.aggregateTransactionBehaviors(transactions)

	// Нормализация данных
	points, ranges := s.normalizeBehaviorData(userBehaviors)

	// Выполнение кластеризации методом K-means++
	centroids, assignments := s.kmeansClusteringPlusPlus(points)

	// Денормализация центроидов
	denormalizedCentroids := s.denormalizeCentroids(centroids, ranges)

	// Анализ кластеров для определения их характеристик
	clusterCharacteristics := s.analyzeBehaviorClusters(denormalizedCentroids)

	// Создание объектов сегментов на основе результатов кластеризации
	segments := make([]entities.Segment, s.numClusters)
	for i := 0; i < s.numClusters; i++ {
		// Подсчитаем количество пользователей в сегменте
		userCount := 0
		for _, a := range assignments {
			if a == i {
				userCount++
			}
		}

		segments[i] = entities.Segment{
			ID:        uuid.New(),
			Name:      s.generateBehaviorSegmentName(clusterCharacteristics[i]),
			Type:      "Behavioral",
			Algorithm: "KMeans",
			CentroidData: map[string]interface{}{
				"values": denormalizedCentroids[i],
				"metrics": []string{
					"avg_transaction_amount",
					"transaction_frequency",
					"preferred_time",
				},
				"characteristics": clusterCharacteristics[i],
				"user_count":      userCount,
				"percentage":      float64(userCount) / float64(len(userIDs)) * 100,
			},
			CreatedAt: time.Now(),
		}
	}

	return segments, nil
}

// AssignUserToSegment определяет сегмент пользователя на основе его RFM-метрик
func (s *KMeansSegmentation) AssignUserToSegment(userID uuid.UUID, metrics entities.UserMetrics, segments []entities.Segment) (entities.Segment, error) {
	if len(segments) == 0 {
		return entities.Segment{}, errors.New("no segments provided")
	}

	// Преобразование метрик пользователя в точку
	point := []float64{
		float64(metrics.Recency),
		float64(metrics.Frequency),
		metrics.Monetary,
	}

	// Извлечение центроидов из сегментов
	centroids := make([][]float64, len(segments))
	for i, segment := range segments {
		if values, ok := segment.CentroidData["values"].([]float64); ok && len(values) >= 3 {
			centroids[i] = values
		} else {
			// Если центроид не найден, используем нулевые значения
			centroids[i] = make([]float64, 3)
		}
	}

	// Найти ближайший центроид
	closestIdx := s.closestCentroid(point, centroids)

	return segments[closestIdx], nil
}

// prepareAndNormalizeRFMData подготавливает и нормализует RFM данные для кластеризации
func (s *KMeansSegmentation) prepareAndNormalizeRFMData(users []entities.UserMetrics) ([][]float64, [][]float64) {
	// Получаем исходные ненормализованные данные
	points := make([][]float64, len(users))
	for i, u := range users {
		points[i] = []float64{
			float64(u.Recency),
			float64(u.Frequency),
			u.Monetary,
		}
	}

	// Находим минимальные и максимальные значения для каждой метрики
	minValues := make([]float64, 3)
	maxValues := make([]float64, 3)

	for j := 0; j < 3; j++ {
		minValues[j] = math.MaxFloat64
		maxValues[j] = -math.MaxFloat64

		for i := 0; i < len(points); i++ {
			if points[i][j] < minValues[j] {
				minValues[j] = points[i][j]
			}
			if points[i][j] > maxValues[j] {
				maxValues[j] = points[i][j]
			}
		}
	}

	// Особая обработка для recency: инвертируем значения, чтобы меньшие значения (недавние покупки) были более значимыми
	for i := 0; i < len(points); i++ {
		// Инвертируем recency: чем меньше дней прошло, тем выше значение
		if maxValues[0] > minValues[0] {
			points[i][0] = 1.0 - (points[i][0]-minValues[0])/(maxValues[0]-minValues[0])
		} else {
			points[i][0] = 0.5 // Если все значения одинаковые
		}

		// Нормализуем frequency и monetary
		for j := 1; j < 3; j++ {
			if maxValues[j] > minValues[j] {
				points[i][j] = (points[i][j] - minValues[j]) / (maxValues[j] - minValues[j])
			} else {
				points[i][j] = 0.5 // Если все значения одинаковые
			}
		}
	}

	// Возвращаем нормализованные данные и диапазоны для денормализации
	ranges := [][]float64{minValues, maxValues}
	return points, ranges
}

// aggregateTransactionBehaviors агрегирует транзакции по пользователям и вычисляет поведенческие метрики
func (s *KMeansSegmentation) aggregateTransactionBehaviors(transactions []entities.Transaction) ([][]float64, []uuid.UUID) {
	userBehaviors := make(map[uuid.UUID][]float64)
	userTransactionCounts := make(map[uuid.UUID]int)

	// Предварительная обработка: сбор всех метрик
	for _, t := range transactions {
		if _, exists := userBehaviors[t.UserID]; !exists {
			userBehaviors[t.UserID] = make([]float64, 3) // [avg_amount, count, preferred_time]
		}

		userBehaviors[t.UserID][0] += t.Amount                    // Сумма транзакции
		userBehaviors[t.UserID][2] += float64(t.Timestamp.Hour()) // Время дня
		userTransactionCounts[t.UserID]++
	}

	// Формирование списка ID пользователей и финальная обработка метрик
	userIDs := make([]uuid.UUID, 0, len(userBehaviors))
	points := make([][]float64, 0, len(userBehaviors))

	for userID, behavior := range userBehaviors {
		count := userTransactionCounts[userID]
		if count > 0 {
			// Вычисляем средние значения
			behavior[0] /= float64(count) // Средняя сумма транзакции
			behavior[1] = float64(count)  // Количество транзакций
			behavior[2] /= float64(count) // Среднее время дня
		}

		userIDs = append(userIDs, userID)
		points = append(points, behavior)
	}

	return points, userIDs
}

// normalizeBehaviorData нормализует поведенческие данные для кластеризации
func (s *KMeansSegmentation) normalizeBehaviorData(points [][]float64) ([][]float64, [][]float64) {
	normalizedPoints := make([][]float64, len(points))
	for i := range points {
		normalizedPoints[i] = make([]float64, len(points[i]))
		copy(normalizedPoints[i], points[i])
	}

	// Находим минимальные и максимальные значения для каждой метрики
	minValues := make([]float64, 3)
	maxValues := make([]float64, 3)

	for j := 0; j < 3; j++ {
		minValues[j] = math.MaxFloat64
		maxValues[j] = -math.MaxFloat64

		for i := 0; i < len(points); i++ {
			if points[i][j] < minValues[j] {
				minValues[j] = points[i][j]
			}
			if points[i][j] > maxValues[j] {
				maxValues[j] = points[i][j]
			}
		}
	}

	// Нормализация
	for i := 0; i < len(normalizedPoints); i++ {
		for j := 0; j < 3; j++ {
			if maxValues[j] > minValues[j] {
				normalizedPoints[i][j] = (normalizedPoints[i][j] - minValues[j]) / (maxValues[j] - minValues[j])
			} else {
				normalizedPoints[i][j] = 0.5 // Если все значения одинаковые
			}
		}
	}

	// Возвращаем нормализованные данные и диапазоны для денормализации
	ranges := [][]float64{minValues, maxValues}
	return normalizedPoints, ranges
}

// kmeansClusteringPlusPlus выполняет кластеризацию методом K-means++
func (s *KMeansSegmentation) kmeansClusteringPlusPlus(points [][]float64) ([][]float64, []int) {
	// Инициализация генератора случайных чисел
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Выбор начальных центроидов методом K-means++
	centroids := s.kmeansppInitialization(points, r)

	// Выполнение кластеризации K-means
	assignments := make([]int, len(points))
	for iter := 0; iter < s.maxIter; iter++ {
		// Назначение точек ближайшим центроидам
		for i, p := range points {
			assignments[i] = s.closestCentroid(p, centroids)
		}

		// Обновление центроидов
		newCentroids := s.updateCentroids(points, assignments)

		// Проверка сходимости
		if s.converged(centroids, newCentroids) {
			break
		}

		centroids = newCentroids
	}

	return centroids, assignments
}

// kmeansppInitialization инициализирует центроиды методом K-means++
func (s *KMeansSegmentation) kmeansppInitialization(points [][]float64, r *rand.Rand) [][]float64 {
	centroids := make([][]float64, s.numClusters)

	// Выбираем первый центроид случайным образом
	firstIdx := r.Intn(len(points))
	centroids[0] = make([]float64, len(points[0]))
	copy(centroids[0], points[firstIdx])

	// Выбираем остальные центроиды
	for c := 1; c < s.numClusters; c++ {
		// Для каждой точки вычисляем минимальное расстояние до ближайшего центроида
		distances := make([]float64, len(points))
		sum := 0.0

		for i, p := range points {
			minDist := math.MaxFloat64
			for j := 0; j < c; j++ {
				dist := s.euclideanDistanceSquared(p, centroids[j])
				if dist < minDist {
					minDist = dist
				}
			}
			distances[i] = minDist
			sum += minDist
		}

		// Выбираем следующий центроид с вероятностью, пропорциональной квадрату расстояния
		target := r.Float64() * sum
		currentSum := 0.0
		selectedIdx := 0

		for i, dist := range distances {
			currentSum += dist
			if currentSum >= target {
				selectedIdx = i
				break
			}
		}

		centroids[c] = make([]float64, len(points[0]))
		copy(centroids[c], points[selectedIdx])
	}

	return centroids
}

// closestCentroid находит индекс ближайшего центроида для заданной точки
func (s *KMeansSegmentation) closestCentroid(point []float64, centroids [][]float64) int {
	minDist := math.MaxFloat64
	minIdx := 0

	for i, c := range centroids {
		dist := s.euclideanDistance(point, c)
		if dist < minDist {
			minDist = dist
			minIdx = i
		}
	}

	return minIdx
}

// updateCentroids обновляет центроиды на основе текущих назначений
func (s *KMeansSegmentation) updateCentroids(points [][]float64, assignments []int) [][]float64 {
	// Инициализация новых центроидов
	centroids := make([][]float64, s.numClusters)
	counts := make([]int, s.numClusters)

	for i := range centroids {
		centroids[i] = make([]float64, len(points[0]))
	}

	// Суммирование точек по кластерам
	for i, p := range points {
		cluster := assignments[i]
		for j := range p {
			centroids[cluster][j] += p[j]
		}
		counts[cluster]++
	}

	// Вычисление средних значений для каждого кластера
	for i := range centroids {
		if counts[i] > 0 {
			for j := range centroids[i] {
				centroids[i][j] /= float64(counts[i])
			}
		} else {
			// Если кластер пуст, выбираем случайную точку
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			randomIdx := r.Intn(len(points))
			copy(centroids[i], points[randomIdx])
		}
	}

	return centroids
}

// denormalizeCentroids конвертирует нормализованные центроиды обратно в исходный масштаб
func (s *KMeansSegmentation) denormalizeCentroids(centroids [][]float64, ranges [][]float64) [][]float64 {
	minValues := ranges[0]
	maxValues := ranges[1]

	denormalizedCentroids := make([][]float64, len(centroids))
	for i := range centroids {
		denormalizedCentroids[i] = make([]float64, len(centroids[i]))

		// Особая обработка для recency (первое значение), так как мы его инвертировали
		denormalizedCentroids[i][0] = maxValues[0] - (centroids[i][0] * (maxValues[0] - minValues[0]))

		// Денормализация других метрик
		for j := 1; j < len(centroids[i]); j++ {
			denormalizedCentroids[i][j] = centroids[i][j]*(maxValues[j]-minValues[j]) + minValues[j]
		}
	}

	return denormalizedCentroids
}

// analyzeRFMClusters анализирует характеристики кластеров RFM
func (s *KMeansSegmentation) analyzeRFMClusters(centroids [][]float64) []map[string]string {
	characteristics := make([]map[string]string, len(centroids))

	for i, c := range centroids {
		char := make(map[string]string)

		// Анализ Recency
		if c[0] < 30 {
			char["recency"] = "recent"
		} else if c[0] < 90 {
			char["recency"] = "moderate"
		} else {
			char["recency"] = "old"
		}

		// Анализ Frequency
		if c[1] > 10 {
			char["frequency"] = "high"
		} else if c[1] > 3 {
			char["frequency"] = "medium"
		} else {
			char["frequency"] = "low"
		}

		// Анализ Monetary
		if c[2] > 1000 {
			char["monetary"] = "high"
		} else if c[2] > 300 {
			char["monetary"] = "medium"
		} else {
			char["monetary"] = "low"
		}

		// Определение общего типа сегмента
		if char["recency"] == "recent" && char["frequency"] == "high" && char["monetary"] == "high" {
			char["type"] = "VIP"
		} else if char["recency"] == "recent" && char["frequency"] == "medium" && char["monetary"] == "medium" {
			char["type"] = "Loyal"
		} else if char["recency"] == "old" && char["frequency"] == "low" && char["monetary"] == "low" {
			char["type"] = "Churn Risk"
		} else if char["recency"] == "recent" && char["frequency"] == "low" {
			char["type"] = "New"
		} else if char["recency"] == "old" && char["frequency"] == "high" {
			char["type"] = "Dormant"
		} else {
			char["type"] = "Regular"
		}

		characteristics[i] = char
	}

	return characteristics
}

// analyzeBehaviorClusters анализирует характеристики поведенческих кластеров
func (s *KMeansSegmentation) analyzeBehaviorClusters(centroids [][]float64) []map[string]string {
	characteristics := make([]map[string]string, len(centroids))

	for i, c := range centroids {
		char := make(map[string]string)

		// Анализ средней суммы транзакции
		if c[0] > 500 {
			char["avg_amount"] = "high"
		} else if c[0] > 100 {
			char["avg_amount"] = "medium"
		} else {
			char["avg_amount"] = "low"
		}

		// Анализ частоты транзакций
		if c[1] > 10 {
			char["frequency"] = "high"
		} else if c[1] > 3 {
			char["frequency"] = "medium"
		} else {
			char["frequency"] = "low"
		}

		// Анализ предпочтительного времени
		hour := int(math.Round(c[2]))
		if hour >= 5 && hour < 12 {
			char["preferred_time"] = "morning"
		} else if hour >= 12 && hour < 17 {
			char["preferred_time"] = "day"
		} else if hour >= 17 && hour < 22 {
			char["preferred_time"] = "evening"
		} else {
			char["preferred_time"] = "night"
		}

		// Определение общего типа сегмента на основе поведения
		if char["avg_amount"] == "high" && char["frequency"] == "high" {
			char["type"] = "Premium"
		} else if char["avg_amount"] == "low" && char["frequency"] == "high" {
			char["type"] = "Frequent Saver"
		} else if char["avg_amount"] == "high" && char["frequency"] == "low" {
			char["type"] = "Big Spender"
		} else if char["frequency"] == "low" && char["preferred_time"] == "night" {
			char["type"] = "Night Owl"
		} else if char["frequency"] == "medium" && char["preferred_time"] == "morning" {
			char["type"] = "Morning Regular"
		} else {
			char["type"] = "Standard"
		}

		characteristics[i] = char
	}

	return characteristics
}

// generateRFMSegmentName генерирует название сегмента на основе его характеристик RFM
func (s *KMeansSegmentation) generateRFMSegmentName(characteristics map[string]string) string {
	if segmentType, ok := characteristics["type"]; ok {
		switch segmentType {
		case "VIP":
			return "VIP-клиенты"
		case "Loyal":
			return "Лояльные клиенты"
		case "Churn Risk":
			return "Риск оттока"
		case "New":
			return "Новые клиенты"
		case "Dormant":
			return "Дремлющие клиенты"
		default:
			return "Обычные клиенты"
		}
	}

	return "Сегмент " + string(rune('A'+rand.Intn(26)))
}

// generateBehaviorSegmentName генерирует название сегмента на основе поведенческих характеристик
func (s *KMeansSegmentation) generateBehaviorSegmentName(characteristics map[string]string) string {
	if segmentType, ok := characteristics["type"]; ok {
		switch segmentType {
		case "Premium":
			return "Премиум клиенты"
		case "Frequent Saver":
			return "Экономные постоянные"
		case "Big Spender":
			return "Крупные покупатели"
		case "Night Owl":
			return "Ночные клиенты"
		case "Morning Regular":
			return "Утренние постоянные"
		default:
			return "Стандартные клиенты"
		}
	}

	return "Поведенческий сегмент " + string(rune('A'+rand.Intn(26)))
}

// euclideanDistance вычисляет евклидово расстояние между двумя точками
func (s *KMeansSegmentation) euclideanDistance(a, b []float64) float64 {
	return math.Sqrt(s.euclideanDistanceSquared(a, b))
}

// euclideanDistanceSquared вычисляет квадрат евклидова расстояния между двумя точками
func (s *KMeansSegmentation) euclideanDistanceSquared(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return sum
}

// converged проверяет, сошлась ли кластеризация (центроиды больше не меняются)
func (s *KMeansSegmentation) converged(old, new [][]float64) bool {
	for i := range old {
		for j := range old[i] {
			if math.Abs(old[i][j]-new[i][j]) > s.epsilon {
				return false
			}
		}
	}
	return true
}