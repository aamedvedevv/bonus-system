package service

import (
	"context"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

type ScoringSystemRepository interface {
	GetOrderStatus(ctx context.Context) ([]string, error)
	UpdateOrder(ctx context.Context, order domain.ScoringSystem) error
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
func (s *ScoringSystem) GetOrderStatus(ctx context.Context) ([]string, error) {
	return s.repo.GetOrderStatus(ctx)
}

// UpdateOrder обновляет данные заказа. Начисляет бонусы и меняет статус.
func (s *ScoringSystem) UpdateOrder(ctx context.Context, order domain.ScoringSystem) error {
	return s.repo.UpdateOrder(ctx, order)
}
