// internal/domain/entities/product_full_segmentation.go
package entities

// ProductFullSegmentation содержит детальную информацию о сегментации продукта
type ProductFullSegmentation struct {
	ProductID       string  `json:"product_id"`
	RevenueSegment  Segment `json:"revenue_segment"`
	QuantitySegment Segment `json:"quantity_segment"`
	ProfitSegment   Segment `json:"profit_segment"`
	FinalSegment    Segment `json:"final_segment"`
	Score           float64 `json:"score"`
}
