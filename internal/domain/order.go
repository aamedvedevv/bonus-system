package domain

import (
	"errors"
	"time"
)

type OrderStatus string

var (
	ErrAlreadyUploadedByThisUser    = errors.New("the order number has already been uploaded by this user")
	ErrAlreadyUploadedByAnotherUser = errors.New("the order number has already been uploaded by another user")
	ErrIncorrectOrder               = errors.New("incorrect order id")
)

const (
	// заказ загружен в систему, но не попал в обработку.
	Registered OrderStatus = "NEW"
	// вознаграждение за заказ рассчитывается.
	Invalid OrderStatus = "PROCESSING"
	// система расчёта вознаграждений отказала в расчёте.
	Processing OrderStatus = "INVALID"
	// данные по заказу проверены и информация о расчёте успешно получена.
	Processed OrderStatus = "PROCESSED"
)

type Order struct {
	OrderID    string      `json:"order_id"`
	Status     OrderStatus `json:"status"`
	UploadedAt time.Time   `json:"uploaded_at"`
	Bonuses    float32     `json:"bonuses"`
	UserID     int64       `json:"user_id"`
}
