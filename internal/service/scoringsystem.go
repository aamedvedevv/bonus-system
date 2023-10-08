package service

import (
	"github.com/AlexCorn999/bonus-system/internal/domain"
)

type ScoringSystemRepository interface {
	GetOrderStatus() (string, error)
	UpdateOrder(order domain.ScoringSystem) error
}

type ScoringSystem struct {
	repo ScoringSystemRepository
}

func NewScoringSystem(repo ScoringSystemRepository) *ScoringSystem {
	return &ScoringSystem{
		repo: repo,
	}
}

// GetOrderStatus получает orderID если его статус не PROCESSED или INVALID.
func (s *ScoringSystem) GetOrderStatus() (string, error) {
	return s.repo.GetOrderStatus()
}

// UpdateOrder обновляет данные заказа. Начисляет бонусы и меняет статус.
func (s *ScoringSystem) UpdateOrder(order domain.ScoringSystem) error {
	return s.repo.UpdateOrder(order)
}
