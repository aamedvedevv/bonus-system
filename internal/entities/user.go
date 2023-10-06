package entities

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"loggin"`
	Password     string    `json:"password"`
	RegisteredAt time.Time `json:"registered_at"`
}

func init() {
	validate = validator.New()
}

type SighUpInput struct {
	Login    string `json:"login" validate:"required,gte=2"`
	Password string `json:"password" validate:"required,gte=4"`
}

func (i *SighUpInput) Validate() error {
	return validate.Struct(i)
}
