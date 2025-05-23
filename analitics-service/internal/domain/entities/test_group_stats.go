// internal/domain/entities/test_group_stats.go
package entities

// TestGroupStats расширяет GroupStats для тестовой группы
type TestGroupStats struct {
	GroupStats
	DiscountPct float64 `json:"discount_pct,omitempty"`
	CouponUsed  bool    `json:"coupon_used"`
}
