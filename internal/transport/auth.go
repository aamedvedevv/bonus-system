package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
	"github.com/AlexCorn999/bonus-system/internal/repository"
)

// SighUp отвечает за регистрацию пользователя по логину и паролю. Автоматически производит аутентификацию.
func (s *APIServer) SighUp(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logError("signUp", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.users.SignUp(r.Context(), usr); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			w.WriteHeader(http.StatusConflict)
			return
		} else {
			logError("signUp", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	// атоматический авторизируем пользователя
	r.Body = io.NopCloser(bytes.NewBuffer(data))
	s.SighIn(w, r)
}

// SighIn отвечает за аутентификацию пользователя по логину и паролю. Проверяет наличие токена.
func (s *APIServer) SighIn(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var usr domain.SighUpAndInInput
	if err := json.Unmarshal(data, &usr); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := usr.Validate(); err != nil {
		logError("signIn", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := s.users.SignIn(r.Context(), usr)
	if err != nil {
		// пользователь не найден.
		if errors.Is(err, domain.ErrUserNotFound) {
			logError("signIn", err)
			w.WriteHeader(http.StatusUnauthorized)
		}
		logError("signIn", err)
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
