package transport

import (
	"errors"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// OrderUploading загружает номер заказа в систему.
func (s *APIServer) OrderUploading(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		logError("orderUploading", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := s.orders.AddOrderID(r.Context(), string(data)); err != nil {
		if errors.Is(err, domain.ErrAlreadyUploadedByThisUser) {
			logError("orderUploading", err)
			w.WriteHeader(http.StatusOK)
			return
		} else if errors.Is(err, domain.ErrAlreadyUploadedByAnotherUser) {
			logError("orderUploading", err)
			w.WriteHeader(http.StatusConflict)
			return
		} else if errors.Is(err, domain.ErrIncorrectOrder) {
			logError("orderUploading", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		logError("orderUploading", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
