package services

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"time"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/repositories"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
	"gonum.org/v1/gonum/mat"
)

// CoxSurvivalAnalysis реализует модель пропорциональных рисков Кокса
// для прогнозирования оттока пользователей
type CoxSurvivalAnalysis struct {
	// Коэффициенты модели для различных предикторов
	coefficients map[string]float64

	// Базовая функция кумулятивного риска (baseline hazard)
	baselineCumulativeHazard []float64

	// Временные точки для базовой функции риска
	timePoints []float64

	// Кэш кривых выживания для пользователей
	survivalCurves map[uuid.UUID][]float64

	// Метки события (1 - отток произошел, 0 - цензурированные данные)
	eventStatus map[uuid.UUID]int

	// Максимальное число итераций для оптимизации
	maxIter int

	// Порог сходимости для оптимизации
	tolerance float64

	// Регуляризационный параметр (L2 - ridge regression)
	lambda float64

	userRep repositories.UserRepository
}

// NewCoxSurvivalAnalysis создает новый экземпляр сервиса анализа выживаемости
func NewCoxSurvivalAnalysis() services.SurvivalAnalysisService {
	return &CoxSurvivalAnalysis{
		coefficients:             make(map[string]float64),
		baselineCumulativeHazard: []float64{},
		timePoints:               []float64{},
		survivalCurves:           make(map[uuid.UUID][]float64),
		eventStatus:              make(map[uuid.UUID]int),
		maxIter:                  100,
		tolerance:                1e-6,
		lambda:                   0.01, // L2 регуляризация
	}
}

// BuildCoxModel строит модель Кокса на основе данных пользователей
func (s *CoxSurvivalAnalysis) BuildCoxModel(segment entities.Segment, users []entities.UserMetrics) error {
	if len(users) == 0 {
		return errors.New("no users provided for model building")
	}

	// Подготовка данных для обучения модели
	features, timeToEvent, eventOccurred, err := s.prepareTrainingData(users)
	if err != nil {
		return fmt.Errorf("failed to prepare training data: %w", err)
	}

	// Получаем список предикторов (названия признаков)
	predictors := []string{"recency", "frequency", "monetary", "age", "sessionCount", "avgSessionDuration"}

	// Обучаем модель с использованием метода частичного правдоподобия
	err = s.fitCoxModel(features, timeToEvent, eventOccurred, predictors)
	if err != nil {
		return fmt.Errorf("failed to fit Cox model: %w", err)
	}

	// Вычисляем базовую функцию риска
	err = s.estimateBaselineHazard(features, timeToEvent, eventOccurred)
	if err != nil {
		return fmt.Errorf("failed to estimate baseline hazard: %w", err)
	}

	// Для каждого пользователя вычисляем и сохраняем кривую выживания
	for _, user := range users {
		// Сохраняем информацию о событии
		if user.Churned {
			s.eventStatus[user.UserID] = 1
		} else {
			s.eventStatus[user.UserID] = 0
		}

		// Вычисляем линейный предиктор для пользователя
		linearPredictor := s.calculateLinearPredictor(user)

		// Вычисляем кривую выживания
		survivalCurve := s.calculateSurvivalCurve(linearPredictor)

		// Сохраняем кривую выживания
		s.survivalCurves[user.UserID] = survivalCurve
	}

	return nil
}

