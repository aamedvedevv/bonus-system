package repository

import (
	"database/sql"

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

// Balance возвращает весь баланс пользователя.
func (s *Storage) Balance(userID int64) (float32, error) {
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
func (s *Storage) WithdrawBalance(userID int64) (float32, error) {
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
