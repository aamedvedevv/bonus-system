package repository

import (
	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// Create добавляет пользователя в базу данных.
func (s *Storage) Create(user domain.User) error {
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

// GetUser возвращает пользователя из базы данных.
func (s *Storage) GetUser(login, password string) (domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT id, login, password, registered_at FROM users WHERE login=$1 AND password=$2", login, password).
		Scan(&user.ID, &user.Login, &user.Password, &user.RegisteredAt)
	return user, err
}