// prepareTrainingData подготавливает данные для обучения модели Кокса
func (s *CoxSurvivalAnalysis) prepareTrainingData(users []entities.UserMetrics) ([][]float64, []float64, []int, error) {
	if len(users) == 0 {
		return nil, nil, nil, errors.New("no users provided")
	}
	// Матрица признаков
	features := make([][]float64, len(users))

	// Вектор времени до события (или цензурирования)
	timeToEvent := make([]float64, len(users))

	// Вектор индикатора события (1 - отток произошел, 0 - цензурированные данные)
	eventOccurred := make([]int, len(users))

	// Для каждого пользователя извлекаем признаки и метки
	for i, user := range users {

		// Формируем вектор признаков для пользователя
		features[i] = []float64{
			float64(user.Recency),
			float64(user.Frequency),
			user.Monetary,
			float64(user.Age),
			float64(user.SessionCount),
			user.AvgSessionDuration,
		}
		usrEnt, err := s.userRep.Get(context.Background(), user.UserID)
		DaysSinceRegistration := 5.0
		if err != nil {
			DaysSinceRegistration = float64(usrEnt.RegistrationDate.Day() - time.Now().Day())
		}
		// Вычисляем время до события (в днях)
		timeToEvent[i] = float64(DaysSinceRegistration)

		// Устанавливаем индикатор события
		if user.Churned {
			eventOccurred[i] = 1
		} else {
			eventOccurred[i] = 0
		}
	}

	return features, timeToEvent, eventOccurred, nil
}

// fitCoxModel обучает модель Кокса методом частичного правдоподобия
func (s *CoxSurvivalAnalysis) fitCoxModel(features [][]float64, time []float64, status []int, predictors []string) error {
	nSamples := len(features)
	nFeatures := len(predictors)

	if nSamples == 0 || nFeatures == 0 {
		return errors.New("empty features or predictors")
	}

	// Создаем матрицу признаков для библиотеки gonum
	x := mat.NewDense(nSamples, nFeatures, nil)
	for i := 0; i < nSamples; i++ {
		// Проверяем длину вектора признаков
		if len(features[i]) != nFeatures {
			return fmt.Errorf("inconsistent feature vector length at index %d", i)
		}

		// Заполняем матрицу признаков
		for j := 0; j < nFeatures; j++ {
			x.Set(i, j, features[i][j])
		}
	}

	// Инициализируем коэффициенты модели нулями
	beta := make([]float64, nFeatures)

	// Создаем структуру для хранения индексов рисковых наборов
	riskSets := s.createRiskSets(time, status)

	// Градиентный спуск для оптимизации частичного правдоподобия
	converged := false
	iter := 0

	for iter < s.maxIter && !converged {
		// Вычисляем градиент и гессиан логарифма частичного правдоподобия
		gradient, hessian := s.calculateGradientAndHessian(x, status, riskSets, beta)

		// Добавляем регуляризацию
		for j := 0; j < nFeatures; j++ {
			gradient[j] -= s.lambda * beta[j]
			hessian.Set(j, j, hessian.At(j, j)+s.lambda)
		}

		// Решаем систему линейных уравнений для обновления коэффициентов
		// H * delta_beta = gradient
		var deltaInv mat.Dense
		err := deltaInv.Inverse(hessian)
		if err != nil {
			// Если гессиан вырожден, используем псевдоинверсию или регуляризацию
			for j := 0; j < nFeatures; j++ {
				hessian.Set(j, j, hessian.At(j, j)+0.1) // Увеличиваем регуляризацию
			}
			err = deltaInv.Inverse(hessian)
			if err != nil {
				return fmt.Errorf("failed to invert Hessian: %w", err)
			}
		}

		// Вычисляем шаг обновления
		delta := make([]float64, nFeatures)
		for i := 0; i < nFeatures; i++ {
			for j := 0; j < nFeatures; j++ {
				delta[i] += deltaInv.At(i, j) * gradient[j]
			}
		}

		// Обновляем коэффициенты
		newBeta := make([]float64, nFeatures)
		for j := 0; j < nFeatures; j++ {
			newBeta[j] = beta[j] + delta[j]
		}

		// Проверяем сходимость
		maxChange := 0.0
		for j := 0; j < nFeatures; j++ {
			change := math.Abs(newBeta[j] - beta[j])
			if change > maxChange {
				maxChange = change
			}
		}

		beta = newBeta

		// Проверка условия сходимости
		if maxChange < s.tolerance {
			converged = true
		}

		iter++
	}

	// Сохраняем оптимальные коэффициенты
	for i, name := range predictors {
		s.coefficients[name] = beta[i]
	}

	return nil
}

