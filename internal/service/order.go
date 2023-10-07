package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

type OrderRepository interface {
	AddOrder(order domain.Order) error
	GetAllOrders(userID int64) ([]domain.Order, error)
}

type Orders struct {
	repo OrderRepository
}

func NewOrders(repo OrderRepository) *Orders {
	return &Orders{
		repo: repo,
	}
}

func (o *Orders) AddOrderID(ctx context.Context, orderID string) error {

	trimmedStr := strings.TrimSpace(orderID)
	if len(trimmedStr) == 0 {
		return domain.ErrIncorrectOrder
	}

	if !checkOrderNumber(orderID) {
		return domain.ErrIncorrectOrder
	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return errors.New("incorrect user id")
	}

	order := domain.Order{
		OrderID:    orderID,
		Status:     domain.Registered,
		UploadedAt: time.Now().Format(time.RFC3339),
		Bonuses:    0,
		UserID:     userID,
	}

	return o.repo.AddOrder(order)
}

func (o *Orders) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	orders, err := o.repo.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
