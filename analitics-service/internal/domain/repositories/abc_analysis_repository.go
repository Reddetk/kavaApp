package repositories

import (
	"context"
	"time"

	"analitics-service/internal/domain/entities"
)

// ABCAnalysisRepository определяет интерфейс для работы с результатами ABC-анализа
type ABCAnalysisRepository interface {
	// SaveAnalysisResult сохраняет результат ABC-анализа
	SaveAnalysisResult(ctx context.Context, result entities.ABCAnalysisResult) error

	// GetAnalysisResultByDate возвращает результат ABC-анализа на указанную дату
	GetAnalysisResultByDate(ctx context.Context, date time.Time) (entities.ABCAnalysisResult, error)

	// GetLatestAnalysisResult возвращает последний результат ABC-анализа
	GetLatestAnalysisResult(ctx context.Context) (entities.ABCAnalysisResult, error)

	// GetAnalysisHistory возвращает историю результатов ABC-анализа
	GetAnalysisHistory(ctx context.Context, startDate, endDate time.Time) ([]entities.ABCAnalysisResult, error)

	// SaveAnalysisCriteria сохраняет критерии ABC-анализа
	SaveAnalysisCriteria(ctx context.Context, criteria entities.ABCAnalysisCriteria) error

	// GetLatestAnalysisCriteria возвращает последние использованные критерии ABC-анализа
	GetLatestAnalysisCriteria(ctx context.Context) (entities.ABCAnalysisCriteria, error)
}
