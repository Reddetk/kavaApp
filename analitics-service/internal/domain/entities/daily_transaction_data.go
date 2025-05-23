// internal/domain/enteties/daily_transaction_data.go
package entities

import (
	"time"
)

// DailyTransactionData представляет агрегированные данные по транзакциям за день
type DailyTransactionData struct {
	Date         time.Time           `json:"date"`
	Sales        float64             `json:"sales"`
	TotalPrice   float64             `json:"total_price"`
	AvgDiscount  float64             `json:"avg_discount"`
	TotalTx      int                 `json:"total_tx"`
	DiscountedTx int                 `json:"discounted_tx"`
	IsHoliday    bool                `json:"is_holiday"`
	ProductCount int                 `json:"product_count"`
	ProductIDs   map[string]struct{} `json:"-"`
}