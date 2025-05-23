// internal/domain/entities/product_recommendation.go
package entities

// ProductRecommendation представляет рекомендацию товара на основе ассоциативных правил
type ProductRecommendation struct {
	Product Product `json:"product"` // Рекомендуемый товар
	Score   float64 `json:"score"`   // Оценка релевантности рекомендации (обычно confidence)
	Lift    float64 `json:"lift"`    // Показатель lift для рекомендации
	Support float64 `json:"support"` // Поддержка правила, на основе которого сделана рекомендация
}
