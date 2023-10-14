package repository

import (
	"context"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// GetOrderStatus получает orderID если его статус не PROCESSED или INVALID.
func (s *Storage) GetOrderStatus(ctx context.Context) (string, error) {
	var orderID string
	err := s.db.QueryRowContext(ctx, "SELECT order_id FROM orders WHERE status NOT IN ('PROCESSED', 'INVALID') LIMIT 1").
		Scan(&orderID)
	return orderID, err
}

// UpdateOrder обновляет данные заказа. Начисляет бонусы и меняет статус.
func (s *Storage) UpdateOrder(ctx context.Context, order domain.ScoringSystem) error {
	_, err := s.db.ExecContext(ctx, "UPDATE orders SET status=$1, bonuses=$2 WHERE order_id=$3", order.Status, order.Bonuses, order.OrderID)
	return err
}
