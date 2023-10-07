package transport

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// Balance выводит сумму баллов лояльности и использованных за весь период регистрации баллов пользователя.
func (s *APIServer) Balance(w http.ResponseWriter, r *http.Request) {
	balance, err := s.withdraw.Balance(r.Context())
	if err != nil {
		logError("balance", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	balanceJSON, err := json.Marshal(balance)
	if err != nil {
		logError("balance", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(balanceJSON)
}

// Withdraw реализует списание бонусов пользователя в учет суммы нового заказа.
func (s *APIServer) Withdraw(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var withdraw domain.Withdraw
	if err := json.Unmarshal(data, &withdraw); err != nil {
		logError("withdraw", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.withdraw.Withdraw(r.Context(), withdraw); err != nil {
		if errors.Is(err, domain.ErrAlreadyUploadedByThisUser) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, domain.ErrAlreadyUploadedByAnotherUser) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, domain.ErrIncorrectOrder) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		} else if errors.Is(err, domain.ErrNoBonuses) {
			logError("withdraw", err)
			w.WriteHeader(http.StatusPaymentRequired)
			return
		}
		logError("withdraw", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}