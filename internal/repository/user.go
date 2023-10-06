package repository

import (
	"github.com/AlexCorn999/bonus-system/internal/entities"
)

// Create добавляет пользователя в базу данных.
func (s *Storage) Create(user entities.User) error {
	result, err := s.db.Exec("INSERT INTO users (login, password, registered_at) values ($1, $2, $3) on conflict (login) do nothing",
		user.Login, user.Password, user.RegisteredAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDuplicate
	}

	return nil
}
