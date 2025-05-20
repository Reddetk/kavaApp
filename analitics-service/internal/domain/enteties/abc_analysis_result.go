// internal/domain/enteties/abc_analysis_result.go
package entities

// ABCAnalysisResult содержит результаты ABC-анализа
type ABCAnalysisResult struct {
	AnalysisMetadata
	ProductsSegmentation map[string]ProductFullSegmentation `json:"products_segmentation"`
	Summary              *ABCSegmentSummary                 `json:"summary"`
}