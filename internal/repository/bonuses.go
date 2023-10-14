package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// Withdraw добавляет списания бонусов пользователя.
func (s *Storage) Withdraw(withdraw domain.Withdraw) error {
	result, err := s.db.Exec("INSERT INTO withdrawals (order_id, bonuses, uploaded_at, user_id) values ($1, $2, $3, $4) on conflict (order_id) do nothing",
		withdraw.OrderID, withdraw.Bonuses, withdraw.UploadedAt, withdraw.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// проверка при возникновении конфликта.
		userID, err := s.checkWithdraw(withdraw)
		if err != nil {
			return err
		}

		if userID == withdraw.UserID {
			return domain.ErrAlreadyUploadedByThisUser
		} else {
			return domain.ErrAlreadyUploadedByAnotherUser
		}
	}

	return nil
}

// Withdrawals возвращает все списания бонусов пользователя.
func (s *Storage) Withdrawals(ctx context.Context) ([]domain.Withdraw, error) {

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return nil, errors.New("incorrect user id")
	}

	var withdrawals []domain.Withdraw
	rows, err := s.db.Query("SELECT order_id, bonuses, uploaded_at FROM withdrawals WHERE user_id = $1 ORDER BY uploaded_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var withdraw domain.Withdraw
		err := rows.Scan(&withdraw.OrderID, &withdraw.Bonuses, &withdraw.UploadedAt)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdraw)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(withdrawals) == 0 {
		return nil, domain.ErrNoWithdraws
	}

	return withdrawals, nil
}

// Balance возвращает весь баланс пользователя.
func (s *Storage) Balance(ctx context.Context) (float32, error) {

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return 0, errors.New("incorrect user id")
	}

	var nullableBalance sql.NullFloat64
	err := s.db.QueryRow("SELECT SUM(bonuses) FROM orders WHERE user_id=$1", userID).
		Scan(&nullableBalance)
	if err != nil {
		return 0, err
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, err
}

// WithdrawBalance возвращает сумму списанных баллов пользователя.
func (s *Storage) WithdrawBalance(ctx context.Context) (float32, error) {

	userID, ok := ctx.Value(domain.UserIDKeyForContext).(int64)
	if !ok {
		return 0, errors.New("incorrect user id")
	}

	var nullableBalance sql.NullFloat64
	err := s.db.QueryRow("SELECT SUM(bonuses) FROM withdrawals WHERE user_id=$1", userID).
		Scan(&nullableBalance)
	if err != nil {
		return 0, err
	}
	if !nullableBalance.Valid {
		return 0, nil
	}

	balance := float32(nullableBalance.Float64)
	return balance, err
}

// heckWithdraw проверяет на конфликт списание.
func (s *Storage) checkWithdraw(withdraw domain.Withdraw) (int64, error) {
	var userID int64
	err := s.db.QueryRow("SELECT user_id FROM withdrawals WHERE order_id=$1", withdraw.OrderID).
		Scan(&userID)
	return userID, err
}
