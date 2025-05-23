// internal/domain/entities/association_rule.go
package entities

// AssociationRule представляет ассоциативное правило между товарами
type AssociationRule struct {
	Antecedent []Item     `json:"antecedent"`  // Предшествующие товары (если эти товары в корзине)
	Consequent []Item     `json:"consequent"`  // Следующие товары (то эти товары также могут быть интересны)
	Support    float64    `json:"support"`     // Поддержка правила (от 0 до 1)
	Confidence float64    `json:"confidence"`  // Достоверность правила (от 0 до 1)
	Lift       float64    `json:"lift"`        // Показатель lift (> 1 означает положительную корреляцию)
	Items      []string   `json:"items"`       // Все товары в правиле (для удобства поиска)
	Categories []string   `json:"categories"`  // Категории товаров в правиле
	PriceRange [2]float64 `json:"price_range"` // Диапазон цен товаров в правиле
}
