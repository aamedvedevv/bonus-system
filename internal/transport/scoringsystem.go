package transport

import (
	"fmt"
	"io"
	"net/http"
)

// функция должна делать гет запросы по адресу и читать тело затем обновлять данные в базе
func (s *APIServer) ScoringSystem() {
	addr := fmt.Sprintf("%s/api/orders/5555555555554444", s.config.ScoringSystemPort)
	resp, err := http.Get(addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response body:", string(body))

	// берем заказ из системы если его нету то ошибка 204 заказа нет в системе

	// если заказ есть то делаем гет запрос

	// читаем тело ответа и заносим в поля заказа

	// 429 — превышено количество запросов к сервису.
}
