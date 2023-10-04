package main

import (
	"log"

	"github.com/AlexCorn999/bonus-system/internal/config"
	"github.com/AlexCorn999/bonus-system/internal/http"
)

func main() {
	server := http.NewAPIServer(config.NewConfig())
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
