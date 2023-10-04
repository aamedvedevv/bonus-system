package http

import (
	"net/http"

	"github.com/AlexCorn999/bonus-system/internal/config"
	"github.com/AlexCorn999/bonus-system/internal/hash"
	"github.com/AlexCorn999/bonus-system/internal/logger"
	"github.com/AlexCorn999/bonus-system/internal/repository"
	"github.com/AlexCorn999/bonus-system/internal/service"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type APIServer struct {
	config  *config.Config
	router  *chi.Mux
	logger  *log.Logger
	storage *repository.Storage
	users   *service.Users
}

func NewAPIServer(config *config.Config) *APIServer {
	return &APIServer{
		config: config,
		router: chi.NewRouter(),
		logger: log.New(),
	}
}

func (s *APIServer) Start() error {
	s.configureRouter()

	if err := s.configureLogger(); err != nil {
		return err
	}

	db, err := s.configureStore()
	if err != nil {
		return err
	}
	s.storage = db
	defer s.storage.Close()

	hasher := hash.NewSHA1Hasher("salt")
	s.users = service.NewUsers(db, hasher, []byte("sample secret"))

	s.logger.Info("starting api server")

	return http.ListenAndServe(s.config.Port, s.router)
}

func (s *APIServer) configureRouter() {
	s.router.Use(logger.WithLogging)
	s.router.Post("/api/user/register", s.SighUp)
	//s.router.Post("/api/user/login", s.SighIn)
}

func (s *APIServer) configureLogger() error {
	level, err := log.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configureStore() (*repository.Storage, error) {
	db, err := repository.NewStorage("host=127.0.0.1 port=5432 user=postgres sslmode=disable password=1234")
	if err != nil {
		return nil, err
	}
	return db, nil
}
