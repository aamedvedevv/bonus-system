package transport

import (
	"encoding/json"
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
		switch {
		case errors.Is(err, domain.ErrAlreadyUploadedByThisUser):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusOK)
			return
		case errors.Is(err, domain.ErrAlreadyUploadedByAnotherUser):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusConflict)
			return
		case errors.Is(err, domain.ErrIncorrectOrder):
			logError("orderUploading", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		default:
			logError("orderUploading", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

// GetAllOrders выводит отсортированный по дате список заказов пользователя.
func (s *APIServer) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := s.orders.GetAllOrders(r.Context())
	if err != nil {
		if errors.Is(err, domain.ErrNoData) {
			logError("getAllOrders", err)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		logError("getAllOrders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		logError("getAllOrders", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(ordersJSON)
}
