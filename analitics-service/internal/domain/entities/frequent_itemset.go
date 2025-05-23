// internal/domain/entities/frequent_itemset.go
package entities

// FrequentItemset представляет частый набор товаров, найденный алгоритмом Apriori
type FrequentItemset struct {
	Items   []Item  `json:"items"`   // Товары в наборе
	Support float64 `json:"support"` // Поддержка набора (от 0 до 1)
	Count   int     `json:"count"`   // Количество транзакций, содержащих этот набор
}
