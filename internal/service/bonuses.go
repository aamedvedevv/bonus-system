package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AlexCorn999/bonus-system/internal/domain"
	"github.com/shopspring/decimal"
)

type BonusesRepository interface {
	Balance(ctx context.Context, userID int64) (float32, error)
	WithdrawBalance(ctx context.Context, userID int64) (float32, error)
	Withdraw(ctx context.Context, withdraw domain.Withdraw) error
	Withdrawals(ctx context.Context, userID int64) ([]domain.Withdraw, error)
}

type Bonuses struct {
	repo BonusesRepository
}

func NewBonuses(repo BonusesRepository) *Bonuses {
	return &Bonuses{
		repo: repo,
	}
}

// Balance выводит сумму баллов лояльности и использованных за весь период регистрации баллов пользователя.
func (b *Bonuses) Balance(ctx context.Context) (*domain.BalanceOutput, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	// узнаем баланс бонусов пользователя
	balanceUser, err := b.repo.Balance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// узнаем баланс списанных бонусов пользователя
	balanceWithdraws, err := b.repo.WithdrawBalance(ctx, userID)
	if err != nil {
		return nil, err
	}

	// чтобы узнать баланс пользователя вычитаем кол-во использованных бонусов
	newBalanceUser := decimal.NewFromFloat32(balanceUser).Sub(decimal.NewFromFloat32(balanceWithdraws))

	var balance domain.BalanceOutput
	balance.Bonuses = float32(newBalanceUser.InexactFloat64())
	balance.Withdraw = balanceWithdraws
	return &balance, nil
}

// Withdraw реализует списание бонусов пользователя в учет суммы нового заказа.
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

	with := domain.Withdraw{
		OrderID:    withdraw.OrderID,
		Bonuses:    withdraw.Bonuses,
		UploadedAt: time.Now().Format(time.RFC3339),
		UserID:     userID,
	}

	// узнаем баланс бонусов пользователя
	balanceUser, err := b.repo.Balance(ctx, userID)
	if err != nil {
		return err
	}

	// узнаем баланс списанных бонусов пользователя
	balanceWithdraws, err := b.repo.WithdrawBalance(ctx, userID)
	if err != nil {
		return err
	}

	// проверка для проведения списания бонусов
	sum := decimal.NewFromFloat32(balanceUser).Sub(decimal.NewFromFloat32(balanceWithdraws))
	if sum.Cmp(decimal.NewFromFloat32(with.Bonuses)) >= 0 {
		return b.repo.Withdraw(ctx, with)
	} else {
		return domain.ErrNoBonuses
	}

}

// Withdrawals выводит отсортированный по дате список списаний бонусов пользователя. Не больше 10 последних записей.
func (b *Bonuses) Withdrawals(ctx context.Context) ([]domain.Withdraw, error) {
	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	withdrawals, err := b.repo.Withdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Пагинация
	if len(withdrawals) <= 10 {
		return withdrawals, nil
	}

	startIndex := len(withdrawals) - 10
	endIndex := len(withdrawals)
	return withdrawals[startIndex:endIndex], nil
}
