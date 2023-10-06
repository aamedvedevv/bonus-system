package service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

type OrderRepository interface {
	AddOrder(order domain.Order) error
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
	if checkOrderNumber(orderID) {
		return domain.ErrIncorrectOrder
	}

	userID, ok := ctx.Value("userID").(int64)
	fmt.Println(userID)
	if !ok {
		return fmt.Errorf("incorrect user id - %s", orderID)
	}

	order := domain.Order{
		OrderID:    orderID,
		Status:     domain.Registered,
		UploadedAt: time.Now(),
		Bonuses:    0,
		UserID:     userID,
	}

	return o.repo.AddOrder(order)
}
