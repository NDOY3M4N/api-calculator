package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &wrapperWritter{
			w, http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		logger.Info("Log request",
			slog.Int("statusCode", wrapper.statusCode),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("requestID", generateRequestID()),
			slog.Any("duration", time.Since(start)),
		)
	})
}
