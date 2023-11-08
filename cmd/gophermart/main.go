package main

import (
	"log"

	"github.com/AlexCorn999/bonus-system/internal/config"
	"github.com/AlexCorn999/bonus-system/internal/transport"
)

// @title Накопительная система лояльности «Гофермарт»
// @version 1.0
// @description Накопительная система баллов лояльности. Система использует регистрацию и авторизацию пользователей. Занимается хранением и списанием баллов, параллельно обрабатывая номера заказов, путем обращения к стороннему источнику.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	server := transport.NewAPIServer(config.NewConfig())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