// createRiskSets создает наборы риска для всех временных точек
func (s *CoxSurvivalAnalysis) createRiskSets(time []float64, status []int) map[float64][]int {
	// Сортируем уникальные времена событий
	uniqueTimes := make([]float64, 0)
	timeSet := make(map[float64]bool)

	for i, t := range time {
		if status[i] == 1 { // Учитываем только события (не цензурированные данные)
			if _, exists := timeSet[t]; !exists {
				timeSet[t] = true
				uniqueTimes = append(uniqueTimes, t)
			}
		}
	}

	sort.Float64s(uniqueTimes)
	s.timePoints = uniqueTimes

	// Для каждой временной точки определяем набор риска
	// (индексы субъектов, находящихся под риском в этот момент)
	riskSets := make(map[float64][]int)

	for _, t := range uniqueTimes {
		riskSet := make([]int, 0)
		for i, subjectTime := range time {
			if subjectTime >= t { // Субъект находится под риском, если его время >= t
				riskSet = append(riskSet, i)
			}
		}
		riskSets[t] = riskSet
	}

	return riskSets
}

// calculateGradientAndHessian вычисляет градиент и гессиан логарифма частичного правдоподобия
func (s *CoxSurvivalAnalysis) calculateGradientAndHessian(x *mat.Dense, status []int, riskSets map[float64][]int, beta []float64) ([]float64, *mat.Dense) {
	nSamples, nFeatures := x.Dims()

	// Инициализируем градиент и гессиан
	gradient := make([]float64, nFeatures)
	hessian := mat.NewDense(nFeatures, nFeatures, nil)

	// Вычисляем линейные предикторы для всех наблюдений
	linearPredictors := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		for j := 0; j < nFeatures; j++ {
			linearPredictors[i] += x.At(i, j) * beta[j]
		}
	}

	// Для каждой временной точки с событием
	for t, riskSet := range riskSets {
		// Вычисляем первый и второй моменты для набора риска
		zSum := make([]float64, nFeatures)
		zSumExp := 0.0

		for _, idx := range riskSet {
			expLP := math.Exp(linearPredictors[idx])
			zSumExp += expLP

			for j := 0; j < nFeatures; j++ {
				zSum[j] += x.At(idx, j) * expLP
			}
		}

		// Для события в этой временной точке
		eventIndices := make([]int, 0)
		for i, time := range s.timePoints {
			if time == t && status[i] == 1 {
				eventIndices = append(eventIndices, i)
			}
		}

		// Вычисляем вклад в градиент и гессиан
		for _, i := range eventIndices {
			// Градиент
			for j := 0; j < nFeatures; j++ {
				gradient[j] += x.At(i, j) - zSum[j]/zSumExp
			}

			// Гессиан
			for j := 0; j < nFeatures; j++ {
				for k := 0; k < nFeatures; k++ {
					zjk := 0.0
					for _, idx := range riskSet {
						expLP := math.Exp(linearPredictors[idx])
						zjk += x.At(idx, j) * x.At(idx, k) * expLP
					}
					hessValue := -(zjk/zSumExp - (zSum[j]*zSum[k])/(zSumExp*zSumExp))
					hessian.Set(j, k, hessian.At(j, k)+hessValue)
				}
			}
		}
	}

	return gradient, hessian
}

