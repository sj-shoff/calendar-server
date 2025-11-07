package middleware

import (
	"calendar-server/pkg/logger/zappretty"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseWriter - структура для записи статуса ответа
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader - метод для записи статуса ответа
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// LoggingMiddleware - middleware для логирования запросов
func LoggingMiddleware(log *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rw, r)

		log.Info("HTTP request",
			zappretty.Field("method", r.Method),
			zappretty.Field("path", r.URL.Path),
			zappretty.Field("status", rw.statusCode),
			zappretty.Field("duration", time.Since(start)),
			zappretty.Field("user_agent", r.UserAgent()),
			zappretty.Field("remote_addr", r.RemoteAddr),
		)
	})
}

// CORS middleware добавляет заголовки для CORS
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем заголовки CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
