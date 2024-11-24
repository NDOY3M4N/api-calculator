package main

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

type wrapperWritter struct {
	http.ResponseWriter
	statusCode int
}

func generateRequestID() string {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// Handle error
		return ""
	}
	return hex.EncodeToString(b)
}

func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapper := &wrapperWritter{
			w, http.StatusOK,
		}

		next.ServeHTTP(wrapper, r)

		logger.Info("handling request",
			slog.Int("statusCode", wrapper.statusCode),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("requestID", generateRequestID()),
			slog.Any("duration", time.Since(start)),
		)
	})
}
