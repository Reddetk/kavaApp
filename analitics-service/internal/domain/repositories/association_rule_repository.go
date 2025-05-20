package repositories

import (
	"context"

	"analitics-service/internal/domain/enteties"
)

// AssociationRuleRepository определяет интерфейс для работы с ассоциативными правилами
type AssociationRuleRepository interface {
	// SaveRules сохраняет ассоциативные правила
	SaveRules(ctx context.Context, rules []entities.AssociationRule) error

	// GetRulesByProduct возвращает ассоциативные правила, связанные с конкретным продуктом
	GetRulesByProduct(ctx context.Context, productID string) ([]entities.AssociationRule, error)

	// GetRulesByCategory возвращает ассоциативные правила для продуктов определенной категории
	GetRulesByCategory(ctx context.Context, category string) ([]entities.AssociationRule, error)

	// GetRulesByConfidence возвращает ассоциативные правила с уверенностью выше указанного порога
	GetRulesByConfidence(ctx context.Context, minConfidence float64) ([]entities.AssociationRule, error)

	// GetRulesBySupport возвращает ассоциативные правила с поддержкой выше указанного порога
	GetRulesBySupport(ctx context.Context, minSupport float64) ([]entities.AssociationRule, error)

	// GetRulesByLift возвращает ассоциативные правила с подъемом выше указанного порога
	GetRulesByLift(ctx context.Context, minLift float64) ([]entities.AssociationRule, error)
}
