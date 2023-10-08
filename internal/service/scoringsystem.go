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

func (s *ScoringSystem) GetOrderStatus() (string, error) {
	return s.repo.GetOrderStatus()
}

func (s *ScoringSystem) UpdateOrder(order domain.ScoringSystem) error {
	return s.repo.UpdateOrder(order)
}
