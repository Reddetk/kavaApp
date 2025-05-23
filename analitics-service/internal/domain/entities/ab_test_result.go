// internal/domain/enteties/ab_test_result.go
package entities

import (
	"time"
)

// ABTestResult представляет результат A/B теста
type ABTestResult struct {
	TestID       string         `json:"test_id"`
	StartDate    time.Time      `json:"start_date"`
	EndDate      time.Time      `json:"end_date"`
	Description  string         `json:"description"`
	ControlGroup GroupStats     `json:"control_group"`
	TestGroup    TestGroupStats `json:"test_group"`
	Lift         float64        `json:"lift"`
	Significance float64        `json:"significance"`
	IsSignificant bool          `json:"is_significant"`
}