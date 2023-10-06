package service

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/AlexCorn999/bonus-system/internal/domain"
	"github.com/golang-jwt/jwt/v4"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UserRepository interface {
	Create(user domain.User) error
	GetUser(login, password string) (domain.User, error)
}

type Users struct {
	repo   UserRepository
	hasher PasswordHasher

	hmacSecret []byte
	tokenTtl   time.Duration
}

func NewUsers(repo UserRepository, hasher PasswordHasher, secret []byte, ttl time.Duration) *Users {
	return &Users{
		repo:       repo,
		hasher:     hasher,
		hmacSecret: secret,
		tokenTtl:   ttl,
	}
}

// SignUp хэширует пароль пользователя и добавляет пользователя в базу данных.
func (u *Users) SignUp(usr domain.SighUpAndInInput) error {
	password, err := u.hasher.Hash(usr.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Login:        usr.Login,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	return u.repo.Create(user)
}

// SignIn проверяет наличие пользователя в базе данных и выписывает token.
func (u *Users) SignIn(usr domain.SighUpAndInInput) (string, error) {
	password, err := u.hasher.Hash(usr.Password)
	if err != nil {
		return "", err
	}

	user, err := u.repo.GetUser(usr.Login, password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrUserNotFound
		}
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.tokenTtl)),
	})

	return token.SignedString(u.hmacSecret)
}