// estimateBaselineHazard оценивает базовую функцию риска
func (s *CoxSurvivalAnalysis) estimateBaselineHazard(features [][]float64, time []float64, status []int) error {
	nSamples := len(features)
	if nSamples == 0 {
		return errors.New("no samples provided")
	}

	// Вычисляем линейные предикторы для всех наблюдений
	linearPredictors := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		linearPredictors[i] = 0.0
		for name, coef := range s.coefficients {
			switch name {
			case "recency":
				linearPredictors[i] += coef * features[i][0]
			case "frequency":
				linearPredictors[i] += coef * features[i][1]
			case "monetary":
				linearPredictors[i] += coef * features[i][2]
			case "age":
				linearPredictors[i] += coef * features[i][3]
			case "sessionCount":
				linearPredictors[i] += coef * features[i][4]
			case "avgSessionDuration":
				linearPredictors[i] += coef * features[i][5]
			}
		}
	}

	// Экспоненцируем линейные предикторы
	expLP := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		expLP[i] = math.Exp(linearPredictors[i])
	}

	// Создаем временные точки и сортируем их
	timeEvents := make([]struct {
		time   float64
		status int
		index  int
	}, nSamples)

	for i := 0; i < nSamples; i++ {
		timeEvents[i] = struct {
			time   float64
			status int
			index  int
		}{time[i], status[i], i}
	}

	// Сортируем по времени
	sort.Slice(timeEvents, func(i, j int) bool {
		return timeEvents[i].time < timeEvents[j].time
	})

	// Вычисляем базовую функцию риска
	baselineHazard := make([]float64, 0)
	cumulativeHazard := make([]float64, 0)
	eventTimes := make([]float64, 0)

	// Для каждой временной точки с событием
	for _, te := range timeEvents {
		if te.status == 1 { // Если произошло событие
			// Находим всех субъектов под риском в этот момент
			riskSet := make([]int, 0)
			riskSum := 0.0

			for i := 0; i < nSamples; i++ {
				if time[i] >= te.time {
					riskSet = append(riskSet, i)
					riskSum += expLP[i]
				}
			}

			// Вычисляем приращение базовой функции риска
			if riskSum > 0 {
				h0 := 1.0 / riskSum
				baselineHazard = append(baselineHazard, h0)

				// Обновляем кумулятивную функцию риска
				cumH0 := 0.0
				if len(cumulativeHazard) > 0 {
					cumH0 = cumulativeHazard[len(cumulativeHazard)-1]
				}
				cumulativeHazard = append(cumulativeHazard, cumH0+h0)
				eventTimes = append(eventTimes, te.time)
			}
		}
	}

	// Сохраняем результаты
	s.baselineCumulativeHazard = cumulativeHazard
	s.timePoints = eventTimes

	return nil
}

// PredictChurnProbability предсказывает вероятность оттока пользователя
func (s *CoxSurvivalAnalysis) PredictChurnProbability(userID uuid.UUID, metrics entities.UserMetrics) (float64, error) {
	// Проверяем, есть ли уже кривая выживания для пользователя
	survival, ok := s.survivalCurves[userID]
	if !ok {
		// Если нет, вычисляем ее
		linearPredictor := s.calculateLinearPredictor(metrics)
		survival = s.calculateSurvivalCurve(linearPredictor)
		s.survivalCurves[userID] = survival
	}

	// Вероятность оттока - это 1 минус вероятность выживания в последней точке
	if len(survival) == 0 {
		return 0.5, nil // Значение по умолчанию, если кривая пуста
	}

	last := survival[len(survival)-1]
	return 1 - last, nil
}

// PredictTimeToEvent предсказывает ожидаемое время до события (оттока)
func (s *CoxSurvivalAnalysis) PredictTimeToEvent(userID uuid.UUID, metrics entities.UserMetrics) (float64, error) {
	// Проверяем, есть ли уже кривая выживания для пользователя
	survival, ok := s.survivalCurves[userID]
	if !ok {
		// Если нет, вычисляем ее
		linearPredictor := s.calculateLinearPredictor(metrics)
		survival = s.calculateSurvivalCurve(linearPredictor)
		s.survivalCurves[userID] = survival
	}

	if len(survival) == 0 || len(s.timePoints) == 0 {
		return 0, errors.New("insufficient data for prediction")
	}

	// Ожидаемое время до события вычисляется как интеграл функции выживания
	// Для дискретного случая используем численное интегрирование
	var expectedTime float64

	// Используем метод трапеций для интегрирования
	for i := 0; i < len(survival)-1; i++ {
		if i < len(s.timePoints)-1 {
			dt := s.timePoints[i+1] - s.timePoints[i]
			avgSurvival := (survival[i] + survival[i+1]) / 2.0
			expectedTime += avgSurvival * dt
		}
	}

	return expectedTime, nil
}

