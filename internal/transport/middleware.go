package transport

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type (
	ResponseData struct {
		Status int
		Size   int
	}

	LoggingResponseWriter struct {
		http.ResponseWriter
		ResponseData *ResponseData
	}
)

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}

func logFields(handler string) log.Fields {
	return log.Fields{
		"handler": handler,
	}
}

func logError(handler string, err error) {
	log.WithFields(logFields(handler)).Error(err)
}

// withLogging выполняет функцию middleware с логированием.
// Содержит сведения о URI, методе запроса и времени, затраченного на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func withLogging(next http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &ResponseData{
			Status: 0,
			Size:   0,
		}
		lw := LoggingResponseWriter{
			ResponseWriter: w,
			ResponseData:   responseData,
		}

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.WithFields(log.Fields{
			"uri":      r.RequestURI,
			"method":   r.Method,
			"duration": duration,
			"status":   responseData.Status,
			"size":     responseData.Size,
		}).Info("request details: ")
	}
	return http.HandlerFunc(logFn)
}

// authMiddleware выполняет функцию middleware авторизации.
// Получает токен из запроса и передает в контекст userID который совершает данный запрос.
func (s *APIServer) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID, err := s.users.ParseToken(r.Context(), token)
		if err != nil {
			logError("authMiddleware", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// getTokenFromRequest получает token из cookie.
func getTokenFromRequest(r *http.Request) (string, error) {
	token, err := r.Cookie("token")
	if err != nil {
		return "", err

	}
	return token.Value, nil
}
