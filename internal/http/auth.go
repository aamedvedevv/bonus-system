package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/entities"
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

	var usr entities.SighUpInput
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

/*
// Login отвечает за аутентификацию пользователя по логину и паролю. Проверяет наличие токена.
func (s *APIServer) SighIn(w http.ResponseWriter, r *http.Request) {
	user, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// обработать разные статусы ответов
	token, err := s.auth.Login(user, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	payload, err := json.Marshal(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}*/
