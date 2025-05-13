package utils

import (
	"math"
	"sort"
)

// Round округляет число до указанного количества знаков после запятой
func Round(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(value*ratio) / ratio
}

// RoundUp округляет число вверх до указанного количества знаков после запятой
func RoundUp(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Ceil(value*ratio) / ratio
}

// RoundDown округляет число вниз до указанного количества знаков после запятой
func RoundDown(value float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Floor(value*ratio) / ratio
}

// Mean вычисляет среднее арифметическое массива чисел
func Mean(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	
	return sum / float64(len(values))
}

// Median вычисляет медиану массива чисел
func Median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	// Создаем копию массива, чтобы не изменять оригинал
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)
	
	sort.Float64s(valuesCopy)
	
	length := len(valuesCopy)
	if length%2 == 0 {
		// Если четное количество элементов, берем среднее двух средних
		return (valuesCopy[length/2-1] + valuesCopy[length/2]) / 2
	}
	
	// Если нечетное количество элементов, берем средний элемент
	return valuesCopy[length/2]
}

// Mode вычисляет моду массива чисел (наиболее часто встречающееся значение)
func Mode(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	// Подсчитываем частоту каждого значения
	frequency := make(map[float64]int)
	for _, v := range values {
		frequency[v]++
	}
	
	// Находим значение с наибольшей частотой
	maxFreq := 0
	var mode float64
	for value, freq := range frequency {
		if freq > maxFreq {
			maxFreq = freq
			mode = value
		}
	}
	
	return mode
}

// StandardDeviation вычисляет стандартное отклонение массива чисел
func StandardDeviation(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	
	mean := Mean(values)
	
	// Вычисляем сумму квадратов отклонений от среднего
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	
	// Вычисляем дисперсию и стандартное отклонение
	variance := sumSquaredDiff / float64(len(values)-1) // Используем n-1 для несмещенной оценки
	return math.Sqrt(variance)
}

// Variance вычисляет дисперсию массива чисел
func Variance(values []float64) float64 {
	if len(values) < 2 {
		return 0
	}
	
	mean := Mean(values)
	
	// Вычисляем сумму квадратов отклонений от среднего
	sumSquaredDiff := 0.0
	for _, v := range values {
		diff := v - mean
		sumSquaredDiff += diff * diff
	}
	
	// Вычисляем дисперсию
	return sumSquaredDiff / float64(len(values)-1) // Используем n-1 для несмещенной оценки
}

// Min возвращает минимальное значение в массиве
func Min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	
	return min
}

// Max возвращает максимальное значение в массиве
func Max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	
	return max
}

// Percentile вычисляет указанный процентиль массива чисел
func Percentile(values []float64, percentile float64) float64 {
	if len(values) == 0 {
		return 0
	}
	
	// Проверяем, что процентиль в диапазоне [0, 100]
	if percentile < 0 {
		percentile = 0
	} else if percentile > 100 {
		percentile = 100
	}
	
	// Создаем копию массива и сортируем
	valuesCopy := make([]float64, len(values))
	copy(valuesCopy, values)
	sort.Float64s(valuesCopy)
	
	// Вычисляем индекс для процентиля
	index := (percentile / 100) * float64(len(valuesCopy)-1)
	
	// Если индекс целое число, возвращаем значение по этому индексу
	if index == float64(int(index)) {
		return valuesCopy[int(index)]
	}
	
	// Иначе интерполируем между двумя ближайшими значениями
	lower := valuesCopy[int(math.Floor(index))]
	upper := valuesCopy[int(math.Ceil(index))]
	fraction := index - math.Floor(index)
	
	return lower + fraction*(upper-lower)
}

// IQR вычисляет межквартильный размах (разница между 75-м и 25-м процентилями)
func IQR(values []float64) float64 {
	return Percentile(values, 75) - Percentile(values, 25)
}

// Normalize нормализует значение в диапазоне [min, max] к диапазону [0, 1]
func Normalize(value, min, max float64) float64 {
	if max == min {
		return 0.5 // Если min и max равны, возвращаем середину диапазона
	}
	
	normalized := (value - min) / (max - min)
	
	// Ограничиваем значение диапазоном [0, 1]
	if normalized < 0 {
		return 0
	} else if normalized > 1 {
		return 1
	}
	
	return normalized
}

// NormalizeArray нормализует массив значений к диапазону [0, 1]
func NormalizeArray(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{}
	}
	
	min := Min(values)
	max := Max(values)
	
	normalized := make([]float64, len(values))
	for i, v := range values {
		normalized[i] = Normalize(v, min, max)
	}
	
	return normalized
}

// ZScore вычисляет z-оценку (количество стандартных отклонений от среднего)
func ZScore(value float64, mean float64, stdDev float64) float64 {
	if stdDev == 0 {
		return 0 // Избегаем деления на ноль
	}
	
	return (value - mean) / stdDev
}

// ZScoreArray вычисляет z-оценки для массива значений
func ZScoreArray(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{}
	}
	
	mean := Mean(values)
	stdDev := StandardDeviation(values)
	
	zScores := make([]float64, len(values))
	for i, v := range values {
		zScores[i] = ZScore(v, mean, stdDev)
	}
	
	return zScores
}

// IsOutlier проверяет, является ли значение выбросом по методу IQR
// Выбросом считается значение, которое меньше Q1 - 1.5*IQR или больше Q3 + 1.5*IQR
func IsOutlier(value float64, values []float64) bool {
	q1 := Percentile(values, 25)
	q3 := Percentile(values, 75)
	iqr := q3 - q1
	
	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr
	
	return value < lowerBound || value > upperBound
}

// RemoveOutliers удаляет выбросы из массива значений по методу IQR
func RemoveOutliers(values []float64) []float64 {
	if len(values) == 0 {
		return []float64{}
	}
	
	q1 := Percentile(values, 25)
	q3 := Percentile(values, 75)
	iqr := q3 - q1
	
	lowerBound := q1 - 1.5*iqr
	upperBound := q3 + 1.5*iqr
	
	result := make([]float64, 0, len(values))
	for _, v := range values {
		if v >= lowerBound && v <= upperBound {
			result = append(result, v)
		}
	}
	
	return result
}