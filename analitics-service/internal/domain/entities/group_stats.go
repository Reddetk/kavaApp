// internal/domain/entities/group_stats.go
package entities

// GroupStats представляет статистику по группе
type GroupStats struct {
	Size        int     `json:"size"`
	Conversion  float64 `json:"conversion"`
	AvgPurchase float64 `json:"avg_purchase"`
	Revenue     float64 `json:"revenue"`
}
