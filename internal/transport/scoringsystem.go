package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/domain"
)

// функция должна делать гет запросы по адресу и читать тело затем обновлять данные в базе
func (s *APIServer) ScoringSystem() {

	// получаем номер заказа из системы если его статус не PROCESSED . INVALID .
	orderID, err := s.scoringsystem.GetOrderStatus()
	if err != nil {
		//fmt.Println("Error GET ORDER ID:", err)
	}

	addr := fmt.Sprintf("%s/api/orders/%s", s.config.ScoringSystemPort, orderID)
	resp, err := http.Get(addr)
	if err != nil {
		//fmt.Println("Error GET ЗАПРОС:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("УСПЕХ")
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			//fmt.Println("Error READ BODY:", err)
		}

		var orderScoring domain.ScoringSystem
		if err := json.Unmarshal(data, &orderScoring); err != nil {
			//fmt.Println("Error JSON:", err)
		}

		if err := s.scoringsystem.UpdateOrder(orderScoring); err != nil {
			//fmt.Println("Error UPDATE:", err)
		}

	} else if resp.StatusCode == http.StatusNoContent {
		fmt.Println("Не зарегистрированный заказ")
	} else {
		fmt.Println("ДРУГАЯ ОШИБКА")
	}
}
