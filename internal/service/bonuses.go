package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

type BonusesRepository interface {
	Balance(userID int64) (float32, error)
	WithdrawBalance(userID int64) (float32, error)
	Withdraw(withdraw domain.Withdraw) error
}

type Bonuses struct {
	repo BonusesRepository
}

func NewBonuses(repo BonusesRepository) *Bonuses {
	return &Bonuses{
		repo: repo,
	}
}

func (b *Bonuses) Balance(ctx context.Context) (*domain.BalanceOutput, error) {

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	// узнаем баланс бонусов пользователя
	balanceUser, err := b.repo.Balance(userID)
	if err != nil {
		return nil, err
	}

	// узнаем баланс списанных бонусов пользователя
	balanceWithdraws, err := b.repo.WithdrawBalance(userID)
	if err != nil {
		return nil, err
	}

	var balance domain.BalanceOutput
	balance.Bonuses = balanceUser
	balance.Withdraw = balanceWithdraws

	return &balance, nil
}

func (b *Bonuses) Withdraw(ctx context.Context, withdraw domain.Withdraw) error {

	trimmedStr := strings.TrimSpace(withdraw.OrderID)
	if len(trimmedStr) == 0 {
		return domain.ErrIncorrectOrder
	}

	if !checkOrderNumber(withdraw.OrderID) {
		return domain.ErrIncorrectOrder
	}

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return errors.New("incorrect user id")
	}

	// проверка что внутри
	fmt.Println(withdraw)

	with := domain.Withdraw{
		OrderID:    withdraw.OrderID,
		Bonuses:    withdraw.Bonuses,
		UploadedAt: time.Now().Format(time.RFC3339),
		UserID:     userID,
	}

	// проверка что внутри
	fmt.Println(with)

	// узнаем баланс бонусов пользователя
	balanceUser, err := b.repo.Balance(userID)
	if err != nil {
		return err
	}

	// узнаем баланс списанных бонусов пользователя
	balanceWithdraws, err := b.repo.WithdrawBalance(userID)
	if err != nil {
		return err
	}

	// проверка для проведения списания бонусов
	sum := balanceUser - balanceWithdraws
	if sum >= with.Bonuses {
		b.repo.Withdraw(with)
	} else {
		return domain.ErrNoBonuses
	}

	return b.repo.Withdraw(with)
}
