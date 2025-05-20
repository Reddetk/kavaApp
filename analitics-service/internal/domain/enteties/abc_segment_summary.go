// internal/domain/enteties/abc_segment_summary.go
package entities

// ABCSegmentSummary содержит сводную информацию по сегментам A, B и C
type ABCSegmentSummary struct {
	SegmentCounts      map[Segment]int     `json:"segment_counts"`
	SegmentPercentages map[Segment]float64 `json:"segment_percentages"`
}
