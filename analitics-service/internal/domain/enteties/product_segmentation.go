// internal/domain/entities/product_segmentation.go
package entities

import "time"

// ProductSegmentation содержит информацию о сегментации продукта
type ProductSegmentation struct {
	ProductID    string    `json:"product_id"`
	Segment      Segment   `json:"segment"`
	Score        float64   `json:"score"`
	AnalysisDate time.Time `json:"analysis_date"`
}
