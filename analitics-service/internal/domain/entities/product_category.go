// internal/domain/entities/product_category.go
package entities

// ProductCategory представляет категорию товаров по ABC-анализу
type ProductCategory struct {
	CategoryName         string  `json:"category_name"`
	TotalRevenue         float64 `json:"total_revenue"`
	TotalSales           int     `json:"total_sales"`
	ABCCategory          Segment `json:"abc_category"`
	Percentage           float64 `json:"percentage"`
	CumulativePercentage float64 `json:"cumulative_percentage"`
}
