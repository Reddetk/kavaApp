// internal/domain/enteties/association_rule.go
package entities

// AssociationRule представляет ассоциативное правило между товарами
type AssociationRule struct {
	Antecedent []string   `json:"antecedent"`
	Consequent []string   `json:"consequent"`
	Support    float64    `json:"support"`
	Confidence float64    `json:"confidence"`
	Lift       float64    `json:"lift"`
	Items      []string   `json:"items"`
	Categories []string   `json:"categories"`
	PriceRange [2]float64 `json:"price_range"`
}