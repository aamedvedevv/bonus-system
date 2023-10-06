package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
	"github.com/AlexCorn999/bonus-system/internal/logger"
	"github.com/AlexCorn999/bonus-system/internal/repository"
)

func (s *APIServer) SighUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logger.LogError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.users.SignUp(usr); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			logger.LogError("signUp", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	w.WriteHeader(http.StatusOK)
}

// Login отвечает за аутентификацию пользователя по логину и паролю. Проверяет наличие токена.
func (s *APIServer) SighIn(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.users.SignIn(usr)
	if err != nil {
		// пользователь не найден.
		if errors.Is(err, domain.ErrUserNotFound) {
			logger.LogError("signIn", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		logger.LogError("signIn", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
	w.WriteHeader(http.StatusOK)
}
