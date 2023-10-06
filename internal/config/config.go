package config

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port     string
	DBPort   string
	TokenTTL time.Duration
	LogLevel string
}

func NewConfig() *Config {
	return &Config{
		Port:     ":8080",
		TokenTTL: time.Minute * 30,
		LogLevel: "debug",
	}
}

// NetAddress структура для проверки флага -a.
type NetAddr struct {
	Host string
	Port int
}

func (a NetAddr) String() string {
	return a.Host + ":" + strconv.Itoa(a.Port)
}

func (a *NetAddr) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	port, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.Host = hp[0]
	a.Port = port
	return nil
}

// ParseFlags обрабатывает аргументы командной строки и сохраняет их значения в соответствующих переменных.
func (c *Config) ParseFlags() {
	port := new(NetAddr)
	_ = flag.Value(port)

	flag.Var(port, "a", "net address host:port")
	// host=127.0.0.1 port=5432 user=postgres sslmode=disable password=1234
	dbPort := flag.String("d", "", "port for database")

	flag.Parse()
	c.DBPort = *dbPort

	// проверка значения addr, чтобы записать в переменную Port.
	if port.String() != ":0" {
		c.Port = port.String()
	}

	// Установка данных адреса запуска HTTP-сервера через переменную окружения.
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		c.Port = envRunAddr
	}

	// Установка порта подключения к базе данных через переменную окружения.
	if envPath := os.Getenv("DATABASE_URI"); envPath != "" {
		c.DBPort = envPath
	}

}
