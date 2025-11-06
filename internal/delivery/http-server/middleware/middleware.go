package middleware

import (
	"calendar-server/pkg/logger/zappretty"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

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