// GetChurnRiskFactors возвращает факторы риска оттока для пользователя
func (s *CoxSurvivalAnalysis) GetChurnRiskFactors(userID uuid.UUID, metrics entities.UserMetrics) (map[string]float64, error) {
	// Проверяем, что модель обучена
	if len(s.coefficients) == 0 {
		return nil, errors.New("model is not trained")
	}

	// Вычисляем вклад каждого фактора в риск оттока
	factors := make(map[string]float64)

	// Получаем значения признаков пользователя
	userFeatures := map[string]float64{
		"recency":            float64(metrics.Recency),
		"frequency":          float64(metrics.Frequency),
		"monetary":           metrics.Monetary,
		"age":                float64(metrics.Age),
		"sessionCount":       float64(metrics.SessionCount),
		"avgSessionDuration": metrics.AvgSessionDuration,
	}

	// Для каждого коэффициента вычисляем вклад
	for name, coef := range s.coefficients {
		if value, ok := userFeatures[name]; ok {
			// Вклад фактора = коэффициент * значение признака
			contribution := coef * value
			factors[name] = contribution
		}
	}

	// Сортируем факторы по абсолютному значению вклада
	type factorContribution struct {
		name         string
		contribution float64
	}

	sortedFactors := make([]factorContribution, 0, len(factors))
	for name, contribution := range factors {
		sortedFactors = append(sortedFactors, factorContribution{name, contribution})
	}

	sort.Slice(sortedFactors, func(i, j int) bool {
		return math.Abs(sortedFactors[i].contribution) > math.Abs(sortedFactors[j].contribution)
	})

	// Формируем результат
	result := make(map[string]float64)
	for _, factor := range sortedFactors {
		result[factor.name] = factor.contribution
	}

	return result, nil
}

// calculateLinearPredictor вычисляет линейный предиктор модели Кокса
func (s *CoxSurvivalAnalysis) calculateLinearPredictor(metrics entities.UserMetrics) float64 {
	// Линейный предиктор - это сумма произведений коэффициентов на значения предикторов
	linearPredictor := 0.0

	// Применяем коэффициенты к метрикам пользователя
	if coef, ok := s.coefficients["recency"]; ok {
		linearPredictor += coef * float64(metrics.Recency)
	}

	if coef, ok := s.coefficients["frequency"]; ok {
		linearPredictor += coef * float64(metrics.Frequency)
	}

	if coef, ok := s.coefficients["monetary"]; ok {
		linearPredictor += coef * metrics.Monetary
	}

	if coef, ok := s.coefficients["age"]; ok {
		linearPredictor += coef * float64(metrics.Age)
	}

	if coef, ok := s.coefficients["sessionCount"]; ok {
		linearPredictor += coef * float64(metrics.SessionCount)
	}

	if coef, ok := s.coefficients["avgSessionDuration"]; ok {
		linearPredictor += coef * metrics.AvgSessionDuration
	}

	return linearPredictor
}

// calculateSurvivalCurve вычисляет кривую выживания для заданного линейного предиктора
func (s *CoxSurvivalAnalysis) calculateSurvivalCurve(linearPredictor float64) []float64 {
	// Кривая выживания S(t) = exp(-H0(t) * exp(LP))
	// где H0(t) - кумулятивная базовая функция риска, LP - линейный предиктор

	// Вычисляем множитель exp(LP)
	expLP := math.Exp(linearPredictor)

	// Вычисляем кривую выживания
	survivalCurve := make([]float64, len(s.baselineCumulativeHazard))
	for i, hazard := range s.baselineCumulativeHazard {
		survivalCurve[i] = math.Exp(-hazard * expLP)
	}

	return survivalCurve
}

// SaveModel сохраняет модель
func (s *CoxSurvivalAnalysis) SaveModel() ([]byte, error) {
	// Здесь может быть реализация сохранения модели в байтовый буфер
	// Например, через gob или JSON
	return nil, errors.New("not implemented")
}

// LoadModel загружает модель
func (s *CoxSurvivalAnalysis) LoadModel(data []byte) error {
	// Здесь может быть реализация загрузки модели из байтового буфера
	return errors.New("not implemented")
}
