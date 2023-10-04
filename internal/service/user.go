package service

import (
	"time"

	"github.com/AlexCorn999/bonus-system/internal/entities"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(user entities.User) error
}

type Users struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSecret []byte
}

func NewUsers(repo UserRepository, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
	}
}

func (u *Users) SignUp(usr entities.SighUpInput) error {
	password, err := u.hasher.Hash(usr.Password)
	if err != nil {
		return err
	}

	user := entities.User{
		Loggin:       usr.Login,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.Create(user)
}
