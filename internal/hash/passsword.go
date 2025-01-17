package hash

import (
	"crypto/sha1"
	"fmt"
)

// SHA1Hasher использует SHA1 для хэширования пароля с добавлением salt.
type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

// Hash хэширует пароль с добавлением salt.
func (h *SHA1Hasher) Hash(password string) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(password)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
