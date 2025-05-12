// internal/infrastructure/services/markov_transition_service.go
package services

import (
	"errors"
	"user-service/internal/domain/entities"
	"user-service/internal/domain/services"

	"github.com/google/uuid"
)

// MarkovTransitionService implements the StateTransitionService interface
type MarkovTransitionService struct {
	// Хранилище состояний пользователей
	userStates map[uuid.UUID]string
	// Матрица переходов между состояниями
	transitionMatrix map[string]map[string]float64
	// Порог вероятности для определения оттока
	churnThreshold float64
}

// NewMarkovTransitionService creates a new instance of MarkovTransitionService
func NewMarkovTransitionService() services.StateTransitionService {
	return &MarkovTransitionService{
		userStates:       make(map[uuid.UUID]string),
		transitionMatrix: make(map[string]map[string]float64),
		churnThreshold:   0.5, // Значение по умолчанию, можно настроить через конфигурацию
	}
}

// BuildTransitionMatrix builds a transition matrix from user transactions
func (s *MarkovTransitionService) BuildTransitionMatrix(segment entities.Segment, transactions []entities.Transaction) (map[string]map[string]float64, error) {
	if len(transactions) == 0 {
		return nil, errors.New("no transactions provided")
	}

	// Определяем возможные состояния
	states := []string{"active", "at_risk", "churned"}

	// Инициализируем матрицу переходов
	matrix := make(map[string]map[string]float64)
	for _, fromState := range states {
		matrix[fromState] = make(map[string]float64)
		for _, toState := range states {
			matrix[fromState][toState] = 0.0
		}
	}

	// Заполняем матрицу на основе данных (упрощенная реализация)
	// В реальном приложении здесь будет более сложная логика анализа транзакций

	// Пример заполнения матрицы:
	matrix["active"]["active"] = 0.7
	matrix["active"]["at_risk"] = 0.2
	matrix["active"]["churned"] = 0.1

	matrix["at_risk"]["active"] = 0.3
	matrix["at_risk"]["at_risk"] = 0.4
	matrix["at_risk"]["churned"] = 0.3

	matrix["churned"]["active"] = 0.1
	matrix["churned"]["at_risk"] = 0.2
	matrix["churned"]["churned"] = 0.7

	// Сохраняем матрицу для использования в других методах
	s.transitionMatrix = matrix

	return matrix, nil
}

// PredictNextState predicts the next state for a user
func (s *MarkovTransitionService) PredictNextState(userID uuid.UUID, currentState string) (string, float64, error) {
	// Проверяем, что матрица переходов инициализирована
	if len(s.transitionMatrix) == 0 {
		return "", 0, errors.New("transition matrix not initialized")
	}

	// Проверяем, что текущее состояние существует в матрице
	stateTransitions, exists := s.transitionMatrix[currentState]
	if !exists {
		return "", 0, errors.New("invalid current state")
	}

	// Находим состояние с наибольшей вероятностью перехода
	var nextState string
	var maxProbability float64 = -1

	for state, probability := range stateTransitions {
		if probability > maxProbability {
			maxProbability = probability
			nextState = state
		}
	}

	return nextState, maxProbability, nil
}

// UpdateUserState updates the state of a user based on churn probability
func (s *MarkovTransitionService) UpdateUserState(userID uuid.UUID, churnProbability float64) error {
	// Определяем состояние на основе вероятности оттока
	var newState string

	if churnProbability >= s.churnThreshold {
		newState = "churned"
	} else if churnProbability >= s.churnThreshold/2 {
		newState = "at_risk"
	} else {
		newState = "active"
	}

	// Обновляем состояние пользователя
	s.userStates[userID] = newState

	return nil
}

// GetUserState returns the current state of a user
func (s *MarkovTransitionService) GetUserState(userID uuid.UUID) (string, error) {
	state, exists := s.userStates[userID]
	if !exists {
		return "unknown", errors.New("user state not found")
	}
	return state, nil
}

// SetChurnThreshold sets the threshold for determining churn
func (s *MarkovTransitionService) SetChurnThreshold(threshold float64) error {
	if threshold <= 0 || threshold >= 1 {
		return errors.New("threshold must be between 0 and 1")
	}
	s.churnThreshold = threshold
	return nil
}
