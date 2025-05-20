// internal/domain/entities/retention_metrics.go
package entities

// RetentionMetrics представляет метрики удержания клиентов
type RetentionMetrics struct {
	Period             TimeRange `json:"period"`
	ChurnRate          float64   `json:"churn_rate"`
	RetentionRate      float64   `json:"retention_rate"`
	NewCustomers       int       `json:"new_customers"`
	LostCustomers      int       `json:"lost_customers"`
	ActiveCustomers    int       `json:"active_customers"`
	RepeatPurchaseRate float64   `json:"repeat_purchase_rate"`
}
