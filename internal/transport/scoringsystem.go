package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// ScoringSystem выполняет GET запрос в систему расчета баллов и обновляет статус и кол-во бонусов за заказ.
func (s *APIServer) ScoringSystem() {
	// получаем номера заказов  до 15 шт из системы если их статус не PROCESSED или INVALID
	orderID, err := s.scoringsystem.GetOrderStatus(context.Background())
	if err != nil {
		logError("scoringSystem", err)
		return
	}

	// создаем ссылку для запроса GET

	for _, id := range orderID {

		addr := fmt.Sprintf("%s/api/orders/%s", s.config.ScoringSystemPort, id)
		resp, err := http.Get(addr)
		if err != nil {
			logError("scoringSystem", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				logError("scoringSystem", err)
				return
			}

			var orderScoring domain.ScoringSystem
			if err := json.Unmarshal(data, &orderScoring); err != nil {
				logError("scoringSystem", err)
				return
			}

			if err := s.scoringsystem.UpdateOrder(context.Background(), orderScoring); err != nil {
				logError("scoringSystem", err)
				return
			}
		}
	}

}
