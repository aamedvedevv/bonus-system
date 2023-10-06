package repository

import "github.com/AlexCorn999/bonus-system/internal/domain"

// Create добавляет пользователя в базу данных.
func (s *Storage) AddOrder(order domain.Order) error {
	result, err := s.db.Exec("INSERT INTO orders (order_id, status, uploaded_at, bonuses, user_id) values ($1, $2, $3, $4, $5) on conflict (order_id) do nothing",
		order.OrderID, order.Status, order.UploadedAt, order.Bonuses, order.UserID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// проверка при возникновении конфликта.
		userID, err := s.checkOrder(order)
		if err != nil {
			return err
		}

		if userID == order.UserID {
			return domain.ErrAlreadyUploadedByThisUser
		} else {
			return domain.ErrAlreadyUploadedByAnotherUser
		}
	}

	return nil
}

// checkOrder проверяет на конфликт номер заказа.
func (s *Storage) checkOrder(order domain.Order) (int64, error) {
	var userID int64
	err := s.db.QueryRow("SELECT user_id FROM orders WHERE order_id=$1", order.OrderID).
		Scan(&userID)
	return userID, err
}
