package repository

import "github.com/AlexCorn999/bonus-system/internal/domain"

// GetOrderStatus получает orderID если его статус не PROCESSED или INVALID.
func (s *Storage) GetOrderStatus() (string, error) {
	var orderID string
	err := s.db.QueryRow("SELECT * FROM orders WHERE status NOT IN ('PROCESSED', 'INVALID') LIMIT 1").
		Scan(&orderID)
	return orderID, err
}

// UpdateOrder обновляет данные заказа. Начисляет бонусы и меняет статус.
func (s *Storage) UpdateOrder(order domain.ScoringSystem) error {
	_, err := s.db.Exec("UPDATE orders SET status=$1, bonuses=$2 WHERE order_id=$3", order.Status, order.Bonuses, order.OrderID)
	return err
}
