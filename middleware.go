package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type contextKey string

const requestIDKey contextKey = "requestID"

type Middleware func(http.HandlerFunc) http.HandlerFunc

type wrapperWritter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrapperWritter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func AddRequestId(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, err := generateRequestID()
		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		r = r.WithContext(ctx)

		if err != nil {
			writeError(w, r, http.StatusInternalServerError, err)
			return
		}

		w.Header().Set("X-Request-ID", reqID)

		next.ServeHTTP(w, r)
	}
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Context().Value(requestIDKey).(string)

		start := time.Now()
		wrappedWritter := &wrapperWritter{w, http.StatusOK}

		logger.Info("Log request",
			slog.Int("statusCode", wrappedWritter.statusCode),
			slog.Duration("duration", time.Since(start)),
			slog.String("remoteAddr", r.RemoteAddr),
			slog.Group("request",
				slog.String("id", reqID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			),
		)

		next.ServeHTTP(wrappedWritter, r)
	}
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}

		return next
	}
}

func generateRequestID() (string, error) {
	b := make([]byte, 3)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("Error generating request ID: %s", err)
	}

	return "req_" + hex.EncodeToString(b), nil
}
